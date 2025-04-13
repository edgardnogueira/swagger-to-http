package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// setupHooksCmd sets up the hooks command and its subcommands
func setupHooksCmd() *cobra.Command {
	hooksCmd := &cobra.Command{
		Use:   "hooks",
		Short: "Manage Git hooks integration",
		Long:  `Install, update, and manage Git hooks for automatic HTTP file generation`,
	}

	// Add subcommands
	hooksCmd.AddCommand(setupInstallHooksCmd())
	hooksCmd.AddCommand(setupStatusHooksCmd())
	hooksCmd.AddCommand(setupEnableHooksCmd())
	hooksCmd.AddCommand(setupDisableHooksCmd())

	return hooksCmd
}

// setupInstallHooksCmd creates the 'hooks install' command
func setupInstallHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install Git hooks",
		Long:  `Install Git hooks for Swagger to HTTP integration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return installHooks()
		},
	}

	return cmd
}

// setupStatusHooksCmd creates the 'hooks status' command
func setupStatusHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check Git hooks status",
		Long:  `Check if Git hooks are installed and enabled`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkHooksStatus()
		},
	}

	return cmd
}

// setupEnableHooksCmd creates the 'hooks enable' command
func setupEnableHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable Git hooks",
		Long:  `Enable Git hooks if they are disabled`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return enableHooks()
		},
	}

	return cmd
}

// setupDisableHooksCmd creates the 'hooks disable' command
func setupDisableHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable Git hooks",
		Long:  `Disable Git hooks temporarily without uninstalling them`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return disableHooks()
		},
	}

	return cmd
}

// installHooks installs the Git hooks
func installHooks() error {
	// Check if we're in a Git repository
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf("not a Git repository (or .git directory not found)")
	}

	// Determine which script to run based on OS
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Check if PowerShell is available
		if _, err := exec.LookPath("powershell.exe"); err == nil {
			cmd = exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", "-File", "hooks/install.ps1")
		} else {
			return fmt.Errorf("PowerShell not found, please install hooks manually")
		}
	} else {
		// For Unix-like systems
		installScript := filepath.Join("hooks", "install.sh")
		cmd = exec.Command("sh", installScript)
		
		// Make the script executable
		if err := os.Chmod(installScript, 0755); err != nil {
			fmt.Printf("Warning: Could not make install script executable: %s\n", err)
		}
	}

	// Set up command to use current directory
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the installation script
	fmt.Println("Installing Git hooks...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Git hooks: %w", err)
	}

	fmt.Println("Git hooks installed successfully!")
	return nil
}

// checkHooksStatus checks if Git hooks are installed and enabled
func checkHooksStatus() error {
	// Check if we're in a Git repository
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf("not a Git repository (or .git directory not found)")
	}

	// Check if pre-commit hook is installed
	preCommitInstalled := false
	if _, err := os.Stat(".git/hooks/pre-commit"); err == nil {
		preCommitInstalled = true
	}

	// Check if post-merge hook is installed
	postMergeInstalled := false
	if _, err := os.Stat(".git/hooks/post-merge"); err == nil {
		postMergeInstalled = true
	}

	// Check if hooks are enabled in config
	hooksEnabled := true
	configFile := ".swagger-to-http/hooks.config"
	if _, err := os.Stat(configFile); err == nil {
		content, err := os.ReadFile(configFile)
		if err == nil {
			if contains(string(content), "HOOKS_ENABLED=false") {
				hooksEnabled = false
			}
		}
	}

	// Print status
	fmt.Println("Git hooks status:")
	fmt.Printf("  Pre-commit hook: %s\n", statusString(preCommitInstalled))
	fmt.Printf("  Post-merge hook: %s\n", statusString(postMergeInstalled))
	fmt.Printf("  Hooks enabled: %s\n", statusString(hooksEnabled))

	return nil
}

// enableHooks enables Git hooks by updating the configuration
func enableHooks() error {
	configFile := ".swagger-to-http/hooks.config"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("hooks configuration file not found, please install hooks first")
	}

	// Read the config file
	content, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read hooks configuration: %w", err)
	}

	// Replace the disabled setting with enabled
	newContent := replaceInContent(string(content), "HOOKS_ENABLED=false", "HOOKS_ENABLED=true")

	// Write the updated config back
	if err := os.WriteFile(configFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to update hooks configuration: %w", err)
	}

	fmt.Println("Git hooks have been enabled!")
	return nil
}

// disableHooks disables Git hooks by updating the configuration
func disableHooks() error {
	configFile := ".swagger-to-http/hooks.config"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("hooks configuration file not found, please install hooks first")
	}

	// Read the config file
	content, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read hooks configuration: %w", err)
	}

	// Replace the enabled setting with disabled
	newContent := replaceInContent(string(content), "HOOKS_ENABLED=true", "HOOKS_ENABLED=false")

	// Write the updated config back
	if err := os.WriteFile(configFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to update hooks configuration: %w", err)
	}

	fmt.Println("Git hooks have been disabled. You can re-enable them with 'swagger-to-http hooks enable'")
	return nil
}

// Helper functions

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// statusString returns a string representation of a boolean status
func statusString(status bool) string {
	if status {
		return "Installed/Enabled"
	}
	return "Not installed/Disabled"
}

// replaceInContent replaces oldStr with newStr in content
func replaceInContent(content, oldStr, newStr string) string {
	return strings.Replace(content, oldStr, newStr, -1)
}
