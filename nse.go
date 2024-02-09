package main

import (
	"fmt"
	"strings"

	"encoding/json"
	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type Log struct {
	Server      string  `json:"server"`
	Host        string  `json:"host"`
	Request     Request  `json:"request"`
	Status      int     `json:"status"`
	Size        int     `json:"body_bytes_sent"`
	RequestTime float64 `json:"request_time"`
}

type Request struct {
	Method string
	Url string
	HttpVersion string
}

func (d *Request) UnmarshalJSON(data []byte) error {
	parts := strings.Split(string(data), " ")
	if len(parts) < 3 {
		return fmt.Errorf("request does not have 3 part")
	}

	d.Method = parts[0]
	d.Url = parts[1]
	d.HttpVersion = parts[2]
	return nil
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

func parser(log format.LogParts) {
	var ms Log
	json.Unmarshal([]byte(log["content"].(string)), &ms)
	fmt.Printf("request: %s status: %d\n", ms.Request.Url, ms.Status)
	//fmt.Println(log["content"].(string))
}
