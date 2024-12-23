package functions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func InsertCode(path string, content string) error {
	// Get the absolute path from the relative path
	// path, err := filepath.Abs(path)
	// if err != nil {
	// 	return fmt.Errorf("failed to get absolute path: %w", err)
	// }
	// modelsPath := fmt.Sprintf("%s\\models\\models.go", path)
	// Read the existing content of the file
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

func ReplaceCode(path string, code string) error {

	// Read the file content
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "// Add code here") {
		return fmt.Errorf("the file does not contain the '// Add code here' comment")
	}

	// Replace the comment with the provided code
	updatedContent := strings.Replace(contentStr, "// Add code here", code, 1)

	// Write the updated content back to the file
	if err := os.WriteFile(path, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated content to file: %w", err)
	}

	return nil
}

func ToAutoMigrate(filePath, modelName string) error {
	// Read the file contents
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Convert to string for manipulation
	content := string(fileContent)

	// Define the AutoMigrate line pattern
	autoMigratePattern := `initializers\.DB\.AutoMigrate\(([^)]*)\)`

	// Use regex to find and update the AutoMigrate line
	re := regexp.MustCompile(autoMigratePattern)
	updatedContent := re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract existing models (if any)
		matches := re.FindStringSubmatch(match)
		existingModels := strings.TrimSpace(matches[1])

		// Handle the case where AutoMigrate is empty
		if existingModels == "" {
			return fmt.Sprintf("initializers.DB.AutoMigrate(%s)", modelName)
		}

		// Check if the model is already listed
		if strings.Contains(existingModels, modelName) {
			return match // Model already exists, return as is
		}

		// Append the new model
		return fmt.Sprintf("initializers.DB.AutoMigrate(%s, %s)", existingModels, modelName)
	})

	// Write the updated content back to the file
	if err := os.WriteFile(filePath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated content: %w", err)
	}

	return nil
}

func WriteMap(m map[string]any) string {
	result := ""
	for key, value := range m {
		result += fmt.Sprintf("\t%s %s\n", strings.Title(key), value.(string))
	}
	return result
}

func CloneRepo(repoURL, targetDir string) error {
	// Clone the repository directly into the target directory
	cmd := exec.Command("git", "clone", repoURL, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Move the files from the cloned folder to the root
	contents, err := os.ReadDir(targetDir)
	if err != nil {
		return fmt.Errorf("failed to read cloned directory: %w", err)
	}

	for _, entry := range contents {
		if entry.Name() == ".git" {
			// Skip moving the .git folder
			continue
		}

		srcPath := filepath.Join(targetDir, entry.Name())
		destPath := filepath.Join(targetDir, "..", entry.Name())

		// Rename (move) each file/directory to the parent directory
		if err := os.Rename(srcPath, destPath); err != nil {
			return fmt.Errorf("failed to move %s: %w", entry.Name(), err)
		}
	}

	// Remove the now-empty cloned directory
	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("failed to clean up temporary directory: %w", err)
	}

	return nil
}
