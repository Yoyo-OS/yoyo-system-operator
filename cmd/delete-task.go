package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: YSO is a utility which allows you to perform maintenance
	tasks on your Yoyo OS installation.
*/

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yoyo-os/yso/core"
)

func deleteTaskUsage(*cobra.Command) error {
	fmt.Print(`Description: 
  Delete a task

Usage:
	yso delete-task [flags] [options]

Flags:
	--help/-h		show this message

Examples:
	yso delete-task my-task
	yso delete-task "my task"
`)

	return nil
}

func NewDeleteTaskCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-task",
		Short: "Delete a task",
		RunE:  deleteTask,
	}
	cmd.SetUsageFunc(deleteTaskUsage)

	return cmd
}

func deleteTask(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	err := core.DeleteTaskByUnitName(args[0])
	if err != nil {
		fmt.Println("Task", args[0], "does not exist")
		return err
	}

	fmt.Println("Task", args[0], "deleted")

	return nil
}
