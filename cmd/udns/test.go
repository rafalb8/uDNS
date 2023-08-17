package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type TestCmd struct {
	Port     string `arg:"" help:"DNS server port" default:"53"`
	HttpPort string `arg:"" help:"HTTP server port" default:"8367"`
}

func (t *TestCmd) Run() error {
	req, err := http.NewRequest("GET", "http://udns.local:"+t.HttpPort, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
	}
	fmt.Printf("%s\n", bodyText)
	return nil
}
