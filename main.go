package main

import (
	"fmt"

	"github.com/pldubouilh/bao/src/ui"
	// "github.com/pldubouilh/bao/src/nw"
	// "github.com/pldubouilh/bao/src/utils"
)

func main() {
	fmt.Println("bao starting")

	// Start in UI mode
	ui.Start()

	// Headless mode - will connect to embedded file or the first file in ~/.ssh/bao
	// cs := *utils.ReadConfigs()
	// go nw.New(cs[0])
	// utils.DummyEventListener(cs[0])
}
