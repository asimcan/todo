package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/jinzhu/now"
	"github.com/peterh/liner"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/lukasdietrich/todo/model"
)

var (
	duePattern      = regexp.MustCompile(`(?i)\s*\bdue\s+(\S+)$`)
	keywordsPattern = regexp.MustCompile(`(\+\S+)\b`)
)

func CmdList(db *database) cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "List pending tasks",
		Action: func(c *cli.Context) error {
			var list []model.Task
			if err := db.pending().Find(&list).Error; err != nil {
				return err
			}

			for i, task := range list {
				check := ' '
				if task.Completed != nil {
					check = 'x'
				}

				desc := keywordsPattern.ReplaceAllStringFunc(task.Description,
					func(s string) string {
						return color.MagentaString("%s", s)
					})

				due := formatDueDate(task.Due.Time(), true)

				fmt.Printf("%s %s %s %s\n",
					color.CyanString("%2d", i+1),
					color.YellowString("(%c)", check),
					color.CyanString("%11s", due),
					desc,
				)
			}

			return nil
		},
	}
}

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

func CmdEdit(db *database) cli.Command {
	return cli.Command{
		Name:      "edit",
		Usage:     "Edit an existing task",
		ArgsUsage: "[task number]",
		Action: func(c *cli.Context) error {
			task, err := findTask(db, c.Args())
			if err != nil {
				return err
			}

			line := liner.NewLiner()
			line.SetCtrlCAborts(true)

			defer line.Close()

			taskString := fmt.Sprintf("%s due %s",
				task.Description,
				formatDueDate(task.Due.Time(), false))

			input, err := line.PromptWithSuggestion("", taskString, -1)
			if err != nil {
				return errors.Wrap(err, "edit canceled")
			}

			edited, err := parseInput(input)
			if err != nil {
				return err
			}

			update := model.Task{
				Metadata: model.Metadata{
					Modified: model.Now(),
				},
				Content: *edited,
			}

			err = db.Model(task).Updates(&update).Error
			if err != nil {
				return errors.Wrap(err, "could not update task")
			}

			fmt.Printf("edited task <%s>\n",
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

func CmdArchive(db *database) cli.Command {
	return cli.Command{
		Name:  "archive",
		Usage: "Archive completed tasks",
		Action: func(c *cli.Context) error {
			res := db.Model(&model.Task{}).
				Where("archived is null").
				Where("completed is not null").
				UpdateColumn("archived", model.Now())
			if err := res.Error; err != nil {
				return errors.Wrap(err, "could not archive tasks")
			}

			fmt.Printf("%s tasks archived\n",
				color.CyanString("%d", res.RowsAffected))
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
