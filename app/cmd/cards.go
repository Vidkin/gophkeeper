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
	cardID int64
	card   proto.BankCard
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

var addCardCmd = &cobra.Command{
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

var getCardCmd = &cobra.Command{
	Use:   "get [flags]",
	Short: "Get bank card by ID from GophKeeper",
	Long: `This command allows you to get bank card info by ID from your account in GophKeeper. For example:
	- client cards get --id 9`,
	Run: func(cmd *cobra.Command, args []string) {
		if cardID < 0 {
			fmt.Println("You must provide a bank card ID")
			os.Exit(1)
		}
		if err := client.GetCard(cardID); err != nil {
			fmt.Println(err)
		}
	},
}

var removeCardCmd = &cobra.Command{
	Use:   "remove [flags]",
	Short: "Remove bank card by ID from GophKeeper",
	Long: `This command allows you to remove bank card info by ID from your account in GophKeeper. For example:
	- client cards remove --id 9`,
	Run: func(cmd *cobra.Command, args []string) {
		if cardID < 0 {
			fmt.Println("You must provide a bank card ID")
			os.Exit(1)
		}
		if err := client.RemoveCard(cardID); err != nil {
			fmt.Println(err)
		}
	},
}

var getAllCardsCmd = &cobra.Command{
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
	addCardCmd.PersistentFlags().StringVar(&card.Owner, "owner", "", "bank card owner")
	addCardCmd.PersistentFlags().StringVar(&card.Cvv, "cvv", "", "bank card CVV")
	addCardCmd.PersistentFlags().StringVar(&card.ExpireDate, "expire", "", "bank card expire date")
	addCardCmd.PersistentFlags().StringVar(&card.Number, "number", "", "bank card number")
	addCardCmd.PersistentFlags().StringVar(&card.Description, "desc", "", "bank card description")

	getCardCmd.PersistentFlags().Int64Var(&cardID, "id", -1, "bank card id")

	removeCardCmd.PersistentFlags().Int64Var(&cardID, "id", -1, "bank card id")

	cardsCmd.AddCommand(getCardCmd)
	cardsCmd.AddCommand(removeCardCmd)
	cardsCmd.AddCommand(addCardCmd)
	cardsCmd.AddCommand(getAllCardsCmd)
	rootCmd.AddCommand(cardsCmd)
}
