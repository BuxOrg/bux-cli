/*
BUX CLI

Author: MrZ © 2023 github.com/BuxOrg/bux-cli

This CLI tool can help you interact with a BUX server or database.

Help contribute via GitHub!
*/
package main

import (
	"github.com/BuxOrg/bux-cli/cmd"
)

// main will load the all the commands and kick-start the application
func main() {
	cmd.Execute()
}
