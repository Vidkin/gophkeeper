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

const DefaultConfigPath = "./cfgclient.yaml"

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
			fmt.Println("Can't read config (see --help), error:", err)
			panic("Can't read default config, see --help")
		}
	} else {
		viper.SetConfigFile(DefaultConfigPath)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read default config (see --help), error:", err)
			panic("Can't read default config, see --help")
		}
	}

	if hashKey != "" {
		if err := viper.BindPFlag("hash_key", rootCmd.PersistentFlags().Lookup("hash_key")); err != nil {
			fmt.Println("Can't bind hash key flag to viper: ", err)
			panic("Can't bind hash key flag to viper")
		}
	} else {
		hashKey = viper.GetString("hash_key")
		if hashKey == "" {
			fmt.Println("You must provide the hash_key flag or set it in the config file, see --help")
			panic("You must provide the hash_key flag or set it in the config file, see --help")
		}
	}

	if secretKey != "" {
		if err := viper.BindPFlag("secret_key", rootCmd.PersistentFlags().Lookup("secret_key")); err != nil {
			fmt.Println("Can't bind secret key flag to viper: ", err)
			panic("Can't bind secret key flag to viper")
		}
	} else {
		secretKey = viper.GetString("secret_key")
		if secretKey == "" {
			fmt.Println("You must provide the secret_key flag or set it in the config file, see --help")
			panic("You must provide the secret_key flag or set it in the config file, see --help")
		}
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
