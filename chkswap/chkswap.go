package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

func main() {
	token := os.Getenv("SLACK_TOKEN")
	if token != "" {
		host, _ := os.Hostname()
		channel := "#build_status"
		subject := "High swap on " + host
		body := ""
		size := 0.0
		used := 0.0
		percentUsed := 0.0

		file, err := os.Open("/proc/swaps")
		checkErr(err)

		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Filename") {
				continue
			}
			words := strings.Fields(line)
			used, _ = strconv.ParseFloat(words[3], 64)
			size, _ = strconv.ParseFloat(words[2], 64)
			percentUsed = used / size
			body += words[3] + "\n"
		}

		if percentUsed > 10.00 {
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
