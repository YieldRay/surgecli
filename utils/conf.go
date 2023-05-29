package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

var ConfPath string

func init() {
	ConfPath = AtHome(".surgecli")
}

type SurgeCliConf struct {
	Api      string            `json:"api"`
	Accounts map[string]string `json:"accounts"`
}

func readConf() SurgeCliConf {
	var c SurgeCliConf
	b, _ := os.ReadFile(ConfPath)
	json.Unmarshal(b, &c)
	if c.Accounts == nil {
		c.Accounts = map[string]string{}
	}
	return c
}

func writeConf(c SurgeCliConf) error {
	b, _ := json.Marshal(c)
	return os.WriteFile(ConfPath, b, 0600)
}

func modConf(cb func(c *SurgeCliConf)) error {
	obj := readConf()
	cb(&obj)
	return writeConf(obj)
}

func ConfSetApi(api string) error {
	return modConf(func(c *SurgeCliConf) {
		if len(api) > 0 {
			c.Api = api
		} else {
			c.Api = "https://surge.surge.sh"
		}
	})
}

func ConfGetApi() string {
	c := readConf()
	api := c.Api
	if len(api) > 0 {
		return api
	} else {
		return "https://surge.surge.sh"
	}
}

func ConfAddAccount(email, token string) error {
	return modConf(func(c *SurgeCliConf) {
		c.Accounts[email] = token
	})
}

func ConfGetEmailList() []string {
	c := readConf()
	list := make([]string, 0, len(c.Accounts))
	for email := range c.Accounts {
		list = append(list, email)
	}
	return list
}

func ConfUseEmail(email string) error {
	c := readConf()
	token, ok := c.Accounts[email]
	if ok {
		return WriteNetrc(email, token)
	} else {
		return fmt.Errorf("email %s is not in the email list read from config file", email)
	}
}
