/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

// Package cmd contains the commands for the GophKeeper client application.
package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFilePath string
	hashKey     string
	secretKey   string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFilePath, "config", "", "client config file")
	rootCmd.PersistentFlags().StringVar(&hashKey, "hash_key", "", "key for calculate request data hash")
	rootCmd.PersistentFlags().StringVar(&secretKey, "secret_key", "", "key to encrypt data in database")
}

func initConfig() {
	if cfgFilePath != "" {
		viper.SetConfigFile(cfgFilePath)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read config:", err)
			panic("Can't read config")
		}
	} else {
		fmt.Println("You must provide the path to config file")
		panic("You must provide the path to config file")
	}
	if hashKey != "" {
		if err := viper.BindPFlag("hash_key", rootCmd.PersistentFlags().Lookup("hash_key")); err != nil {
			fmt.Println("Can't bind hash key flag to viper: ", err)
			panic("Can't bind hash key flag to viper")
		}
	} else {
		fmt.Println("You must provide the hash_key flag, see --help")
		panic("You must provide the hash_key flag, see --help")
	}

	if secretKey != "" {
		if err := viper.BindPFlag("secret_key", rootCmd.PersistentFlags().Lookup("secret_key")); err != nil {
			fmt.Println("Can't bind secret key flag to viper: ", err)
			panic("Can't bind secret key flag to viper")
		}
	} else {
		fmt.Println("You must provide the secret_key flag, see --help")
		panic("You must provide the secret_key flag, see --help")
	}
}

var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "GophKeeper client application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute is the entry point for client commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
