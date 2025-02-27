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

// listidCmd represents the listid command
var listidCmd = &cobra.Command{
	Use:   "listid",
	Short: "listid command adds listid handler to the project.",
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
		if viper.GetString("Framework") == "standard" {

			code = fmt.Sprintf(`
func List%[1]vId(w http.ResponseWriter, r *http.Request) {
id := r.PathValue("id")

var %[2]v models.%[1]v
initializers.DB.First(&%[2]v, id)

w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(%[2]v)
}
	`, capital, model)

			apiPath = fmt.Sprintf("%s\\internal\\routes\\routes.go", viper.GetString("path"))
			apiCode = fmt.Sprintf(`
router.HandleFunc("GET /%[2]v/{id}", handlers.List%[1]vId)
// Add code here
	`, capital, model)
		} else if viper.GetString("Framework") == "gin" {
			code = fmt.Sprintf(`
func List%[1]vId(c *gin.Context) {
	id := c.Param("id")

	var %[2]v models.%[1]v
	initializers.DB.First(&%[2]v, id)

	c.JSON(200, gin.H{
		"%[2]v": %[2]v,
	})
}
	`, capital, model)

			apiPath = fmt.Sprintf("%s\\cmd\\api\\main.go", viper.GetString("path"))
			apiCode = fmt.Sprintf(`
r.GET("/%[2]v/:id", handlers.List%[1]vId)
// Add code here
	`, capital, model)
		}

		functions.InsertCode(controllersPath, code)
		functions.ReplaceCode(apiPath, apiCode, "// Add code here")
		fmt.Printf("Listing all entities for model: %s\n", model)
	},
}

func init() {
	// Remove Args check since we're using flags
	HandlersCmd.AddCommand(listidCmd)

	// Add name flag
	listidCmd.Flags().StringP("name", "n", "", "Name of the model")
}
