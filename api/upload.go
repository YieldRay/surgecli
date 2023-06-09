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
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	surgeUtils "github.com/yieldray/surgecli/utils"
)

// src 必须为绝对路径！故返回值也是绝对路径
// 返回一个目录的父目录，正斜线分隔路径
// 例如：/a/b 返回 /a，/a/b/ 返回 /a
func toParent(src string) string {
	const SEP = string(filepath.Separator)
	p := src + SEP
	p = filepath.Dir(p)
	p = filepath.Dir(p)
	return filepath.ToSlash(p) + "/"
}

// client =  &http.Client{}
// domain = domainplaceholder.surge.sh
// src = <the directory path>
// onEventStream <jsonString=>void>
func Upload(client *http.Client, token, domain, src string, onEventStream func(byteLine []byte)) (err error) {
	if !surgeUtils.IsDir(src) {
		return errors.New("not a directory")
	}
	// 获取绝对路径，保证tar是压缩了一个文件夹而不是其内容（当src为当前目录时）
	src, _ = filepath.Abs(src)
	// 获取src的父目录路径，写入tar时用绝对路径删除此前缀
	parentSrc := toParent(src)

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)

	var fileCount, projectSize int64
	fileCount = 0
	projectSize = 0

	ignoreList := surgeIgnore(src)

	err = filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 路径转换为 unix 格式
		unixPath := filepath.ToSlash(path)

		// 默认跳过文件夹
		if info.IsDir() && info.Name() != "." {
			defaultIgnoreList := []string{".git", ".*", ".*.*~", "node_modules", "bower_components"}
			for _, pattern := range defaultIgnoreList {
				matched, _ := filepath.Match(pattern, info.Name())
				// fmt.Println(info.Name(), pattern, matched)
				if matched {
					return filepath.SkipDir
				}
			}
		}

		// 自定义文件跳过
		if filterForIgnore(ignoreList, unixPath, info) {
			return nil
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

		// 写入文件名：dir是一级根目录
		name, cut := strings.CutPrefix(unixPath, parentSrc)
		hdr.Name = strings.TrimPrefix(name, string(filepath.Separator))
		if !cut {
			// 正常情况下必定能够移除父目录的路径，否则panic
			panic(fmt.Sprintf("path=%s parentPath=%s", unixPath, parentSrc))
			// TODO:应该有更好的方法来移除父目录路径
		}

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
	req.Header.Add("version", "0.23.1")
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

// 读取目录下的 .surgeignore 文件
func surgeIgnore(src string) []string {
	ignoreList := []string{}

	p := filepath.Join(src, ".surgeignore")

	file, err := os.Open(p)
	if err != nil {
		return ignoreList
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		ignoreList = append(ignoreList, scanner.Text())
	}

	return ignoreList
}

// 返回 ture 则被跳过
func filterForIgnore(ignoreList []string, path string, info fs.FileInfo) bool {
	if strings.HasPrefix(info.Name(), ".") {
		return true
	}

	for _, pattern := range ignoreList {
		matched, _ := filepath.Match(pattern, info.Name())
		if matched {
			return true
		}
	}

	return false
}
