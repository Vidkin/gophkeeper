/*
Copyright © 2024 MIKHAIL SIRKIN <skim991@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Vidkin/gophkeeper/internal/client"
)

var (
	filePath    string
	description string
	text        string
	fileID      int64
)

// filesCmd represents the files management command
var filesCmd = &cobra.Command{
	Use:   "files [command] [flags]",
	Short: "Files management",
	Long: `Files management in GophKeeper. For example:
	- client files download --id fileID --path /path/to/file
	- client files upload --path /path/to/file --desc "File description"
	- client files getAll
	- client files remove fileID`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download [flags]",
	Short: "Download file from GophKeeper",
	Long: `This command allows you to download file from your account in GophKeeper. For example:
	- client files download --id 123 --path /path/to/file`,
	Run: func(cmd *cobra.Command, args []string) {
		if fileID < 0 {
			fmt.Println("You must provide a correct file ID")
			os.Exit(1)
		}
		//if err := client.DownloadFile(fileID, filePath); err != nil {
		//	fmt.Println(err)
		//}
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload [flags]",
	Short: "Upload file to GophKeeper",
	Long: `This command allows you to upload file to your account in GophKeeper. For example:
	- client files upload --path /path/to/file --desc "File description" --text "Text for text files"`,
	Run: func(cmd *cobra.Command, args []string) {
		if filePath != "" && text != "" {
			fmt.Println("You must provide either a source file path or text, not both.")
			os.Exit(1)
		}

		if filePath == "" && text == "" {
			fmt.Println("You must provide a source file path or text")
			os.Exit(1)
		}

		if err := client.UploadFile(filePath, description); err != nil {
			fmt.Println(err)
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove [flags]",
	Short: "Remove file by ID from GophKeeper",
	Long: `This command allows you to remove file by ID from your account in GophKeeper. For example:
	- client files remove --id fileID`,
	Run: func(cmd *cobra.Command, args []string) {
		if fileID < 0 {
			fmt.Println("You must provide a file ID")
			os.Exit(1)
		}
		if err := client.RemoveFile(fileID); err != nil {
			fmt.Println(err)
		}
	},
}

var getAllCmd = &cobra.Command{
	Use:   "getAll",
	Short: "Get all files infos from GophKeeper",
	Long: `This command allows you to get all files infos from your account in GophKeeper. For example:
	- client files getAll`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.GetAllFiles(); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	downloadCmd.PersistentFlags().StringVar(&filePath, "path", "", "path to destination file")
	downloadCmd.PersistentFlags().Int64Var(&fileID, "id", -1, "file id to download")

	uploadCmd.PersistentFlags().StringVar(&filePath, "path", "", "path to source file")
	uploadCmd.PersistentFlags().StringVar(&description, "desc", "", "file description")
	uploadCmd.PersistentFlags().StringVar(&text, "text", "", "text to upload as a text file")

	removeCmd.PersistentFlags().Int64Var(&fileID, "id", -1, "file id to remove")

	filesCmd.AddCommand(downloadCmd)
	filesCmd.AddCommand(uploadCmd)
	filesCmd.AddCommand(removeCmd)
	filesCmd.AddCommand(getAllCmd)
	rootCmd.AddCommand(filesCmd)
}
