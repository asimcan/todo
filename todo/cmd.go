package main

import (
	"fmt"
	"regexp"
	"strconv"
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

func CmdDo(db *database) cli.Command {
	return cli.Command{
		Name:      "do",
		Usage:     "Mark a task as completed",
		ArgsUsage: "[task number]",
		Action: func(c *cli.Context) error {
			task, err := findTask(db, c.Args())
			if err != nil {
				return err
			}

			if task.Completed != nil {
				return errors.New("task already completed")
			}

			now := model.Now()
			err = db.Model(&task).Updates(&model.Metadata{Completed: &now}).Error
			if err != nil {
				return errors.Wrap(err, "could not update task")
			}

			fmt.Printf("marked as completed <%s>\n",
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

func findTask(db *database, args []string) (*model.Task, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	num, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.Wrap(err, "could not parse task number")
	}

	var task model.Task
	err = db.pending().
		Offset(num - 1).
		Limit(1).
		First(&task).
		Error

	return &task, err
}
