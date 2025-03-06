/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package features

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/HarshThakur1509/boil/cmd/functions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dockerCmd represents the docker command
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		repoURL := "https://github.com/HarshThakur1509/boilerplate"
		framework := viper.GetString("framework")
		repoFolder := filepath.Join("features", "docker", framework)
		cwd := viper.GetString("path")

		if !viper.IsSet("Features") {
			viper.Set("Features", make(map[string]interface{}))
		}

		// Get the features map properly typed
		features := viper.GetStringMap("Features")
		featureName := "Docker"
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
		tempDir, err := os.MkdirTemp(cwd, "boil-docker-*")
		if err != nil {
			log.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Clone and move files
		if err := functions.CloneRepo(repoURL, tempDir, repoFolder, cwd); err != nil {
			log.Fatalf("Docker setup failed: %v", err)
		}

		// delete initializers.LoadEnv()
		if err := functions.DeleteCode(filepath.Join(cwd, "migrations", "migrate.go"), "initializers.LoadEnv()"); err != nil {
			log.Fatalf("Code deletion failed: %v", err)
		}

		if err := functions.DeleteCode(filepath.Join(cwd, "cmd", "api", "main.go"), "initializers.LoadEnv()"); err != nil {
			log.Fatalf("Code deletion failed: %v", err)
		}

		fmt.Println("✅ Docker files successfully added to project")
	},
}

func init() {
	FeaturesCmd.AddCommand(dockerCmd)

}
