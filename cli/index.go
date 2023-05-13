package cli

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/yieldray/surgecli/surge"
)

type privateSurgeCLI struct { // 单例模式
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