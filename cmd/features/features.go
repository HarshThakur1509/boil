/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package features

import (
	"fmt"

	"github.com/spf13/cobra"
)

// featuresCmd represents the features command
var FeaturesCmd = &cobra.Command{
	Use:   "features",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("features called")

	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// featuresCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// featuresCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
