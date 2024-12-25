/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package controllers

import (
	"fmt"
	"strings"

	"github.com/HarshThakur1509/boil/cmd/functions"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete command adds delete controller to the project.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model := args[0]
		capital := strings.Title(model)

		controllersPath := fmt.Sprintf("%s\\controllers\\controllers.go", viper.GetString("path"))

		code := fmt.Sprintf(`
func Delete%[1]v(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var %[2]v models.%[1]v
	initializers.DB.Delete(&%[2]v, id)
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
		`, capital, model)

		apiPath := fmt.Sprintf("%s\\api\\api.go", viper.GetString("path"))
		apiCode := fmt.Sprintf(`
router.HandleFunc("DELETE /%[2]v/{id}", controllers.Delete%[1]v)
// Add code here
		`, capital, model)

		functions.InsertCode(controllersPath, code)
		functions.ReplaceCode(apiPath, apiCode)
		fmt.Printf("Listing all entities for model: %s\n", model)
	},
}

func init() {
	ControllersCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
