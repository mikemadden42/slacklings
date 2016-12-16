package main

// http://stackoverflow.com/questions/20108520/get-amount-of-free-disk-space-using-go

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/nlopes/slack"
)

func main() {
	file, err := os.Open("/proc/filesystems")
	checkErr(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)
	filesystems := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "nodev") {
			continue
		} else {
			fs := strings.TrimSpace(line)
			filesystems = append(filesystems, fs)
		}
	}

	mtab, err := os.Open("/etc/mtab")
	checkErr(err)

	defer mtab.Close()

	scanner = bufio.NewScanner(mtab)
	mounts := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		mount := fields[1]
		filesys := fields[2]
		for _, i := range filesystems {
			if i == filesys {
				mounts = append(mounts, mount)
			}
		}
	}

	token := os.Getenv("SLACK_TOKEN")
	if token != "" {
		var stat syscall.Statfs_t
		host, _ := os.Hostname()
		channel := "#build_status"
		subject := "High disk utilization on " + host
		body := ""

		for _, m := range mounts {
			if m == "/var/cache/ccache" {
				continue
			}

			err := syscall.Statfs(m, &stat)
			checkErr(err)

			free := stat.Bfree * uint64(stat.Bsize)
			total := stat.Blocks * uint64(stat.Bsize)
			percentUsed := (float64(total-free) / float64(total))

			if percentUsed > 0.90 {
				body += m + "\n"
				body += strconv.FormatFloat(percentUsed, 'f', 2, 64) + "\n"
				postAlert(channel, subject, body, token)
			}

			body = ""
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
