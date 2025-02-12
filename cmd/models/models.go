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
	Run: func(cmd *cobra.Command, args []string) {
		modelName, _ := cmd.Flags().GetString("name")
		fields := args // Get remaining args after flags

		if modelName == "" {
			log.Fatal("Model name is required. Use --name flag")
		}

		if len(fields) == 0 {
			log.Fatal("Fields are required. Use --fields followed by field definitions")
		}

		if len(fields)%2 != 0 {
			log.Fatal("Fields must be provided as name and data type pairs")
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

		// Write configuration
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
		modelStruct := fmt.Sprintf("\ntype %s struct {\n	gorm.Model\n	%s}\n", strings.Title(modelName), structFields)
		modelsPath := fmt.Sprintf("%s\\internal\\models\\models.go", viper.GetString("path"))

		functions.InsertCode(modelsPath, modelStruct)

		migratePath := fmt.Sprintf("%s\\migrations\\migrate.go", viper.GetString("path"))

		functions.ToAutoMigrate(migratePath, strings.Title(modelName))

		fmt.Printf("Model '%s' with fields %v has been saved and YAML updated.\n", modelName, fieldMap)
	},
}

func init() {
	ModelsCmd.Flags().StringP("name", "n", "", "Name of the model")
	ModelsCmd.Flags().BoolP("fields", "f", false, "Indicates that following arguments are field definitions")
}
