package main

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
	"github.com/yoyo-os/yso/cmd"
)

var (
	Version = "1.3.4-1"
)

func help(cmd *cobra.Command, args []string) {
	fmt.Println(`Usage:
	yso [flags] [command] [arguments]

Global Flags:
	-h, --help            	Show this help message and exit

Commands:
	config              	Configure YSO
	create-task             Create a new task
	delete-task             Delete a task
	developer-program   	Join the developers program
	help                	Show this help message and exit
	list-tasks          	List all tasks
	rotate-tasks		Rotate tasks
	trigger-update	  	Trigger a system update
	update-check	  	Check for system updates
	version             	Show version and exit`)
}

func newYsoCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "yso",
		Short:   "YSO is an utility which allows you to perform maintenance tasks on your Yoyo OS installation.",
		Version: Version,
	}
}

func main() {
	rootCmd := newYsoCommand()

	rootCmd.AddCommand(cmd.NewCreateTaskCommand())
	rootCmd.AddCommand(cmd.NewDeleteTaskCommand())
	rootCmd.AddCommand(cmd.NewConfigCommand())
	rootCmd.AddCommand(cmd.NewDevProgramCommand())
	rootCmd.AddCommand(cmd.NewListTasksCommand())
	rootCmd.AddCommand(cmd.NewRotateTasksCommand())
	rootCmd.AddCommand(cmd.NewTriggerUpdateCommand())
	rootCmd.AddCommand(cmd.NewCheckUpdateCommand())
	rootCmd.SetHelpFunc(help)
	rootCmd.Execute()
}
