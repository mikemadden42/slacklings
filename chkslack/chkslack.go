package main

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "item",
			Value: "users",
			Usage: "item to list",
		},
	}

	token := os.Getenv("SLACK_TOKEN")
	if token != "" {
		api := slack.New(token)
		app.Action = func(c *cli.Context) error {
			if c.String("item") == "users" {
				users, err := api.GetUsers()
				checkErr(err)

				fmt.Println("id,username,realname,email,admin,owner,twofactor")
				for _, user := range users {
					fmt.Printf("%s,%s,%s,%s,%t,%t,%t\n", user.ID, user.Name, user.RealName, user.Profile.Email, user.IsAdmin, user.IsOwner, user.Has2FA)
				}
			} else if c.String("item") == "channels" {
				channels, err := api.GetChannels(false)
				checkErr(err)

				fmt.Println("id,name")
				for _, channel := range channels {
					fmt.Printf("%s,%s\n", channel.ID, channel.Name)
				}
			}
			return nil
		}
		app.Run(os.Args)
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
