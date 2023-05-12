package api

import (
	"encoding/json"
	"net/http"

	"github.com/yieldray/surgecli/types"
)

func List(client *http.Client, token string) (list types.List, err error) {
	list = types.List{}

	req, _ := http.NewRequest("GET", "https://surge.surge.sh/list", nil)
	req.SetBasicAuth("token", token)

	if res, err := client.Do(req); err != nil {
		return list, err
	} else {
		json.NewDecoder(res.Body).Decode(&list)
		return list, nil
	}

}
