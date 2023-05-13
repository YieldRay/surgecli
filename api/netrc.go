package api

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/bgentry/go-netrc/netrc"
)

// .netrc should be 0400 in unix
// but it will cause problem for windows
// so we use 0600 here

//```~/.netrc
//machine surge.surge.sh
//    login email@example.net
//    password here_is_token
//```

var netrcPath string

func init() {
	if runtime.GOOS == "windows" {
		home := os.Getenv("USERPROFILE")
		netrcPath = fmt.Sprintf("%s\\%s", home, ".netrc")
	} else {
		home := os.Getenv("HOME")
		netrcPath = path.Join(home, ".netrc")
	}

}

func saveMyNetrc(myNetrc *netrc.Netrc) error {
	b, err := myNetrc.MarshalText()
	if err != nil {
		return err
	}

	if err = os.WriteFile(netrcPath, b, 0600); err != nil {
		return err
	}

	return nil
}

func RemoveNetrc() error {
	if myNetrc, err := netrc.ParseFile(netrcPath); err != nil {
		return err
	} else {
		myNetrc.RemoveMachine("surge.surge.sh")

		// Save token to .netrc
		return saveMyNetrc(myNetrc)
	}
}

func isExist(path string) bool {
	if s, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	} else {
		return s.Mode().IsRegular()
	}
}

func WriteNetrc(login, password string) error {
	if !isExist(netrcPath) {
		// Create .netrc file if not exists
		os.WriteFile(netrcPath, []byte{}, 0600)
	}
	if myNetrc, err := netrc.ParseFile(netrcPath); err != nil {
		return err
	} else {
		// Save or update machine=surge.surge.sh
		if m := myNetrc.FindMachine("surge.surge.sh"); m != nil {
			m.UpdateLogin(login)
			m.UpdatePassword(password)
		} else {
			myNetrc.NewMachine("surge.surge.sh", login, password, "")
		}

		// Save token to .netrc
		return saveMyNetrc(myNetrc)
	}
}

func ReadNetrc() (login, password string, err error) {
	if netrc, err := netrc.ParseFile(netrcPath); err != nil {
		return "", "", err
	} else {
		if m := netrc.FindMachine("surge.surge.sh"); m != nil {
			return m.Login, m.Password, nil
		} else {
			return "", "", errors.New("no existed machine surge.surge.sh")
		}
	}

}
