/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package models

import (
	"fmt"
	"log"
	"strings"

	"github.com/HarshThakur1509/boil/cmd/functions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// modelsCmd represents the models command
var ModelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Adds Models to the project",
	Long:  ``,
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.GetViper().GetString("port"))
		modelName := args[0]
		fields := args[1:]

		if len(fields)%2 != 0 {
			log.Fatalf("Fields must be provided as name and data type pairs")
		}

		// Build field map from input
		fieldMap := make(map[string]any)
		for i := 0; i < len(fields); i += 2 {
			fieldMap[fields[i]] = fields[i+1]
		}

		// Add model to "Models" section
		if !viper.IsSet("Models") {
			viper.Set("Models", make(map[string]interface{}))
		}
		models := viper.GetStringMap("Models")
		models[modelName] = fieldMap
		viper.Set("Models", models)

		// Write changes back to the file
		if err := viper.WriteConfig(); err != nil {
			// If the config file does not exist, create and write to it
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				if err := viper.SafeWriteConfig(); err != nil {
					log.Fatalf("Failed to create and write to config file: %v", err)
				}
			} else {
				log.Fatalf("Failed to write to config file: %v", err)
			}
		}

		// Append struct definition to models.go
		structFields := functions.WriteMap(fieldMap)
		modelStruct := fmt.Sprintf("\ntype %s struct {\ngorm.Model\n%s}\n", strings.Title(modelName), structFields)
		modelsPath := fmt.Sprintf("%s\\models\\models.go", viper.GetString("path"))

		functions.InsertCode(modelsPath, modelStruct)

		migratePath := fmt.Sprintf("%s\\migrate\\migrate.go", viper.GetString("path"))
		migrateCode := fmt.Sprintf(`
		%[1]v := &models.%[1]v{}
		// Add code here
		`, strings.Title(modelName))

		functions.ReplaceCode(migratePath, migrateCode)
		functions.ToAutoMigrate(migratePath, strings.Title(modelName))

		fmt.Printf("Model '%s' with fields %v has been saved and YAML updated.\n", modelName, fieldMap)
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modelsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modelsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
