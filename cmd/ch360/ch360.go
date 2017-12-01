package main

import (
	//"fmt"

	"github.com/CloudHub360/ch360.go"
	"github.com/docopt/docopt-go"
)

func main() {
	usage := `CloudHub360 command-line tool.

Usage:
  ch360 -h | --help
  ch360 --version

Options:
  -h --help     Show this help message.
  --version     Show version.`

	docopt.Parse(usage, nil, true, ch360.Version, false)

}
