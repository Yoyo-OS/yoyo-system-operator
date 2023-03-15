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

func rotateTasksUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	  Rotate tasks

Usage:
	yso rotate-tasks [flags]

Flags:
	--help/-h		show this message

Examples:
	yso rotate-tasks
`)

	return nil
}

func NewRotateTasksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate-tasks",
		Short: "Rotate tasks",
		RunE:  rotateTasks,
	}
	cmd.SetUsageFunc(rotateTasksUsage)
	cmd.Flags().StringP("private-event", "e", "", "private event to rotate tasks for (boot, shutdown, login, logout)")

	return cmd
}

func rotateTasks(cmd *cobra.Command, args []string) error {
	event, err := cmd.Flags().GetString("private-event")
	if err != nil {
		event = "no-system-event"
	}

	return core.RotateTasks(event)
}
