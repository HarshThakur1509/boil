/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package tables

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/HarshThakur1509/boil/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// tableCmd represents the table command
var TablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "Adds Tables to the project",
	Run: func(cmd *cobra.Command, args []string) {
		orm := viper.GetString("orm")
		cwd := viper.GetString("path")
		tableName, _ := cmd.Flags().GetString("name")
		fields := args // Get remaining args after flags

		if tableName == "" {
			log.Fatal("Table name is required. Use --name flag")
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

		// Add model to "Tables" section
		if !viper.IsSet("Tables") {
			viper.Set("Tables", make(map[string]interface{}))
		}
		tables := viper.GetStringMap("Tables")
		tables[tableName] = fieldMap
		viper.Set("Tables", tables)

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

		switch orm {
		case "gorm":
			// Append struct definition to models.go
			structFields := util.WriteMap(fieldMap)

			caser := cases.Title(language.English)
			titleStr := caser.String(tableName)

			modelStruct := fmt.Sprintf("\ntype %s struct {\n	gorm.Model\n	%s}\n", titleStr, structFields)

			util.InsertCode(filepath.Join(cwd, "internal", "models", "models.go"), modelStruct)

			util.ToAutoMigrate(filepath.Join(cwd, "migrations", "migrate.go"), titleStr)

			fmt.Printf("Model '%s' with fields %v has been saved and YAML updated.\n", tableName, fieldMap)
		case "sqlc":
			// Append table definition to migrations/db/migrations/00001_create_table.sql
			fieldSql := util.GenerateFields(fieldMap)
			tableDef := fmt.Sprintf("CREATE TABLE %s\n(id SERIAL PRIMARY KEY,\n%s,\ncreated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),\nupdated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),\ndeleted_at TIMESTAMPTZ)", tableName, fieldSql)

			util.InsertCode(filepath.Join(cwd, "migrations", "db", "migrations", "00001_create_table.sql"), tableDef)

			fmt.Printf("Model '%s' with fields %v has been saved and YAML updated.\n", tableName, fieldMap)

		default:
			log.Fatalf("Unsupported ORM: %s", orm)

		}

	},
}

func init() {
	TablesCmd.Flags().StringP("name", "n", "", "Name of the model")
	TablesCmd.Flags().BoolP("fields", "f", false, "Indicates that following arguments are field definitions")
}
