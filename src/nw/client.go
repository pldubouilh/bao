package nw

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	utils "github.com/pldubouilh/bao/src/utils"
	"golang.org/x/crypto/ssh"
)

func copyConns(m string, i int, connA net.Conn, connB net.Conn) {
	_, err := io.Copy(connA, connB)
	if err != nil && i < 8 {
		i++
		time.Sleep(time.Duration(i) * 200 * time.Millisecond)
		copyConns(m, i, connA, connB)
	} else if err != nil {
		utils.PrintMaybe("m", err)
	}
}

func forward(conn net.Conn, remote string, c *utils.BaoConfig) {
	if !c.Connected || !c.Wanted {
		return
	}

	sshConn, err := c.SSHConn.Dial("tcp", remote)
	if err != nil {
		utils.PrintMaybe("ssh dial failed", err)
		return
	}

	go copyConns("ssh to local failed", 0, sshConn, conn)
	go copyConns("local to ssh failed", 0, conn, sshConn)
}

func acceptConns(l net.Listener, local string, remote string, c *utils.BaoConfig) {
	// defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil && strings.Contains(err.Error(), "already in use") {
			utils.DieMaybe("", err)
		} else if err != nil || !c.Connected || !c.Wanted {
			return
		}

		go forward(conn, remote, c)
	}
}

func cbHostKeyCheck(hostname string, remote net.Addr, key ssh.PublicKey, c *utils.BaoConfig) error {
	hostKey := key.Marshal()

	for _, s := range c.Checksums {
		pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(s))
		utils.DieMaybe("cant read local host checksum", err)

		if bytes.Equal(hostKey, pubKey.Marshal()) {
			return nil
		}
	}

	utils.DieMaybe("", errors.New("@@@@ no valid host @@@@"))
	return nil
}

func healthcheck(c *utils.BaoConfig) {
	hc := func(c *utils.BaoConfig) {
		c.MightBeDead = true
		_, _, err := c.SSHConn.SendRequest("keepalive@golang.org", true, nil)
		if err == nil {
			c.MightBeDead = false
		}
	}

	for {
		time.Sleep(5 * time.Second)

		if !c.Wanted {
			return
		} else if !c.Connected || c.MightBeDead {
			go attemptConn(c)
		} else {
			go hc(c)
		}
	}
}

func attemptConn(c *utils.BaoConfig) {
	var err error
	c.SSHConn, err = ssh.Dial("tcp", c.Addr, c.SSHConfig)

	if err != nil {
		fmt.Println(c.Nickname, "conn seems off:", err)
		c.Connected = false
	} else {
		fmt.Println(c.Nickname, "(re)connected")
		c.Connected = true
		c.MightBeDead = false
	}

	c.Event <- true
}

// Kill kills client
func Kill(c *utils.BaoConfig) {
	fmt.Println(c.Nickname, "will disconnect")
	c.SSHConn.Close()
	c.Wanted = false
	c.Connected = false
	c.MightBeDead = false
	for _, l := range c.LocalConns {
		l.Close()
	}
	c.LocalConns = nil
	c.SSHConfig = nil
	c.SSHConn = nil
	c.Event <- true
}

// New spins a new bao client
func New(c *utils.BaoConfig) {
	fmt.Println(c.Nickname, "will connect")
	c.Wanted = true
	privKey, err := ssh.ParsePrivateKey([]byte(c.Privkey))
	utils.DieMaybe("cant read private key", err)

	c.Event <- true

	c.SSHConfig = &ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privKey),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return cbHostKeyCheck(hostname, remote, key, c)
		},
		Timeout: 4 * time.Second,
	}

	attemptConn(c)
	go healthcheck(c)

	for _, s := range c.Forwards {
		local := "127.0.0.1:" + strings.Split(s, ":")[0]
		remote := strings.SplitN(s, ":", 2)[1]

		l, err := net.Listen("tcp", local)
		utils.DieMaybe("cant listen on local port", err)

		fmt.Println(c.Nickname, "setup", local, "<>", remote)
		c.LocalConns = append(c.LocalConns, l)
		go acceptConns(l, local, remote, c)
	}
}
