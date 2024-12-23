/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package controllers

import (
	"fmt"

	"github.com/spf13/cobra"
)

// controllersCmd represents the controllers command
var ControllersCmd = &cobra.Command{
	Use:   "controllers",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("controllers called")
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
