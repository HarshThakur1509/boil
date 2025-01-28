// init.go
package cmd

import (
	"log"
	"os"

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
		if folder != "" {
			viper.Set("folder", folder)
		}

		// Write configuration
		if err := viper.SafeWriteConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
				log.Fatalf("Config error: %v", err)
			}
		}

		// Clone repository
		if err := functions.CloneRepo(repoURL, tempDir, folder, cwd); err != nil {
			log.Fatalf("Initialization failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("folder", "f", "", "Specific folder to download from the repository")
}
