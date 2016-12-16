package main

// http://stackoverflow.com/questions/32987215/find-numbers-in-string-using-golang-regexp
// http://stackoverflow.com/questions/6878590/the-maximum-value-for-an-int-type-in-go
// http://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"regexp"
	"strconv"

	"github.com/nlopes/slack"
)

func main() {
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		log.Fatal(err)
	}

	var processes map[string]int
	processes = make(map[string]int)

	re := regexp.MustCompile("^[0-9]+")
	for _, file := range files {
		matches := re.FindAllString(file.Name(), -1)
		if file.IsDir() && len(matches) > 0 {
			loginuid, _ := ioutil.ReadFile("/proc/" + file.Name() + "/loginuid")
			uid := string(loginuid[:])
			if uid == "4294967295" {
				continue
			}
			processes[uid]++
		}
	}

	token := os.Getenv("SLACK_TOKEN")
	if token != "" {
		host, _ := os.Hostname()
		channel := "#build_status"
		subject := "High process count on " + host
		body := ""

		for key, value := range processes {
			if value >= 40 {
				user, _ := user.LookupId(key)
				body += "user " + user.Username + " -> " + strconv.Itoa(value) + " processes\n"
				postAlert(channel, subject, body, token)
				body = ""
			}
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
