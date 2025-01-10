/*
Copyright © 2024 MIKHAIL SIRKIN <skim991@gmail.com>
*/

// Package cmd contains the commands for the GophKeeper client application.
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Vidkin/gophkeeper/internal/client"
	"github.com/Vidkin/gophkeeper/proto"
)

var (
	noteID int64
	note   proto.Note
)

// notesCmd represents the user notes management command
var notesCmd = &cobra.Command{
	Use:   "notes [command] [flags]",
	Short: "User notes management",
	Long: `User notes management in GophKeeper. For example:
	- client notes get noteID
	- client notes getAll
	- client notes add note info
	- client notes remove noteID`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var addNoteCmd = &cobra.Command{
	Use:   "add [flags]",
	Short: "Add a new user note to GophKeeper",
	Long: `This command allows you to add a new user note to your account in GophKeeper. For example:
	- client notes add --text NoteText --desc Description`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.AddNote(&note); err != nil {
			fmt.Println(err)
		}
	},
}

var getNoteCmd = &cobra.Command{
	Use:   "get [flags]",
	Short: "Get user note by ID from GophKeeper",
	Long: `This command allows you to get user note by ID from your account in GophKeeper. For example:
	- client notes get --id 9`,
	Run: func(cmd *cobra.Command, args []string) {
		if noteID < 0 {
			fmt.Println("You must provide a note ID")
			os.Exit(1)
		}
		if err := client.GetNote(noteID); err != nil {
			fmt.Println(err)
		}
	},
}

var removeNoteCmd = &cobra.Command{
	Use:   "remove [flags]",
	Short: "Remove user note by ID from GophKeeper",
	Long: `This command allows you to remove user note by ID from your account in GophKeeper. For example:
	- client notes remove --id 9`,
	Run: func(cmd *cobra.Command, args []string) {
		if noteID < 0 {
			fmt.Println("You must provide a note ID")
			os.Exit(1)
		}
		if err := client.RemoveNote(noteID); err != nil {
			fmt.Println(err)
		}
	},
}

var getAllNotesCmd = &cobra.Command{
	Use:   "getAll",
	Short: "Get all user notes from GophKeeper",
	Long: `This command allows you to get all user notes from your account in GophKeeper. For example:
	- client notes getAll`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.GetAllNotes(); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	addNoteCmd.PersistentFlags().StringVar(&note.Text, "text", "", "text")
	addNoteCmd.PersistentFlags().StringVar(&note.Description, "desc", "", "note description")

	getNoteCmd.PersistentFlags().Int64Var(&noteID, "id", -1, "note id")

	removeNoteCmd.PersistentFlags().Int64Var(&noteID, "id", -1, "note id")

	notesCmd.AddCommand(getNoteCmd)
	notesCmd.AddCommand(removeNoteCmd)
	notesCmd.AddCommand(addNoteCmd)
	notesCmd.AddCommand(getAllNotesCmd)
	rootCmd.AddCommand(notesCmd)
}
