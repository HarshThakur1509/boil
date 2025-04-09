// init.go
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/HarshThakur1509/boil/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init downloads the boilerplate code from github repo.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := "https://github.com/HarshThakur1509/boilerplate"
		yaml, _ := cmd.Flags().GetString("yaml")
		if yaml != "" {
			data, err := util.ReadYaml(yaml)
			if err != nil {
				log.Fatalf("Failed to read YAML: %v", err)
			}
			if err := util.YamlExec(data); err != nil {
				log.Fatal(err)
			}
			return
		}
		framework, _ := cmd.Flags().GetString("framework")
		name, _ := cmd.Flags().GetString("name")
		orm, _ := cmd.Flags().GetString("orm")

		if name == "" {
			name = "myapp"
		}

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}

		// Configure Viper settings
		viper.Set("path", cwd)
		if !viper.IsSet("framework") {
			viper.Set("framework", framework)
		}
		if !viper.IsSet("name") {
			viper.Set("name", name)
		}
		if orm != "" {
			viper.Set("orm", orm)
		}

		// Merge with existing config and write
		if err := viper.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				log.Fatalf("Error reading config: %v", err)
			}
		}

		// Force-write the updated config (overwrite existing file)
		if err := viper.WriteConfigAs(filepath.Join(cwd, "boil.yaml")); err != nil {
			log.Fatalf("Failed to write config: %v", err)
		}

		if orm != "" {

			ormRepoFolder := filepath.Join("features", "orm", orm)
			// Create temporary directory
			ormTempDir, err := os.MkdirTemp(cwd, "boil-orm-*")
			if err != nil {
				log.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(ormTempDir)

			// Clone and move files
			if err := util.CloneRepo(repoURL, ormTempDir, ormRepoFolder, cwd); err != nil {
				log.Fatalf("Orm setup failed: %v", err)
			}

			switch orm {
			case "gorm":
				util.ReplaceCode(filepath.Join(cwd, "migrations", "migrate.go"), name, "myapp")
			case "sqlc":
				util.ReplaceCode(filepath.Join(cwd, "migrations", "migrate.go"), name, "myapp")
			default:
				log.Fatalf("Invalid orm: %s", orm)

			}

		}

		if framework != "" {

			// Create temporary directory
			tempDir, err := os.MkdirTemp(cwd, "boil-clone-*")
			if err != nil {
				log.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Clone repository
			repoFolder := filepath.Join("", framework)
			if err := util.CloneRepo(repoURL, tempDir, repoFolder, cwd); err != nil {
				log.Fatalf("Initialization failed: %v", err)
			}

			// Change project name
			util.ReplaceCode(filepath.Join(cwd, "go.mod"), name, "myapp")
			util.ReplaceCode(filepath.Join(cwd, "cmd", "api", "main.go"), name, "myapp")
			util.ReplaceCode(filepath.Join(cwd, "cmd", "api", "main.go"), name, "myapp")
			util.ReplaceCode(filepath.Join(cwd, "internal", "routes", "routes.go"), name, "myapp")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("framework", "f", "", "Specific framework to download from the repository")
	initCmd.Flags().StringP("name", "n", "", "Name of the project")
	initCmd.Flags().StringP("orm", "o", "", "Name of the orm")
	initCmd.Flags().StringP("yaml", "y", "", "Name of the yaml file")
}
