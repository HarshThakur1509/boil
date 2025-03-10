/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package handlers

import (
	"fmt"

	"github.com/spf13/cobra"
)

// controllersCmd represents the controllers command
var HandlersCmd = &cobra.Command{
	Use:   "handlers",
	Short: "handlers adds controllers to the project.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("handlers called")
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// controllersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// controllersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
