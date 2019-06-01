package main

import "github.com/zdoherty/rainyqt/pkg/version"

var Build string

func main() {
	if Build != "" {
		version.RainyqtVersion.Build = Build
	}
}
