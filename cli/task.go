package cli

import (
	"fmt"
	"task-manager/model"
	"log"
	"os"
	"time"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var taskApp *cli.App

func TaskCLI() {
	taskApp = &cli.App{
		Name:     "Task Manager",
		Version:  "v1.0",
		Compiled: time.Now(),
		Usage:    "Task management tool",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Action: func(ctx *cli.Context) error {
					fmt.Println("List all tasks...")
					tasks, err := model.GetAllTasks(bson.D{{}})
					if err != nil {
						log.Fatal(err)
					}
					model.PrintTasks(tasks)
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name: "completed",
						Action: func(ctx *cli.Context) error {
							tasks, err := model.FilterTasks(true)
							if err != nil {
								log.Fatal(err)
							}
							model.PrintTasks(tasks)
							return nil
						},
					},
					{
						Name: "uncompleted",
						Action: func(ctx *cli.Context) error {
							tasks, err := model.FilterTasks(false)
							if err != nil {
								log.Fatal(err)
							}
							model.PrintTasks(tasks)
							return nil
						},
					},
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() == 0 {
						color.Red.Println("Please provide a task title!")
						return nil
					}

					title := ctx.Args().First()
					task := model.Task{
						ID:        primitive.NewObjectID(),
						Title:     title,
						Completed: false,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}

					if err := model.AddTask(task); err != nil {
						color.Red.Println("Error adding task", err)
					} else {
						color.Green.Println("Task added successfully!")
					}
					return nil
				},
			},
			{
				Name:    "complete",
				Aliases: []string{"cpt"},
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() == 0 {
						color.Red.Println("Provide task id to complete task!")
						return nil
					}

					id := ctx.Args().First()
					if err := model.CompleteTask(id); err != nil {
						color.Red.Println("Error completing task", err)
					} else {
						color.Green.Println("Task completed successfully")
					}
					return nil
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"rm"},
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() == 0 {
						color.Red.Println("Provide task id to delete task!")
						return nil
					}

					id := ctx.Args().First()
					if err := model.DeleteTask(id); err != nil {
						color.Red.Println("Error deleting task", err)
					} else {
						color.Green.Println("Task deleted successfully")
					}
					return nil
				},
			},
		},
	}
}

func Run() error {
	if err := taskApp.Run(os.Args); err != nil {
		return err
	}
	return nil
}
