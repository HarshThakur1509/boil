/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/HarshThakur1509/boil/cmd/controllers"
	"github.com/HarshThakur1509/boil/cmd/features"
	"github.com/HarshThakur1509/boil/cmd/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "boil",
	Short:   "Boil is CLI tool to create boilerplate code for golang rest api which use go standard library.",
	Version: "1.3.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(controllers.ControllersCmd)
	rootCmd.AddCommand(models.ModelsCmd)
	rootCmd.AddCommand(features.FeaturesCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Use the current working directory
		cwd, err := os.Getwd()
		cobra.CheckErr(err)
		viper.Set("path", cwd)

		// Search config in home directory with name ".boil" (without extension).
		viper.AddConfigPath(cwd)
		viper.SetConfigType("yaml")
		viper.SetConfigName("boil")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
