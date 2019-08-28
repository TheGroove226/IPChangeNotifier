package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	externalip "github.com/glendc/go-external-ip"
)

// SlackRequestBody ... Struct for JSON message request to Slack
type SlackRequestBody struct {
	Text string `json:"text"`
}

// CheckCurrentIPAddress ... function to check current public IP address
func CheckCurrentIPAddress() string {

	consensus := externalip.DefaultConsensus(nil, nil)

	ip, err := consensus.ExternalIP()
	if err == nil {
		fmt.Println(ip.String()) // print IPv4/IPv6 in string format
	}
	return ip.String()

}

// SlackNotification ... Function to inform you that IP has changed
func SlackNotification(webhookURL string, msg string) error {
	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Error!")
	}
	return nil
}

func main() {

	fmt.Println("Starting Program...\n")
	fmt.Println("First Time Checking ...\n")

	oldIP := CheckCurrentIPAddress()

	fmt.Println("IP Address is:", oldIP)

	for {
		time.Sleep(2 * time.Second)
		currentIPAddress := CheckCurrentIPAddress()

		if currentIPAddress == oldIP {
			fmt.Println("No changes detected.\n")
		} else {
			webhookURL := "https://hooks.slack.com/services/<YOUR TOKEN HERE>"
			msg := "IP address has changed. New IP Address: " + currentIPAddress
			err := SlackNotification(webhookURL, msg)
			if err != nil {
				fmt.Println("Could Not Send Message!")
			}
			oldIP = currentIPAddress
		}
	}
}
