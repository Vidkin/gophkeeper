/*
Copyright © 2024 MIKHAIL SIRKIN <skim991@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Vidkin/gophkeeper/internal/client"
	"github.com/Vidkin/gophkeeper/proto"
)

var (
	credID      int64
	credentials proto.Credentials
)

// credentialsCmd represents the user credentials management command
var credentialsCmd = &cobra.Command{
	Use:   "credentials [command] [flags]",
	Short: "User credentials management",
	Long: `User credentials management in GophKeeper. For example:
	- client credentials get credID
	- client credentials getAll
	- client credentials add credentials info
	- client credentials remove credID`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var addCredentialCmd = &cobra.Command{
	Use:   "add [flags]",
	Short: "Add a new user credentials to GophKeeper",
	Long: `This command allows you to add a new user credentials to your account in GophKeeper. For example:
	- client credentials add --login Login --pass Password --desc Description`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.AddCredentials(&credentials); err != nil {
			fmt.Println(err)
		}
	},
}

var getCredentialsCmd = &cobra.Command{
	Use:   "get [flags]",
	Short: "Get user credentials by ID from GophKeeper",
	Long: `This command allows you to get user credentials by ID from your account in GophKeeper. For example:
	- client credentials get --id 9`,
	Run: func(cmd *cobra.Command, args []string) {
		if credID < 0 {
			fmt.Println("You must provide a credential ID")
			os.Exit(1)
		}
		if err := client.GetCredentials(credID); err != nil {
			fmt.Println(err)
		}
	},
}

var removeCredentialsCmd = &cobra.Command{
	Use:   "remove [flags]",
	Short: "Remove user credentials by ID from GophKeeper",
	Long: `This command allows you to remove user credentials info by ID from your account in GophKeeper. For example:
	- client credentials remove --id 9`,
	Run: func(cmd *cobra.Command, args []string) {
		if credID < 0 {
			fmt.Println("You must provide a credential ID")
			os.Exit(1)
		}
		if err := client.RemoveCredentials(credID); err != nil {
			fmt.Println(err)
		}
	},
}

var getAllCredentialsCmd = &cobra.Command{
	Use:   "getAll",
	Short: "Get all user credentials from GophKeeper",
	Long: `This command allows you to get all user credentials from your account in GophKeeper. For example:
	- client credentials getAll`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.GetAllCredentials(); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	addCredentialCmd.PersistentFlags().StringVar(&credentials.Login, "login", "", "login")
	addCredentialCmd.PersistentFlags().StringVar(&credentials.Password, "pass", "", "password")
	addCredentialCmd.PersistentFlags().StringVar(&credentials.Description, "desc", "", "credentials description")

	getCredentialsCmd.PersistentFlags().Int64Var(&credID, "id", -1, "credentials id")
	removeCredentialsCmd.PersistentFlags().Int64Var(&credID, "id", -1, "credentials id")

	credentialsCmd.AddCommand(getCredentialsCmd)
	credentialsCmd.AddCommand(removeCredentialsCmd)
	credentialsCmd.AddCommand(addCredentialCmd)
	credentialsCmd.AddCommand(getAllCredentialsCmd)
	rootCmd.AddCommand(credentialsCmd)
}
