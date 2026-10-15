/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"

	"github.com/CALEXCO/goautomate/cmd"
	"github.com/CALEXCO/goautomate/pkg/server"
	"github.com/goccy/go-yaml"
)

var tagetServers = server.NewTargetServers()

func main() {
	cmd.Execute()
}

func example() {
	var server = server.NewServer("1.1.1.1", "", "")
	tagetServers.Servers = append(tagetServers.Servers, server)
	tagetServers.Servers = append(tagetServers.Servers, server)
	tagetServers.Servers = append(tagetServers.Servers, server)

	bytes, err := yaml.Marshal(tagetServers)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
	fmt.Printf("%v", tagetServers)
}
