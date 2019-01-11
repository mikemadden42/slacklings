package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
	gojenkins "github.com/yosida95/golang-jenkins"
)

type jsonObject struct {
	Config configType
}

type configType struct {
	User  string
	Token string
	URL   string
}

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "item",
			Value: "jobs",
			Usage: "item to list",
		},
	}

	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var jsontype jsonObject
	json.Unmarshal(file, &jsontype)
	user := jsontype.Config.User
	token := jsontype.Config.Token
	url := jsontype.Config.URL

	auth := &gojenkins.Auth{
		Username: user,
		ApiToken: token,
	}
	jenkins := gojenkins.NewJenkins(auth, url)

	app.Action = func(c *cli.Context) error {
		if c.String("item") == "jobs" {
			jobs, err := jenkins.GetJobs()
			checkErr(err)
			for _, job := range jobs {
				fmt.Println(job.Name, job.Color)
			}
		} else if c.String("item") == "nodes" {
			computers, err := jenkins.GetComputers()
			checkErr(err)
			for _, computer := range computers {
				fmt.Println(computer.DisplayName, computer.Offline, computer.MonitorData.ArchitectureMonitor)
			}
		}
		return nil
	}
	app.Run(os.Args)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
