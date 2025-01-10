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

// initGitRepo initializes a new git repository in the target directory
func initGitRepo(targetDir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = targetDir
	return cmd.Run()
}

// setupRemote adds the remote repository URL
func setupRemote(targetDir, repoURL string) error {
	cmd := exec.Command("git", "remote", "add", "origin", repoURL)
	cmd.Dir = targetDir
	return cmd.Run()
}

// setupSparseCheckout enables and configures sparse checkout
func setupSparseCheckout(targetDir, folder string) error {
	// Enable sparse checkout
	cmd := exec.Command("git", "config", "core.sparseCheckout", "true")
	cmd.Dir = targetDir
	if err := cmd.Run(); err != nil {
		return err
	}

	// Write sparse checkout pattern
	sparseCheckoutPath := filepath.Join(targetDir, ".git", "info", "sparse-checkout")
	if err := os.MkdirAll(filepath.Dir(sparseCheckoutPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(sparseCheckoutPath, []byte(folder+"/*"), 0644)
}

// pullRepo pulls the repository from remote
func pullRepo(targetDir string) error {
	cmd := exec.Command("git", "pull", "origin", "main")
	cmd.Dir = targetDir
	return cmd.Run()
}

// moveContents moves files from source to destination directory
func moveContents(srcDir, destDir string) error {
	contents, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range contents {
		if entry.Name() == ".git" {
			continue
		}
		src := filepath.Join(srcDir, entry.Name())
		dest := filepath.Join(destDir, entry.Name())
		if err := os.Rename(src, dest); err != nil {
			return err
		}
	}
	return nil
}

// cloneFullRepo clones the entire repository
func cloneFullRepo(repoURL, targetDir, parentDir string) error {
	cmd := exec.Command("git", "clone", repoURL, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return moveContents(targetDir, parentDir)
}

// cloneSparseRepo clones specific folder using sparse checkout
func cloneSparseRepo(repoURL, targetDir, parentDir, folder string) error {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	steps := []struct {
		name string
		fn   func() error
	}{
		{"init git repo", func() error { return initGitRepo(targetDir) }},
		{"setup remote", func() error { return setupRemote(targetDir, repoURL) }},
		{"setup sparse checkout", func() error { return setupSparseCheckout(targetDir, folder) }},
		{"pull repo", func() error { return pullRepo(targetDir) }},
		{"move contents", func() error { return moveContents(filepath.Join(targetDir, folder), parentDir) }},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			return fmt.Errorf("failed to %s: %w", step.name, err)
		}
	}

	return nil
}

// CloneRepo clones either the entire repository or a specific folder
func CloneRepo(repoURL, targetDir, folder string) error {
	parentDir := filepath.Dir(targetDir)
	var err error

	if folder == "gin" || folder == "standard" {
		err = cloneSparseRepo(repoURL, targetDir, parentDir, folder)
	} else {
		// err = cloneFullRepo(repoURL, targetDir, parentDir)
		fmt.Println("Please use provided flag for the folder name")
	}

	if err != nil {
		return err
	}

	// Clean up temp directory
	return os.RemoveAll(targetDir)
}
