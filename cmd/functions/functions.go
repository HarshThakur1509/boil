package functions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func InsertCode(path string, content string) error {

	existingContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Append the new content to the existing content
	updatedContent := string(existingContent) + content

	// Write the updated content back to the file
	if err := os.WriteFile(path, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func ReplaceCode(path string, code string, replace string) error {

	// Read the file content
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, replace) {
		return fmt.Errorf("the file does not contain the '%s' comment", replace)
	}

	// Replace the comment with the provided code
	updatedContent := strings.Replace(contentStr, replace, code, 1)

	// Write the updated content back to the file
	if err := os.WriteFile(path, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated content to file: %w", err)
	}

	return nil
}

func DeleteCode(filePath string, codeToDelete string) error {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to access file: %w", err)
	}

	contentStr := string(content)

	// Check if code exists in file
	if !strings.Contains(contentStr, codeToDelete) {
		return fmt.Errorf("code pattern not found in %s", filepath.Base(filePath))
	}

	// Remove all instances of the code
	modifiedContent := strings.ReplaceAll(contentStr, codeToDelete, "")

	// Preserve line endings
	if modifiedContent == contentStr {
		return fmt.Errorf("no changes made - code pattern not found")
	}

	// Write modified content back
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

	// srcFolder := filepath.Join(tempDir, folder)
	// if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
	// 	return fmt.Errorf("folder %q not found in repository", folder)
	// }

	// return moveContents(srcFolder, parentDir)
	srcFolder := filepath.Join(tempDir, folder)
	if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
		return fmt.Errorf("folder %q not found in repository", folder)
	}

	return moveContents(srcFolder, parentDir)
}

func writeSparseConfig(tempDir, folder string) error {
	configPath := filepath.Join(tempDir, ".git", "info", "sparse-checkout")
	content := "/*\n!" + folder + "/*\n" + folder + "/**"
	return os.WriteFile(configPath, []byte(content), 0644)
}

func moveContents(src, dest string) error {
	// Create destination directory if needed
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("create destination dir failed: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("read dir failed: %w", err)
	}

	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}

		oldPath := filepath.Join(src, entry.Name())
		newPath := filepath.Join(dest, entry.Name())
		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("move %q failed: %w", entry.Name(), err)
		}
	}
	return nil
}

func ReadRepoFile(repoURL, filePath string) ([]byte, error) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "boil-read-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize git repository
	initCmd := exec.Command("git", "init")
	initCmd.Dir = tempDir
	if output, err := initCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git init failed: %s\n%w", output, err)
	}

	// Enable sparse checkout
	configCmd := exec.Command("git", "config", "core.sparseCheckout", "true")
	configCmd.Dir = tempDir
	if output, err := configCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("sparse config failed: %s\n%w", output, err)
	}

	// Write sparse checkout pattern
	sparsePath := filepath.Join(tempDir, ".git", "info", "sparse-checkout")
	if err := os.MkdirAll(filepath.Dir(sparsePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create sparse dir: %w", err)
	}
	if err := os.WriteFile(sparsePath, []byte(filePath+"\n"), 0644); err != nil {
		return nil, fmt.Errorf("failed to write sparse config: %w", err)
	}

	// Add remote
	remoteCmd := exec.Command("git", "remote", "add", "origin", repoURL)
	remoteCmd.Dir = tempDir
	if output, err := remoteCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("remote add failed: %s\n%w", output, err)
	}

	// Fetch and checkout
	pullCmd := exec.Command("git", "pull", "origin", "main", "--depth=1")
	pullCmd.Dir = tempDir
	if output, err := pullCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git pull failed: %s\n%w", output, err)
	}

	// Verify and read file
	targetFile := filepath.Join(tempDir, filePath)
	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %q not found in repository", filePath)
	}

	return os.ReadFile(targetFile)
}
