package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// SlackRequestBody ... Struct for JSON message request to Slack
type SlackRequestBody struct {
	Text string `json:"text"`
}

// CheckCurrentIPAddress ... function to check current public IP address
func CheckCurrentIPAddress() string {

	ip, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		fmt.Println("Connection Lost\n")
	}
	defer ip.Body.Close()
	responseData, err := ioutil.ReadAll(ip.Body)
	responseString := string(responseData)

	return responseString
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
		return errors.New("Slack Did Not Responded!")
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
			oldIP = currentIPAddress
			webhookURL := "https://hooks.slack.com/services/<YOUR TOKEN HERE>"
			msg := "New IP Address: " + currentIPAddress
			err := SlackNotification(webhookURL, msg)
			if err != nil {
				fmt.Println("Could Not Send Message!")
			}
		}

	}
}
