/*
Copyright © 2024 MIKHAIL SIRKIN <skim991@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Vidkin/gophkeeper/internal/client"
	"github.com/Vidkin/gophkeeper/proto"
)

// cardsCmd represents the bank cards management command
var cardsCmd = &cobra.Command{
	Use:   "cards [command] [flags]",
	Short: "Bank cards management",
	Long: `Bank cards management in GophKeeper. For example:
	- client cards get cardID
	- client cards getAll
	- client cards add cardInfo
	- client cards remove cardID`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var card proto.BankCard
var addCmd = &cobra.Command{
	Use:   "add [flags]",
	Short: "Add a new bank card to GophKeeper",
	Long: `This command allows you to add a new bank card to your account in GophKeeper. For example:
	- client cards add --owner "Name Surname" --cvv 123 --expire 2024-12-26 --number 78878877 --desc "Test card"`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.AddCard(&card); err != nil {
			fmt.Println(err)
		}
	},
}

var getAllCmd = &cobra.Command{
	Use:   "getAll",
	Short: "Get all bank cards from GophKeeper",
	Long: `This command allows you to get all bank cards from your account in GophKeeper. For example:
	- client cards getAll`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.GetAllCards(); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	addCmd.PersistentFlags().StringVar(&card.Owner, "owner", "", "bank card owner")
	addCmd.PersistentFlags().StringVar(&card.Cvv, "cvv", "", "bank card CVV")
	addCmd.PersistentFlags().StringVar(&card.ExpireDate, "expire", "", "bank card expire date")
	addCmd.PersistentFlags().StringVar(&card.Number, "number", "", "bank card number")
	addCmd.PersistentFlags().StringVar(&card.Description, "desc", "", "bank card description")

	cardsCmd.AddCommand(addCmd)
	cardsCmd.AddCommand(getAllCmd)
	rootCmd.AddCommand(cardsCmd)
}
