/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package handlers

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/HarshThakur1509/boil/cmd/functions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// listidCmd represents the listid command
var listidCmd = &cobra.Command{
	Use:   "listid",
	Short: "listid command adds listid handler to the project.",
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
				initializers.DB.First(&%[2]v, id)
				`, capital, model)
		case "sqlc":
			ormCode = fmt.Sprintf(`%[2]v, err := db.New(initializers.DB).List%[1]vId(r.Context(), id)
				if err != nil {
				http.Error(w, "Failed to fetch %[1]v", http.StatusInternalServerError)
				return
				}`, capital, model)
		default:
			log.Fatal("Invalid ORM. Use --orm flag")
		}

		switch framework {
		case "standard":

			code = fmt.Sprintf(`
func List%[1]vId(w http.ResponseWriter, r *http.Request) {
strId := r.PathValue("id")
id, _ := strconv.Atoi(strId)

%[3]v

w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(%[2]v)
}
	`, capital, model, ormCode)

			routesPath = filepath.Join(cwd, "internal", "routes", "routes.go")
			routesCode = fmt.Sprintf(`
router.HandleFunc("GET /%[2]v/{id}", handlers.List%[1]vId)
// Add code here
	`, capital, model)

		case "gin":
			code = fmt.Sprintf(`
func List%[1]vId(c *gin.Context) {
	id := c.Param("id")

%[3]v

	c.JSON(200, gin.H{
		"%[2]v": %[2]v,
	})
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
		fmt.Printf("Listing all entities for model: %s\n", model)
	},
}

func init() {
	// Remove Args check since we're using flags
	HandlersCmd.AddCommand(listidCmd)

	// Add name flag
	listidCmd.Flags().StringP("name", "n", "", "Name of the model")
}
