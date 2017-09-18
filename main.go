package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Name:  "boom",
		Usage: "make an explosive entrance",
		Action: func(c *cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
	}

	app.Run(os.Args)
}
