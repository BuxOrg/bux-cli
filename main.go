/*
BUX CLI

This CLI app is used for interacting with BUX databases or servers.

Learn more about BUX: https://GetBux.io
*/
package main

import (
	"github.com/BuxOrg/bux-cli/cmd"
)

// main will load the all the commands and kick-start the application
func main() {
	cmd.Execute()
}
