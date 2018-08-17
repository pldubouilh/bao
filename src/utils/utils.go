package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)

// BaoConfigStr is used to hardcode a config file
var BaoConfigStr = ``

// BaoConfig is the main type for all things config
type BaoConfig struct {
	Event       chan (bool)
	MightBeDead bool
	Connected   bool
	Wanted      bool
	SSHConfig   *ssh.ClientConfig
	SSHConn     *ssh.Client
	LocalConns  []net.Listener
	Nickname    string   `json:"nickname"`
	Username    string   `json:"username"`
	Addr        string   `json:"addr"`
	Forwards    []string `json:"forwards"`
	Privkey     string   `json:"privkey"`
	Checksums   []string `json:"checksums"`
}

// WriteConfig writes config
func WriteConfig(path string, w BaoConfig) {
	b, err := json.Marshal(w)
	DieMaybe("cant marshall object", err)
	ioutil.WriteFile(path, b, 0644)
}

// ReadConfig reads config
func parseConfig(f []byte) (*BaoConfig, error) {
	var c BaoConfig
	errM := json.Unmarshal(f, &c)

	if errM != nil {
		errM = errors.New("uh cant read config")
		fmt.Println("uh cant read config")
		return &c, errM
	}

	c.Wanted = false
	c.Connected = false
	c.Event = make(chan bool)
	c.MightBeDead = false
	return &c, errM
}

// ReadConfigs Read all configs either in base dir or embedded
func ReadConfigs() *[]*BaoConfig {
	var cs []*BaoConfig

	pathLocal, _ := homedir.Expand("~/.ssh/bao")
	pathLocal = pathLocal + "/"
	filesLocal, er := ioutil.ReadDir(pathLocal)

	pathMac, _ := os.Executable()
	pathMac = strings.TrimSuffix(pathMac, "bao")
	filesMac, erm := ioutil.ReadDir(pathMac)

	if len(BaoConfigStr) != 0 {
		fmt.Println("using embedded config")
		c, erp := parseConfig([]byte(BaoConfigStr))
		if erp == nil {
			cs = append(cs, c)
		}
	} else if er == nil && len(filesLocal) != 0 {
		fmt.Println("reading local filesLocal")
		for _, f := range filesLocal {
			payload, err := ioutil.ReadFile(pathLocal + f.Name())
			c, erp := parseConfig(payload)
			if err == nil && erp == nil {
				cs = append(cs, c)
			}
		}
	} else if erm == nil && len(filesMac) != 0 {
		fmt.Println("reading local filesMac")
		for _, f := range filesMac {
			payload, err := ioutil.ReadFile(pathMac + f.Name())
			c, erp := parseConfig(payload)
			if err == nil && erp == nil {
				cs = append(cs, c)
			}
		}
	} else {
		DieMaybe("", errors.New("no config found"))
	}

	return &cs
}

// DieMaybe dies if errs
func DieMaybe(m string, err error) {
	if err != nil {
		log.Fatalf(m, err)
	}
}

// PrintMaybe prints if errs
func PrintMaybe(m string, err error) {
	if err != nil {
		fmt.Println(m, err)
	}
}

// DummyEventListener spins a dummy event listener if not used (e.g. cli mode)
func DummyEventListener(c *BaoConfig) {
	for {
		<-c.Event // some event !
	}
}
