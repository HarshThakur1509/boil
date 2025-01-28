/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package features

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/HarshThakur1509/boil/cmd/functions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		repoURL := "https://github.com/HarshThakur1509/boilerplate"
		repoFolder := "features/auth/standard/"
		cwd := viper.GetString("path")

		if !viper.IsSet("Features") {
			viper.Set("Features", make(map[string]interface{}))
		}

		// Get the features map properly typed
		features := viper.GetStringMap("Features")
		featureName := "Auth"
		features[featureName] = true
		viper.Set("Features", features)

		// Write configuration correctly
		if err := viper.WriteConfig(); err != nil {
			// If config file doesn't exist, create it
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				if err := viper.SafeWriteConfigAs(viper.ConfigFileUsed()); err != nil {
					log.Fatalf("Failed to create config file: %v", err)
				}
			} else {
				log.Fatalf("Failed to update config file: %v", err)
			}
		}

		// Create temporary directory
		tempDir, err := os.MkdirTemp(cwd, "boil-auth-*")
		if err != nil {
			log.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Clone and move files

		// api content
		content, err := functions.ReadRepoFile(repoURL, repoFolder+`api/api.go`)
		if err != nil {
			fmt.Println(err)
		}
		apiContent := fmt.Sprintf(`%v
		// Add code here
		`, string(content))
		functions.ReplaceCode(cwd+`\api\api.go`, apiContent)

		// models content

		fieldMap := make(map[string]any)

		fieldMap["Email"] = "string"
		fieldMap["Password"] = "string"
		fieldMap["Name"] = "string"
		fieldMap["ResetToken"] = "string"
		fieldMap["TokenExpiry"] = "time.Time"

		// Add model to "Models" section
		if !viper.IsSet("Models") {
			viper.Set("Models", make(map[string]interface{}))
		}
		models := viper.GetStringMap("Models")
		modelName := "User"
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

		if err := functions.DeleteCode(cwd+`\migrate\migrate.go`, "initializers.LoadEnv()"); err != nil {
			log.Fatalf("Code deletion failed: %v", err)
		}

		// controllers content
		if err := functions.CloneRepo(repoURL, tempDir, repoFolder+"controllers", cwd+`\controllers`); err != nil {
			log.Fatalf("Docker setup failed: %v", err)
		}

		// main content
		content, err = functions.ReadRepoFile(repoURL, repoFolder+`main.go`)
		if err != nil {
			fmt.Println(err)
		}
		functions.ReplaceCode(cwd+`\main.go`, string(content))

		if err := functions.DeleteCode(cwd+`\main.go`, "initializers.LoadEnv()"); err != nil {
			log.Fatalf("Code deletion failed: %v", err)
		}

	},
}

func init() {
	FeaturesCmd.AddCommand(authCmd)

}
