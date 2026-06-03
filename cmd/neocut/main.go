package main

import (
	"neocut/internal/cmd"
	"neocut/internal/config"
)

var (
	Version        string
	Commit         string
	PublisherName  string
	PublisherEmail string
)

func main() {
	config.SetBuildInfo(Version, Commit, PublisherName, PublisherEmail)
	cmd.Execute()
}
