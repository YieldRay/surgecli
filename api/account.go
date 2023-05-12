package api

import (
	"encoding/json"
	"net/http"

	"github.com/yieldray/surgecli/types"
)

func Account(client *http.Client, token string) (account types.Account, err error) {
	account = types.Account{}

	req, _ := http.NewRequest("GET", "https://surge.surge.sh/account", nil)
	req.SetBasicAuth("token", token)

	if res, err := client.Do(req); err != nil {
		return account, err
	} else {
		json.NewDecoder(res.Body).Decode(&account)
		return account, nil
	}

}
