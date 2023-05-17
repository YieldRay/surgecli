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
	surgeUtils "github.com/yieldray/surgecli/utils"
)

type privateSurgeCLI struct { // Singleton, only create once
	surgesh *surge.Surge
	DEBUG   int
}

var SurgeCLI *privateSurgeCLI // Singleton instance

type transport struct { // only for *http.client
}

func (tsc transport) RoundTrip(req *http.Request) (*http.Response, error) {
	API_HOST := surgeUtils.ConfGetApi()
	// debug http
	if SurgeCLI.DEBUG > 0 {
		log.Printf("API_HOST = %s\n", API_HOST)
	}

	// replace api host
	if u, e := url.Parse(strings.Replace(req.URL.String(), "https://surge.surge.sh", API_HOST, 1)); e != nil {
		fmt.Println(e)
	} else {
		req.URL = u
		req.Host = u.Host
	}

	// debug http request
	if SurgeCLI.DEBUG > 0 {
		if b, e := httputil.DumpRequest(req, true); e != nil {
			log.Fatalln(e)
		} else {
			log.Println(string(b))
		}
	}

	// send the actual request
	res, err := http.DefaultClient.Do(req)

	// debug http response
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

// the only http client for all http request
var customHttpClient *http.Client

func init() {
	customHttpClient = &http.Client{Transport: &transport{}}
}

// customHttpClient is readonly, so wrap it by a function
func CustomHttpClient() *http.Client {
	return customHttpClient
}

func init() {
	surgesh := surge.New()
	surgesh.SetHTTPClient(customHttpClient)
	SurgeCLI = &privateSurgeCLI{
		surgesh: surgesh,
		DEBUG:   0,
	}
}

// use reflect to get all sub commands as a slice
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
