/*
Copyright © 2024 MIKHAIL SIRKIN <skim991@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Vidkin/gophkeeper/internal/client"
)

// authCmd represents the authorize command
var authCmd = &cobra.Command{
	Use:   "auth [login] [password]",
	Short: "Authorize user",
	Long: `Authorize user in GophKeeper and get JWT. For example:
	client auth login password`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.Auth(args[0], args[1]); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
