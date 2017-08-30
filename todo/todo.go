package main

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func main() {
	db, err := initDatabase()
	if err != nil {
		color.Red("%s\n", err)
		return
	}
	defer db.Close()

	app := cli.NewApp()

	app.HideVersion = true
	app.Usage = "a minimal commandline todo list"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "no-color",
			Usage:       "Disable colored output",
			Destination: &color.NoColor,
		},
	}

	app.Commands = []cli.Command{
		CmdList(db),
		CmdAdd(db),
		CmdDo(db),
	}

	if err := app.Run(os.Args); err != nil {
		color.Red("%s\n", err)
	}
}

func initDatabase() (*database, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	file := filepath.Join(user.HomeDir, ".todo", "todo.sqlite")
	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return nil, err
	}

	return open(file)
}
