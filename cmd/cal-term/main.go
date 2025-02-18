package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "cal-term"
	app.Usage = "display your calendar in the terminal"

	app.Action = func(c *cli.Context) {
		fmt.Println("This is the main part")
    if err := viewAction(); err != nil {
      fmt.Fprintln(os.Stderr, "Error:", err)
    }
	}

	app.Commands = []cli.Command{
		{
			Name:      "authenticate",
			ShortName: "auth",
			Usage:     "authorize and store google credentials via Oauth",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "clientId",
					Usage:  "google app client ID",
					EnvVar: "GOOGLE_CLIENT_ID",
				},
				cli.StringFlag{
					Name:   "clientSecret",
					Usage:  "google app client secrtet",
					EnvVar: "GOOGLE_CLIENT_SECRET",
				},
			},
			Action: func(c *cli.Context) {
				if err := authAction(c.String("clientId"), c.String("clientSecret")); err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err)
				}
				fmt.Fprintln(os.Stdout, "🎉 Successfully authenticated!")
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
