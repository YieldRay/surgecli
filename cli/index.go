package cli

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/surge"
)

type privateSurgeCLI struct { // 单例模式，此结构体仅实现一次
	surgesh  *surge.Surge
	API_HOST string
	DEBUG    int
}

var SurgeCLI *privateSurgeCLI // 单例模式

type transport struct { // 自定义http客户端所需
}

func (tsc transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 若 SurgeCLI.DEBUG > 0 开启 DEBUG，则打印 http 请求及响应
	if SurgeCLI.DEBUG > 0 {
		log.Printf("API_HOST = %s\n", SurgeCLI.API_HOST)
	}

	// 替换为命令行参数的API
	if u, e := url.Parse(strings.Replace(req.URL.String(), "https://surge.surge.sh", SurgeCLI.API_HOST, 1)); e != nil {
		fmt.Println(e)
	} else {
		req.URL = u
		req.Host = u.Host
	}

	if SurgeCLI.DEBUG > 0 {
		if b, e := httputil.DumpRequest(req, true); e != nil {
			log.Fatalln(e)
		} else {
			log.Println(string(b))
		}
	}

	// 发送实际请求
	res, err := http.DefaultClient.Do(req)

	if SurgeCLI.DEBUG > 0 {
		if b, e := httputil.DumpResponse(res, true); e != nil {
			log.Fatalln(e)
		} else {
			log.Println(string(b))
			log.Println()
		}
	}

	return res, err
}

func init() {
	surgesh := surge.New()

	surgesh.SetHTTPClient(&http.Client{Transport: &transport{}})
	SurgeCLI = &privateSurgeCLI{
		surgesh:  surgesh,
		API_HOST: "",
		DEBUG:    0,
	}
}


// 通过反射获取命令配置数组
func (surgecli *privateSurgeCLI) Commands() []*cli.Command {
	cmds := []*cli.Command{}

	t := reflect.TypeOf(surgecli)

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if createCommandFunc, ok := m.Func.Interface().(func(c *privateSurgeCLI) *cli.Command); ok {
			cmds = append(cmds, createCommandFunc(surgecli))
		}
	}

	return cmds
}
