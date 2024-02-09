package main

import (
	"fmt"

	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
	"encoding/json"
)

type Log struct {
	Server string
	Host string
	Request string
	Status int
	Size int
	RequestTime float64
}


func main() {
	fmt.Println("Start syslog server")
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC3164)
	server.SetHandler(handler)
	server.ListenUDP("0.0.0.0:5114")

	server.Boot()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			parser(logParts)
		}
	}(channel)

	server.Wait()
}

func parser(log format.LogParts ) {
	var ms Log
	json.Unmarshal([]byte(log["content"].(string)), &ms)
	fmt.Printf("request: %s status: %d\n", ms.Request, ms.Status)
}
