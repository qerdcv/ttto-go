package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/urfave/cli/v2"

	_ "github.com/qerdcv/ttto/docs"
)

//	@title			Swagger of ttt-online
//	@version		0.0.1
//	@description	Tic-Tac-Toe online
func main() {
	app := &cli.App{
		Name:  "ttto",
		Usage: "ttto cli",
		Commands: cli.Commands{
			{
				Name:   "run",
				Usage:  "run an application",
				Action: run,
			},
			{
				Name:   "migrate",
				Usage:  "migrates application database",
				Action: runMigration,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
