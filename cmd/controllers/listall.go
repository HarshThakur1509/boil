/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package controllers

import (
	"boil/cmd/functions"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listallCmd represents the listall command
var listallCmd = &cobra.Command{
	Use:   "listall",
	Short: "A brief description of your command",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model := args[0]
		// lower := strings.ToLower(model)
		capital := strings.Title(model)

		controllersPath := fmt.Sprintf("%s\\controllers\\controllers.go", viper.GetString("path"))

		code := fmt.Sprintf(`
func List%[1]v(w http.ResponseWriter, r *http.Request) {
	var %[2]v []models.%[1]v
	initializers.DB.Find(&%[2]v)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(%[2]v)
}
		`, capital, model)

		apiPath := fmt.Sprintf("%s\\api\\api.go", viper.GetString("path"))
		apiCode := fmt.Sprintf(`
router.HandleFunc("GET /%[2]v", controllers.List%[1]v)
// Add code here
		`, capital, model)

		functions.InsertCode(controllersPath, code)
		functions.ReplaceCode(apiPath, apiCode)
		fmt.Printf("Listing all entities for model: %s\n", model)
	},
}

func init() {
	ControllersCmd.AddCommand(listallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
