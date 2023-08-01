package surge

import (
	"errors"
	"net/http"
	"os"

	"github.com/yieldray/surgecli/api"
	"github.com/yieldray/surgecli/types"
	surgeUtils "github.com/yieldray/surgecli/utils"
)

type Surge struct {
	email      string
	token      string
	httpClient *http.Client
}

func (surge *Surge) GetEmailAndToken() (string, string) {
	return surge.email, surge.token
}

func (surge *Surge) SetEmailAndToken(email, token string) {
	surge.email = email
	surge.token = token
}

// new a Surge instance with default http client
// and automatically load emain&token from `~/.netrc`
// if there is still no token, please call `surge.Login()` to login
func New() *Surge {
	surge := &Surge{}
	surge.httpClient = http.DefaultClient

	if token := os.Getenv("SURGECLI_TOKEN"); token != "" {
		// try to load token from environment variables
		surge.email = `<no local email, use "surgecli account" to check from remote!>`
		surge.token = token
	} else {
		// try to load email and token from .netrc file
		email, token, _ := surgeUtils.ReadNetrc()
		surge.email = email
		surge.token = token
	}

	return surge
}

// change the default http client
func (surge *Surge) SetHTTPClient(client *http.Client) {
	surge.httpClient = client
}

// login and store email&token to `~/.netrc`
// just like what the official surge cli do
// if the username do not exists, surge.sh will create one for you :)
func (surge *Surge) Login(username, password string) (email string, err error) {
	t, err := api.Token(surge.httpClient, username, password)
	if err != nil {
		return "", err
	}

	surge.email = t.Email
	surge.token = t.Token

	// custom config file
	surgeUtils.ConfAddAccount(t.Email, t.Token)

	return t.Email, surgeUtils.WriteNetrc(t.Email, t.Token)
}

// logout and clear the `~/.netrc`
func (surge *Surge) Logout() (email string, err error) {
	email = surge.email

	if email == "" {
		return email, errors.New("not logged-in")
	}

	if err = surgeUtils.RemoveNetrc(); err != nil {
		return
	} else {
		surge.email = ""
		surge.token = ""
		return
	}
}

func (surge *Surge) List() (types.List, error) {
	return api.List(surge.httpClient, surge.token)
}

func (surge *Surge) Upload(domain, src string, onEventStream func(byteLine []byte)) error {
	return api.Upload(surge.httpClient, surge.token, domain, src, onEventStream)
}

func (surge *Surge) Account() (types.Account, error) {
	return api.Account(surge.httpClient, surge.token)
}

func (surge *Surge) Teardown(domain string) (types.Teardown, error) {
	return api.Teardown(surge.httpClient, surge.token, domain)
}

func (surge *Surge) Whoami() string {
	return surge.email
}
