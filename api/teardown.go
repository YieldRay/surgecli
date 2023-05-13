package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yieldray/surgecli/types"
)

func Teardown(client *http.Client, token, domain string) (types.Teardown, error) {
	teardown := types.Teardown{}

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("https://surge.surge.sh/%s", domain), nil)
	req.SetBasicAuth("token", token)

	if res, err := client.Do(req); err != nil {
		return teardown, err
	} else {
		json.NewDecoder(res.Body).Decode(&teardown)
		return teardown, nil
	}

}
