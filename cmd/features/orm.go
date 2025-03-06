/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package features

import (
	"log"
	"os"
	"path/filepath"

	"github.com/HarshThakur1509/boil/cmd/functions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ormCmd represents the orm command
var ormCmd = &cobra.Command{
	Use:   "orm",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := "https://github.com/HarshThakur1509/boilerplate"
		name := viper.GetString("name")
		orm, _ := cmd.Flags().GetString("name")
		repoFolder := filepath.Join("features", "orm", orm)
		cwd := viper.GetString("path")

		if !viper.IsSet("orm") {
			viper.Set("orm", orm)
		}

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
		tempDir, err := os.MkdirTemp(cwd, "boil-orm-*")
		if err != nil {
			log.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Clone and move files
		if err := functions.CloneRepo(repoURL, tempDir, repoFolder, cwd); err != nil {
			log.Fatalf("Orm setup failed: %v", err)
		}

		switch orm {
		case "gorm":
			functions.ReplaceCode(filepath.Join(cwd, "migrations", "migrate.go"), name, "myapp")
		case "sqlc":
			functions.ReplaceCode(filepath.Join(cwd, "migrations", "migrate.go"), name, "myapp")
		default:
			log.Fatalf("Invalid orm: %s", orm)

		}

	},
}

func init() {
	FeaturesCmd.AddCommand(ormCmd)
	ormCmd.Flags().StringP("name", "n", "", "Name of the orm")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ormCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ormCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
