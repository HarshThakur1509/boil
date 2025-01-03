/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package controllers

import (
	"fmt"
	"log"
	"strings"

	"github.com/HarshThakur1509/boil/cmd/functions"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listallCmd represents the listall command
var listallCmd = &cobra.Command{
	Use:   "listall",
	Short: "listall command adds listall controller to the project.",
	Run: func(cmd *cobra.Command, args []string) {
		model, _ := cmd.Flags().GetString("name")
		if model == "" {
			log.Fatal("Model name is required. Use --name flag")
		}

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
		fmt.Printf("Listall controller added for model: %s\n", model)
	},
}

func init() {
	ControllersCmd.AddCommand(listallCmd)

	// Add name flag to each command
	listallCmd.Flags().StringP("name", "n", "", "Name of the model")
}
