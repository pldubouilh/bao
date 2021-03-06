package ui

import (
	"errors"
	"strings"
	"time"

	"github.com/pldubouilh/bao/src/nw"
	"github.com/pldubouilh/bao/src/utils"

	systray "github.com/getlantern/systray"
)

// Start starts ui
func Start() {
	systray.Run(func() { onReady() }, onExit)
}

func getClicks(m *systray.MenuItem, c *utils.BaoConfig) {
	for {
		<-m.ClickedCh
		// m.Check()
		// m.Disable()
		if !c.Wanted {
			nw.New(c)
		} else {
			nw.Kill(c)
		}
	}
}

func watchEvts(m *systray.MenuItem, c *utils.BaoConfig) {
	for {
		<-c.Event
		if (c.Wanted && !c.Connected) || c.MightBeDead {
			m.SetTitle("❔ " + c.Nickname + " ")
		} else if c.Connected {
			m.SetTitle("✅ " + c.Nickname + " ")
		} else {
			m.SetTitle("❌ " + c.Nickname + " ")
		}
	}
}

func onReady() {
	systray.SetIcon(BaoIcon())
	cs := *utils.ReadConfigs()

	for _, c := range cs {
		m := systray.AddMenuItem("❌ "+c.Nickname, "connect to "+c.Nickname)
		go getClicks(m, c)
		go watchEvts(m, c)
	}

	systray.AddSeparator()

	mService := systray.AddMenuItem("Open first service", "Open first service")
	firstService := "http://127.0.0.1:" + strings.Split(cs[0].Forwards[0], ":")[0]
	go func(f string) {
		for {
			<-mService.ClickedCh
			utils.OpenBrowser(f)
		}
	}(firstService)

	systray.AddSeparator()

	mInfo := systray.AddMenuItem("Info", "More info")
	go func() {
		for {
			<-mInfo.ClickedCh
			utils.OpenBrowser("https://github.com/pldubouilh/bao")
		}
	}()

	mQuit := systray.AddMenuItem("Quit", "Quits this app")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	// connect first service and open client
	nw.New(cs[0])
	go func(f string) {
		time.Sleep(2 * time.Second)
		utils.OpenBrowser(f)
	}(firstService)
}

func onExit() {
	utils.DieMaybe("", errors.New("time to go now"))
}
