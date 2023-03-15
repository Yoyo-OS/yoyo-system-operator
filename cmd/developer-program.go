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
)

func devProgramUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	Join the developers program

Usage:
  	yso developer-program [flags]

Flags:
	--help/-h		show this message

Examples:
	yso developer-program
`)
	return nil
}

func NewDevProgramCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "developer-program",
		Short: "Join the developers program",
		RunE:  devProgram,
	}
	cmd.SetUsageFunc(devProgramUsage)
	return cmd
}

func devProgram(cmd *cobra.Command, args []string) error {
	fmt.Println("Developer program is not yet available.")
	return nil
}
