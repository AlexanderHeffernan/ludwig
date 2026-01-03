package main

import (
	"flag"
	"fmt"
	"ludwig/internal/cli"
)

var version = "dev"

func main() {
	versionFlag := flag.Bool("version", false, "Print the version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println("ludwig version " + version)
		return
	}

	cli.StartInteractive()
}
