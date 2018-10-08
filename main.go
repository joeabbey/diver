package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joeabbey/diver/cmd"
)

// Version is populated from the Makefile and is tied to the release TAG
var Version string

// Build is the last GIT commit
var Build string

// READONLY will disable any commands that write
var READONLY string

func init() {
	if READONLY == "" {
		READONLY = os.Getenv("DIVER_RO")
		if READONLY == "" {
			READONLY = "false"
		}
	}
	b, err := strconv.ParseBool(READONLY)
	if err != nil {
		log.Fatalf("%v", err)
	}
	cmd.DiverVersion = Version
	cmd.DiverBuild = Build
	cmd.DiverRO = b
}

func main() {

	cmd.Execute()
}
