// init.go
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/HarshThakur1509/boil/cmd/functions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init downloads the boilerplate code from github repo.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := "https://github.com/HarshThakur1509/boilerplate"
		framework, _ := cmd.Flags().GetString("framework")
		name, _ := cmd.Flags().GetString("name")

		if name == "" {
			name = "myapp"
		}

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}

		// Create temporary directory
		tempDir, err := os.MkdirTemp(cwd, "boil-clone-*")
		if err != nil {
			log.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Configure Viper settings
		viper.Set("path", cwd)
		if !viper.IsSet("framework") {
			viper.Set("framework", framework)
		}
		if !viper.IsSet("name") {
			viper.Set("name", name)
		}

		// Write configuration
		if err := viper.SafeWriteConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
				log.Fatalf("Config error: %v", err)
			}
		}

		// Clone repository
		repoFolder := filepath.Join("", framework)
		if err := functions.CloneRepo(repoURL, tempDir, repoFolder, cwd); err != nil {
			log.Fatalf("Initialization failed: %v", err)
		}

		// Change project name
		functions.ReplaceCode(filepath.Join(cwd, "go.mod"), name, "myapp")
		functions.ReplaceCode(filepath.Join(cwd, "cmd", "api", "main.go"), name, "myapp")
		functions.ReplaceCode(filepath.Join(cwd, "cmd", "api", "main.go"), name, "myapp")
		functions.ReplaceCode(filepath.Join(cwd, "internal", "routes", "routes.go"), name, "myapp")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("framework", "f", "", "Specific framework to download from the repository")
	initCmd.Flags().StringP("name", "n", "", "Name of the project")
}
