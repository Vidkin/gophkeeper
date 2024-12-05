/*
Copyright © 2024 MIKHAIL SIRKIN <skim991@gmail.com>
*/

// Package cmd contains the commands for the GophKeeper client application.
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Vidkin/gophkeeper/internal/client"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register [login] [password]",
	Short: "Register user",
	Long: `Register user in GophKeeper. For example:
	client register login password`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.Register(args[0], args[1]); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
