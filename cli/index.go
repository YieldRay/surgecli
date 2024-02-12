package cli

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/yieldray/surgecli/surge"
	"github.com/yieldray/surgecli/utils"
)

var FLAG_DEBUG = 0
var Commands = []*cli.Command{}

// the only http client for all http request
var httpClient = &http.Client{Transport: utils.RoundTrip(func(req *http.Request) (*http.Response, error) {
	API_HOST := utils.ConfGetApi()
	// debug http
	if FLAG_DEBUG > 0 {
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
	if FLAG_DEBUG > 0 {
		if b, e := httputil.DumpRequest(req, true); e != nil {
			log.Fatalln(e)
		} else {
			log.Println(string(b))
		}
	}

	// send the actual request
	res, err := http.DefaultClient.Do(req)

	// debug http response
	if FLAG_DEBUG > 0 {
		if b, e := httputil.DumpResponse(res, true); e != nil {
			log.Fatalln(e)
		} else {
			log.Println(string(b))
			log.Println()
		}
	}

	return res, err
})}

var surgesh = surge.New()

func init() {
	surgesh.SetHTTPClient(httpClient)
}
