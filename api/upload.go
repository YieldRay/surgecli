package api

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"k8s.io/helm/pkg/ignore"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	utils "github.com/yieldray/surgecli/utils"
)

// client =  &http.Client{}
// domain = domainplaceholder.surge.sh
// src = <the directory path>
// onEventStream <jsonString=>void>
func Upload(client *http.Client, token, domain, src string, onEventStream func(byteLine []byte)) (err error) {
	if !utils.IsDir(src) {
		return errors.New("not a directory")
	}
	// 获取绝对路径，保证tar是压缩了一个文件夹而不是其内容（当src为当前目录时）
	src, _ = filepath.Abs(src)

	computeTarPath := computeTarPathFn(src)
	computeIgnore := computeIgnoreFn(src)

	buf := new(bytes.Buffer)
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)

	var fileCount, projectSize int64
	fileCount = 0
	projectSize = 0

	err = filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		tarPath := computeTarPath(path)
		isIgnore := computeIgnore(path, tarPath)

		if isIgnore {
			if info.IsDir() {
				return filepath.SkipDir // skip entire dir
			}
			return nil // skip this file
		}

		// 不处理非标准文件
		if !info.Mode().IsRegular() {
			return nil
		}

		// 计算文件数量和大小
		if !info.IsDir() {
			fileCount += 1
			projectSize += info.Size()
		}

		// 写入tar头
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		// 写入文件名
		hdr.Name = tarPath

		// 写入文件信息
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		// 打开文件
		fr, err := os.Open(path)
		if err != nil {
			return err
		}

		// copy 文件数据到 tw
		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}

		return fr.Close()
	})

	if err != nil {
		return err
	}

	// 关闭流
	if err = tw.Close(); err != nil {
		return err
	}
	if err = gw.Close(); err != nil {
		return err
	}

	// 构造 PUT 上传请求
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://surge.surge.sh/%s", domain), buf)

	if err != nil {
		return err
	}

	req.SetBasicAuth("token", token)
	req.Header.Add("file-count", fmt.Sprint(fileCount))
	req.Header.Add("project-size", fmt.Sprint(projectSize))
	req.Header.Add("timestamp", string(time.Now().UTC().Format(time.RFC3339)))
	req.Header.Add("version", Version)
	req.Header.Add("user-agent", "")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	// 若上传不成功，将res读取为字符串返回error
	if res.StatusCode == 403 {
		return fmt.Errorf("%s %s", "you do not have permission to publish to ", domain)
	}
	if res.StatusCode != 200 {
		b, err := io.ReadAll(res.Body)

		if err != nil {
			return err
		}
		return fmt.Errorf("%s", b)

		// TODO：读取返回的json再封装为error

		// m := make(map[string]interface{})
		// if err = json.Unmarshal(b, &m); err != nil {
		// 	return errors.New(res.Status)
		// }
	}

	reader := bufio.NewReader(res.Body)

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		onEventStream(line)

	}

	return nil
}

func computeTarPathFn(src string) func(path string) string {
	const SEP = string(filepath.Separator)
	p := src + SEP
	p = filepath.Dir(p)
	p = filepath.Dir(p)
	p = filepath.ToSlash(p) + "/"

	// 移除父路径，仅留下作为tar名称的路径 (Unix格式)
	return func(filePath string) string {
		unixPath := filepath.ToSlash(filePath)
		name, _ := strings.CutPrefix(unixPath, p)
		return strings.TrimPrefix(name, string(filepath.Separator))
	}
}

func computeIgnoreFn(src string) func(fullPath, tarPath string) bool {
	p := filepath.Join(src, ".surgeignore")
	rules, rulesErr := ignore.ParseFile(p)
	return func(fullPath, tarPath string) bool {
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			return true // we can not access the file, so ignore
		}

		// 这是surge文档规定的默认ignore列表
		defaultIgnoreList := []string{".git", ".*", ".*.*~", "node_modules", "bower_components"}
		for _, pattern := range defaultIgnoreList {
			matched, _ := filepath.Match(pattern, filepath.Base(fullPath))
			if matched {
				return true
			}
		}

		if rulesErr != nil {
			return false // rules got error, do not ignore
		}
		return rules.Ignore(removeTopParent(tarPath), fileInfo)
	}
}

func removeTopParent(path string) string {
	dir := filepath.Dir(path)
	return filepath.ToSlash(filepath.Join(filepath.Base(dir), filepath.Base(path)))
}
