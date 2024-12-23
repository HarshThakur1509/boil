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

// listidCmd represents the listid command
var listidCmd = &cobra.Command{
	Use:   "listid",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model := args[0]
		capital := strings.Title(model)

		controllersPath := fmt.Sprintf("%s\\controllers\\controllers.go", viper.GetString("path"))

		code := fmt.Sprintf(`
func List%[1]vId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var %[2]v models.%[1]v
	initializers.DB.First(&%[2]v, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(%[2]v)
}
		`, capital, model)

		apiPath := fmt.Sprintf("%s\\api\\api.go", viper.GetString("path"))
		apiCode := fmt.Sprintf(`
router.HandleFunc("GET /%[2]v/{id}", controllers.List%[1]vId)
// Add code here
		`, capital, model)

		functions.InsertCode(controllersPath, code)
		functions.ReplaceCode(apiPath, apiCode)
		fmt.Printf("Listing all entities for model: %s\n", model)
	},
}

func init() {
	ControllersCmd.AddCommand(listidCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listidCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listidCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
