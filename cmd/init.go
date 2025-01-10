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
		folder, _ := cmd.Flags().GetString("folder")

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
		}

		tempDir := filepath.Join(cwd, "temp-clone")

		// Add model to "Models" section
		if !viper.IsSet("Folder") {
			viper.Set("Folder", folder)
		}

		// Update boil.yaml configuration
		viper.Set("path", cwd)
		if folder != "" {
			viper.Set("folder", folder)
		} else {
			viper.Set("folder", "root")
		}

		// Write configuration to file
		if err := viper.WriteConfig(); err != nil {
			// If the config file doesn't exist, create it
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				if err := viper.SafeWriteConfig(); err != nil {
					log.Fatalf("Failed to create and write to config file: %v", err)
				}
			} else {
				log.Fatalf("Failed to write to config file: %v", err)
			}
		}

		if err := functions.CloneRepo(repoURL, tempDir, folder); err != nil {
			log.Fatalf("Error initializing project: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("folder", "f", "", "Specific folder to download from the repository")
}
