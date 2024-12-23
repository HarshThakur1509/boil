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

// postCmd represents the post command
var postCmd = &cobra.Command{
	Use:   "post",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model := args[0]
		capital := strings.Title(model)

		viper.ReadInConfig()

		// Check if the model exists
		viper.IsSet(fmt.Sprintf("models.%s", model))

		// Get the model details
		modelData := viper.GetStringMap(fmt.Sprintf("models.%s", model))
		fields := functions.WriteMap(modelData)
		// for key, value := range modelData {
		// 	fields += fmt.Sprintf("\t\t%s %s\n", key, value)
		// }

		controllersPath := fmt.Sprintf("%s\\controllers\\controllers.go", viper.GetString("path"))

		code := fmt.Sprintf(`
	func Post%[1]v(w http.ResponseWriter, r *http.Request) {
		var body struct {
%[2]v		}
		json.NewDecoder(r.Body).Decode(&body)

		%[3]v := models.%[1]v{
`, capital, fields, model)

		for key := range modelData {
			code += fmt.Sprintf("\t\t%[1]v: body.%[1]v,\n", strings.Title(key))
		}

		code += fmt.Sprintf(`
		}
		result := initializers.DB.Create(&%[1]v)

		if result.Error != nil {
			http.Error(w, "Something went wrong!!", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201 Created
		json.NewEncoder(w).Encode(%[1]v)
	}
	`, model)

		apiPath := fmt.Sprintf("%s\\api\\api.go", viper.GetString("path"))
		apiCode := fmt.Sprintf(`
router.HandleFunc("POST /%[2]v", controllers.Post%[1]v)
// Add code here
				`, capital, model)

		functions.InsertCode(controllersPath, code)
		functions.ReplaceCode(apiPath, apiCode)
		fmt.Printf("Listing all entities for model: %s\n", model)
	},
}

func init() {
	ControllersCmd.AddCommand(postCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// postCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// postCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
