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
	fileName    string
	description string
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
	- client files download --name FileName --dir /path/to/file`,
	Run: func(cmd *cobra.Command, args []string) {
		if fileName == "" || filePath == "" {
			fmt.Println("You must provide a correct file name and dir")
			os.Exit(1)
		}
		if err := client.DownloadFile(fileName, filePath); err != nil {
			fmt.Println(err)
		}
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload [flags]",
	Short: "Upload file to GophKeeper",
	Long: `This command allows you to upload file to your account in GophKeeper. For example:
	- client files upload --path /path/to/file --desc "File description" --text "Text for text files"`,
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			fmt.Println("You must provide a source file path")
			os.Exit(1)
		}

		if err := client.UploadFile(filePath, description); err != nil {
			fmt.Println(err)
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove [flags]",
	Short: "Remove file by name from GophKeeper",
	Long: `This command allows you to remove file by name from your account in GophKeeper. For example:
	- client files remove --name FileName`,
	Run: func(cmd *cobra.Command, args []string) {
		if fileName == "" {
			fmt.Println("You must provide a file name")
			os.Exit(1)
		}
		if err := client.RemoveFile(fileName); err != nil {
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
	downloadCmd.PersistentFlags().StringVar(&fileName, "name", "", "file name to download")
	downloadCmd.PersistentFlags().StringVar(&filePath, "dir", "", "dir where to download file")

	uploadCmd.PersistentFlags().StringVar(&filePath, "path", "", "path to source file")
	uploadCmd.PersistentFlags().StringVar(&description, "desc", "", "file description")

	removeCmd.PersistentFlags().StringVar(&fileName, "name", "", "file name to remove")

	filesCmd.AddCommand(downloadCmd)
	filesCmd.AddCommand(uploadCmd)
	filesCmd.AddCommand(removeCmd)
	filesCmd.AddCommand(getAllCmd)
	rootCmd.AddCommand(filesCmd)
}
