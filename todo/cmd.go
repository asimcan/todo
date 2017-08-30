package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/jinzhu/now"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/lukasdietrich/todo/model"
)

var (
	duePattern = regexp.MustCompile(`(?i)\s*\bdue\s+(\S+)$`)
)

func CmdAdd(db *database) cli.Command {
	return cli.Command{
		Name:      "add",
		Usage:     "Add a new task",
		ArgsUsage: "[description] [due ...]",
		Action: func(c *cli.Context) error {
			content, err := parseInput(strings.Join(c.Args(), " "))
			if err != nil {
				return err
			}

			now := model.Now()
			task := model.Task{
				ID: model.NewGUID(),
				Metadata: model.Metadata{
					Created:  now,
					Modified: now,
				},
				Content: *content,
			}

			if err := db.Create(task).Error; err != nil {
				return errors.Wrap(err, "could not save task")
			}

			fmt.Printf("added task <%s>\n",
				color.CyanString("%s", task.ID))
			return nil
		},
	}
}

func parseInput(input string) (*model.Content, error) {
	var (
		index   = duePattern.FindStringSubmatchIndex(input)
		content = model.Content{
			Description: input,
			Due:         model.FromTime(now.EndOfDay()),
		}
	)

	if len(index) == 4 {
		due, err := parseDueDate(input[index[2]:])
		if err != nil {
			return nil, errors.Wrap(err, "invalid due date")
		}

		content.Due = model.FromTime(due)
		content.Description = input[:index[0]]
	}

	return &content, nil
}
