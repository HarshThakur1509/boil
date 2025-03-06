/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package handlers

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/HarshThakur1509/boil/cmd/functions"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete command adds delete handler to the project.",
	Run: func(cmd *cobra.Command, args []string) {
		cwd := viper.GetString("path")
		orm := viper.GetString("orm")
		framework := viper.GetString("framework")
		model, _ := cmd.Flags().GetString("name")
		if model == "" {
			log.Fatal("Model name is required. Use --name flag")
		}
		caser := cases.Title(language.English)
		capital := caser.String(model)
		handlersPath := filepath.Join(cwd, "internal", "handlers", "handlers.go")

		code := ""
		routesPath := ""
		routesCode := ""

		// Setup Orm
		ormCode := ""
		switch orm {
		case "gorm":
			ormCode = fmt.Sprintf(`	var %[2]v []models.%[1]v
				initializers.DB.Delete(&%[2]v, id)
				`, capital, model)
		case "sqlc":
			ormCode = fmt.Sprintf(`%[2]v, err := db.New(initializers.DB).Delete%[1]v(r.Context(), id)
				if err != nil {
				http.Error(w, "Failed to Delete %[1]v", http.StatusInternalServerError)
				return
				}`, capital, model)
		default:
			log.Fatal("Invalid ORM. Use --orm flag")
		}

		// Setup Framework
		switch framework {
		case "standard":

			code = fmt.Sprintf(`
func Delete%[1]v(w http.ResponseWriter, r *http.Request) {
strId := r.PathValue("id")
id, _ := strconv.Atoi(strId)

%[3]v

w.WriteHeader(http.StatusNoContent) // 204 No Content
}
	`, capital, model, ormCode)

			routesPath = filepath.Join(cwd, "internal", "routes", "routes.go")
			routesCode = fmt.Sprintf(`
router.HandleFunc("GET /%[2]v/{id}", handlers.List%[1]vId)
// Add code here
	`, capital, model)

		case "gin":
			code = fmt.Sprintf(`
func Delete%[1]v(c *gin.Context) {
	id := c.Param("id")

%[3]v

	c.Status(200)
}
	`, capital, model, ormCode)

			routesPath = filepath.Join(cwd, "cmd", "api", "main.go")
			routesCode = fmt.Sprintf(`
r.GET("/%[2]v/:id", handlers.List%[1]vId)
// Add code here
	`, capital, model)

		default:
			log.Fatal("Invalid framework. Use --framework flag")
		}

		functions.InsertCode(handlersPath, code)
		functions.ReplaceCode(routesPath, routesCode, "// Add code here")
		fmt.Printf("Delete handler added for model: %s\n", model)
	},
}

func init() {
	HandlersCmd.AddCommand(deleteCmd)

	// Add name flag to each command
	deleteCmd.Flags().StringP("name", "n", "", "Name of the model")
}
