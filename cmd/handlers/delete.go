/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/HarshThakur1509/boil/cmd/functions"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete command adds delete handler to the project.",
	Run: func(cmd *cobra.Command, args []string) {
		model, _ := cmd.Flags().GetString("name")
		if model == "" {
			log.Fatal("Model name is required. Use --name flag")
		}

		capital := strings.Title(model)
		controllersPath := fmt.Sprintf("%s\\internal\\handlers\\handlers.go", viper.GetString("path"))

		code := ""
		apiPath := ""
		apiCode := ""
		if viper.GetString("Folder") == "standard" {

			code = fmt.Sprintf(`
	func Delete%[1]v(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
	
		var %[2]v models.%[1]v
		initializers.DB.Delete(&%[2]v, id)
		w.WriteHeader(http.StatusNoContent) // 204 No Content
	}
			`, capital, model)

			apiPath = fmt.Sprintf("%s\\internal\\routes\\routes.go", viper.GetString("path"))
			apiCode = fmt.Sprintf(`
	router.HandleFunc("DELETE /%[2]v/{id}", handlers.Delete%[1]v)
	// Add code here
			`, capital, model)
		} else if viper.GetString("Folder") == "gin" {
			code = fmt.Sprintf(`
func Delete%[1]v(c *gin.Context) {
	id := c.Param("id")

	initializers.DB.Delete(&models.%[1]v{}, id)

	c.Status(200)
}
			`, capital, model)

			apiPath = fmt.Sprintf("%s\\cmd\\api\\main.go", viper.GetString("path"))
			apiCode = fmt.Sprintf(`
r.DELETE("/%[2]v/:id", handlers.Delete%[1]v)
// Add code here
			`, capital, model)
		}

		functions.InsertCode(controllersPath, code)
		functions.ReplaceCode(apiPath, apiCode, "// Add code here")
		fmt.Printf("Delete handler added for model: %s\n", model)
	},
}

func init() {
	HandlersCmd.AddCommand(deleteCmd)

	// Add name flag to each command
	deleteCmd.Flags().StringP("name", "n", "", "Name of the model")
}
