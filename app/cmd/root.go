/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

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
			os.Exit(1)
		}
	} else {
		fmt.Println("You should provide the path to config file")
		os.Exit(1)
	}

	if hashKey != "" {
		if err := viper.BindPFlag("hash_key", rootCmd.PersistentFlags().Lookup("hash_key")); err != nil {
			fmt.Println("Can't bind hash key flag to viper: ", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("You should provide the hash_key flag, see --help")
		os.Exit(1)
	}

	if secretKey != "" {
		if err := viper.BindPFlag("secret_key", rootCmd.PersistentFlags().Lookup("secret_key")); err != nil {
			fmt.Println("Can't bind secret key flag to viper: ", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("You should provide the secret_key flag, see --help")
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "GophKeeper client application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
