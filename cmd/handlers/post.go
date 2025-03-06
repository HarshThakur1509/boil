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

// postCmd represents the post command
var postCmd = &cobra.Command{
	Use:   "post",
	Short: "post command adds post controller to the project.",
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

		viper.ReadInConfig()

		// Check if the model exists
		viper.IsSet(fmt.Sprintf("models.%s", model))

		// Get the model details
		modelData := viper.GetStringMap(fmt.Sprintf("models.%s", model))
		fields := functions.WriteMap(modelData)

		handlersPath := filepath.Join(cwd, "internal", "handlers", "handlers.go")

		code := ""
		routesPath := ""
		routesCode := ""

		// Setup Orm
		ormCode := ""
		switch orm {
		case "gorm":

			ormCode = fmt.Sprintf(`%[3]v := models.%[1]v{
				`, capital, fields, model)

			for key := range modelData {
				caser := cases.Title(language.English)
				capitalKey := caser.String(key)
				ormCode += fmt.Sprintf("\t\t%[1]v: body.%[1]v,\n", capitalKey)
			}

			ormCode += fmt.Sprintf(`
						}
						result := initializers.DB.Create(&%[1]v)
				
						if result.Error != nil {
							http.Error(w, "Something went wrong!!", http.StatusBadRequest)
							return
						}`, model)

		case "sqlc":
			ormCode = fmt.Sprintf(`%[2]v, err := db.New(initializers.DB).Post%[1]v(r.Context()`, capital, model)

			for key := range modelData {
				caser := cases.Title(language.English)
				capitalKey := caser.String(key)
				ormCode += fmt.Sprintf(", body.%[1]v", capitalKey)
			}

			ormCode += fmt.Sprintf(`if err != nil {
		http.Error(w, "Failed to create %[1]v", http.StatusInternalServerError)
		return
		}`, model)
		default:
			log.Fatal("Invalid ORM. Use --orm flag")
		}

		switch framework {
		case "standard":

			code = fmt.Sprintf(`
			func Post%[1]v(w http.ResponseWriter, r *http.Request) {
				var body struct {
		%[2]v		}
				json.NewDecoder(r.Body).Decode(&body)
		
			%[4]v
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated) // 201 Created
				json.NewEncoder(w).Encode(%[3]v)
			}
			`, capital, fields, model, ormCode)

			routesPath = filepath.Join(cwd, "internal", "routes", "routes.go")
			routesCode = fmt.Sprintf(`
		router.HandleFunc("POST /%[2]v", handlers.Post%[1]v)
		// Add code here
		
						`, capital, model)
		case "gin":

			code = fmt.Sprintf(`
			func Post%[1]v(c *gin.Context) {
				var body struct {
		%[2]v		}
				c.Bind(&body)
		
				%[4]v

	// Return the created post
	c.JSON(200, gin.H{
		"%[1]v": %[1]v,
	})
			}
			`, capital, fields, model, ormCode)

			routesPath = filepath.Join(cwd, "cmd", "api", "main.go")
			routesCode = fmt.Sprintf(`
r.POST("/%[2]v", handlers.Post%[1]v)
// Add code here
				`, capital, model)
		}

		functions.InsertCode(handlersPath, code)
		functions.ReplaceCode(routesPath, routesCode, "// Add code here")
		fmt.Printf("Post handler added for model: %s\n", model)
	},
}

func init() {
	HandlersCmd.AddCommand(postCmd)

	// Add name flag to each command
	postCmd.Flags().StringP("name", "n", "", "Name of the model")
}
