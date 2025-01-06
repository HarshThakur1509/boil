// init.go
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/HarshThakur1509/boil/cmd/functions"
	"github.com/spf13/cobra"
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

		if err := functions.CloneRepo(repoURL, tempDir, folder); err != nil {
			log.Fatalf("Error initializing project: %v", err)
		}

		if folder != "" {
			fmt.Printf("Successfully downloaded folder '%s' from repository into: %s\n", folder, cwd)
		} else {
			fmt.Println("Repository cloned successfully into:", cwd)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("folder", "f", "", "Specific folder to download from the repository")
}
