/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package features

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/HarshThakur1509/boil/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Add authentication system to the project",
	Long:  `Sets up authentication system including models, routes, and handlers`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd := viper.GetString("path")
		framework := viper.GetString("framework")
		orm := viper.GetString("orm")

		viper.SetDefault("Features", make(map[string]interface{}))
		viper.SetDefault("Tables", make(map[string]interface{}))

		features := viper.GetStringMap("Features")
		features["Auth"] = true
		viper.Set("Features", features)

		if err := viper.WriteConfig(); err != nil {
			handleConfigError(err)
		}

		tempDir, err := os.MkdirTemp(cwd, "boil-auth-*")
		if err != nil {
			log.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Setup authentication components
		setupRoutes(cwd, framework)
		setupHandlers(cwd, framework, tempDir)
		setupMain(cwd, framework)

		// Create user model
		fieldMap := map[string]any{
			"Email":       "string",
			"Password":    "string",
			"Name":        "string",
			"ResetToken":  "string",
			"TokenExpiry": "time.Time",
		}

		updateConfigTables(fieldMap)

		switch orm {
		case "gorm":
			createModel(fieldMap, "user", cwd)
		case "sqlc":
			createSQLTable(fieldMap, "user", cwd)
		default:
			log.Fatalf("Unsupported ORM: %s", orm)
		}
	},
}

func init() {
	FeaturesCmd.AddCommand(authCmd)
}

func setupRoutes(cwd, framework string) {
	routesPath := filepath.Join("features", "auth", framework, "internal", "routes", "routes.go")
	content, err := util.ReadRepoFile(
		"https://github.com/HarshThakur1509/boilerplate",
		routesPath,
	)
	if err != nil {
		log.Printf("Failed to read routes template: %v", err)
		return
	}

	targetPath := filepath.Join(cwd, "internal", "routes", "routes.go")
	if err := util.ReplaceCode(targetPath, string(content), "// Add code here"); err != nil {
		log.Printf("Failed to update routes: %v", err)
	}
}

func setupHandlers(cwd, framework string, tempDir string) {
	src := filepath.Join("features", "auth", framework, "internal", "handlers")
	dest := filepath.Join(cwd, "internal", "handlers")

	if err := util.CloneRepo(
		"https://github.com/HarshThakur1509/boilerplate",
		tempDir,
		src,
		dest,
	); err != nil {
		log.Fatalf("Handler setup failed: %v", err)
	}
}

func setupMain(cwd, framework string) {
	mainPath := filepath.Join("features", "auth", framework, "cmd", "api", "main.go")
	content, err := util.ReadRepoFile(
		"https://github.com/HarshThakur1509/boilerplate",
		mainPath,
	)
	if err != nil {
		log.Printf("Failed to read main template: %v", err)
		return
	}

	targetPath := filepath.Join(cwd, "cmd", "api", "main.go")
	if err := util.ReplaceCode(targetPath, string(content), "// Add code here"); err != nil {
		log.Printf("Failed to update main: %v", err)
	}
}

func updateConfigTables(fieldMap map[string]any) {
	tables := viper.GetStringMap("Tables")
	tables["user"] = fieldMap
	viper.Set("Tables", tables)

	if err := viper.WriteConfig(); err != nil {
		handleConfigError(err)
	}
}

func createModel(fieldMap map[string]any, tableName, cwd string) {
	structFields := util.WriteMap(fieldMap)
	caser := cases.Title(language.English)
	titleStr := caser.String(tableName)

	modelStruct := fmt.Sprintf("\ntype %s struct {\n\tgorm.Model\n%s}\n", titleStr, structFields)
	modelsPath := filepath.Join(cwd, "internal", "models", "models.go")

	if err := util.InsertCode(modelsPath, modelStruct); err != nil {
		log.Printf("Failed to insert model: %v", err)
	}

	migratePath := filepath.Join(cwd, "migrations", "migrate.go")
	if err := util.ToAutoMigrate(migratePath, titleStr); err != nil {
		log.Printf("Failed to update migrations: %v", err)
	}
}

func createSQLTable(fieldMap map[string]any, tableName, cwd string) {
	fieldSQL := util.GenerateFields(fieldMap)
	tableDef := fmt.Sprintf("CREATE TABLE %s\n(\n\tid SERIAL PRIMARY KEY,\n%s,\n\tcreated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),\n\tupdated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),\n\tdeleted_at TIMESTAMPTZ\n);", tableName, fieldSQL)

	migrationPath := filepath.Join(cwd, "migrations", "db", "migrations", "000001_create_table.up.sql")
	if err := util.InsertCode(migrationPath, tableDef); err != nil {
		log.Printf("Failed to create SQL table: %v", err)
	}
}

func handleConfigError(err error) {
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		if err := viper.SafeWriteConfig(); err != nil {
			log.Fatalf("Failed to create config file: %v", err)
		}
	} else {
		log.Fatalf("Config error: %v", err)
	}
}
