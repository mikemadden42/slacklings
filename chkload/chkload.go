package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

func main() {
	token := os.Getenv("SLACK_TOKEN")
	if token != "" {
		host, _ := os.Hostname()
		channel := "#build_status"
		subject := "High load on " + host
		body := ""
		load := ""

		file, err := os.Open("/proc/loadavg")
		checkErr(err)

		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			load = strings.Split(line, " ")[2]
			body += load + "\n"
		}

		currentLoad, _ := strconv.ParseFloat(load, 64)
		if currentLoad > float64(runtime.NumCPU()) {
			postAlert(channel, subject, body, token)
		}
	} else {
		fmt.Println("SLACK_TOKEN is not set! Set SLACK_TOKEN before running again.")
		return
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func postAlert(channel string, subject string, body string, token string) {
	api := slack.New(token)
	params := slack.PostMessageParameters{AsUser: true}
	attachment := slack.Attachment{Text: body}
	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage(channel, subject, params)
	checkErr(err)
	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
}
