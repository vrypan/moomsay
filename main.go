package main

import "github.com/vrypan/moomsay/cmd"

var VERSION string

func main() {
	cmd.VERSION = VERSION
	cmd.Execute()
}
