package cli

import (
	"fmt"

	"github.com/edgardnogueira/swagger-to-http/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Display the version, build date, and other information about the swagger-to-http tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Info())
	},
}
