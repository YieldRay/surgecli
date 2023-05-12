package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/yieldray/surgecli/types"
)

func Token(client *http.Client, username, password string) (token types.Token, err error) {
	token = types.Token{}

	req, _ := http.NewRequest("POST", "https://surge.surge.sh/token", nil)
	req.SetBasicAuth(username, password)

	if res, err := client.Do(req); err != nil {
		return token, err
	} else {
		switch res.StatusCode {
		case 200, 201:
			{
				if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
					return token, err
				} else {
					return token, nil
				}
			}
		case 401:
			{
				tokenErr := types.TokenError{}
				if err = json.NewDecoder(res.Body).Decode(&tokenErr); err != nil {
					return token, err
				} else {
					return token, errors.New(strings.Join(tokenErr.Messages, " & "))
				}
			}
		default:
			{
				// test
				fmt.Println(res)
				b, _ := io.ReadAll(res.Body)
				fmt.Println(string(b))

				return token, errors.New(string(b))
			}
		}
	}

}
