package functions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func InsertCode(path string, content string) error {
	existingContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	updatedContent := string(existingContent) + content
	if err := os.WriteFile(path, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func ReplaceCode(path string, code string, replace string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, replace) {
		return fmt.Errorf("the file does not contain the '%s' comment", replace)
	}

	updatedContent := strings.Replace(contentStr, replace, code, 1)
	if err := os.WriteFile(path, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated content to file: %w", err)
	}

	return nil
}

func DeleteCode(filePath string, codeToDelete string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to access file: %w", err)
	}

	contentStr := string(content)
	modifiedContent := strings.ReplaceAll(contentStr, codeToDelete, "")

	if modifiedContent == contentStr {
		return fmt.Errorf("no changes made - code pattern not found")
	}

	if err := os.WriteFile(filePath, []byte(modifiedContent), 0644); err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}

	return nil
}

func ToAutoMigrate(filePath, modelName string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	model := fmt.Sprintf("&models.%s{}", modelName)
	re := regexp.MustCompile(`initializers\s*\.\s*DB\s*\.\s*AutoMigrate\s*\(\s*([^)]*)\s*\)`)

	updated := re.ReplaceAllStringFunc(string(content), func(match string) string {
		parts := re.FindStringSubmatch(match)
		existing := strings.ReplaceAll(parts[1], " ", "")

		if existing == "" {
			return fmt.Sprintf("initializers.DB.AutoMigrate(%s)", model)
		}

		if strings.Contains(existing, model) {
			return match
		}

		return fmt.Sprintf(`initializers.DB.AutoMigrate(
		 %s,
		 %s,
		 )`, existing, model)
	})

	return os.WriteFile(filePath, []byte(updated), 0644)
}

func WriteMap(m map[string]any) string {
	var builder strings.Builder
	caser := cases.Title(language.English)

	for k, v := range m {
		builder.WriteString("\t")
		builder.WriteString(caser.String(k))
		builder.WriteString(" ")
		builder.WriteString(v.(string))
		builder.WriteString("\n")
	}
	return builder.String()
}

func GenerateFields(fieldMap map[string]any) string {
	typeMapping := map[string]string{
		"string":  "VARCHAR(255)",
		"uint":    "INTEGER",
		"int":     "INTEGER",
		"float64": "DOUBLE PRECISION",
		"bool":    "BOOLEAN",
	}

	var columns []string
	for fieldName, goType := range fieldMap {
		var sb strings.Builder
		for i, r := range fieldName {
			if unicode.IsUpper(r) && i > 0 {
				sb.WriteRune('_')
			}
			sb.WriteRune(unicode.ToLower(r))
		}
		columnName := sb.String()

		sqlType, ok := typeMapping[goType.(string)]
		if !ok {
			sqlType = "TEXT"
		}

		columns = append(columns, fmt.Sprintf("%s %s NOT NULL", columnName, sqlType))
	}

	return strings.Join(columns, ",\n")
}

func CloneRepo(repoURL, tempDir, repoFolder, targetDir string) error {
	if err := cloneSparseRepo(repoURL, tempDir, targetDir, repoFolder); err != nil {
		return fmt.Errorf("clone failed: %w", err)
	}
	return nil
}

func cloneSparseRepo(repoURL, tempDir, parentDir, folder string) error {
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}

	steps := []struct {
		name string
		cmd  *exec.Cmd
	}{
		{"init", exec.Command("git", "init")},
		{"remote", exec.Command("git", "remote", "add", "origin", repoURL)},
		{"sparse-config", exec.Command("git", "config", "core.sparseCheckout", "true")},
	}

	for _, step := range steps {
		step.cmd.Dir = tempDir
		if output, err := step.cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%s failed: %s\n%w", step.name, string(output), err)
		}
	}

	if err := writeSparseConfig(tempDir, folder); err != nil {
		return err
	}

	pull := exec.Command("git", "pull", "origin", "main")
	pull.Dir = tempDir
	if output, err := pull.CombinedOutput(); err != nil {
		return fmt.Errorf("git pull failed: %s\n%w", string(output), err)
	}

	srcFolder := filepath.Join(tempDir, folder)
	if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
		return fmt.Errorf("folder %q not found in repository", folder)
	}

	targetFolder := filepath.Join(parentDir, filepath.Base(folder))
	if err := os.RemoveAll(targetFolder); err != nil {
		return fmt.Errorf("cleanup failed: %w", err)
	}

	return moveContents(srcFolder, parentDir)
}

func writeSparseConfig(tempDir, folder string) error {
	configPath := filepath.Join(tempDir, ".git", "info", "sparse-checkout")
	// Convert to forward slashes for Git compatibility
	folder = filepath.ToSlash(folder)
	content := fmt.Sprintf("/*\n!%s/*\n%s/**", folder, folder)
	return os.WriteFile(configPath, []byte(content), 0644)
}

func moveContents(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == src {
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		input, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, input, info.Mode())
	})
}

func ReadRepoFile(repoURL, filePath string) ([]byte, error) {
	tempDir, err := os.MkdirTemp("", "boil-read-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	initCmd := exec.Command("git", "init")
	initCmd.Dir = tempDir
	if output, err := initCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git init failed: %s\n%w", output, err)
	}

	configCmd := exec.Command("git", "config", "core.sparseCheckout", "true")
	configCmd.Dir = tempDir
	if output, err := configCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("sparse config failed: %s\n%w", output, err)
	}

	sparsePath := filepath.Join(tempDir, ".git", "info", "sparse-checkout")
	if err := os.MkdirAll(filepath.Dir(sparsePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create sparse dir: %w", err)
	}

	// Ensure correct path format for Git
	gitPath := filepath.ToSlash(filePath)
	if err := os.WriteFile(sparsePath, []byte(gitPath+"\n"), 0644); err != nil {
		return nil, fmt.Errorf("failed to write sparse config: %w", err)
	}

	remoteCmd := exec.Command("git", "remote", "add", "origin", repoURL)
	remoteCmd.Dir = tempDir
	if output, err := remoteCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("remote add failed: %s\n%w", output, err)
	}

	pullCmd := exec.Command("git", "pull", "origin", "main", "--depth=1")
	pullCmd.Dir = tempDir
	if output, err := pullCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git pull failed: %s\n%w", output, err)
	}

	targetFile := filepath.Join(tempDir, filePath)
	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %q not found in repository", filePath)
	}

	return os.ReadFile(targetFile)
}
