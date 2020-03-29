package cmd

import (
	"fmt"

	"github.com/danmx/sigil/pkg/aws"

	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify if all external dependencies are available",
	Long: fmt.Sprint(`This command will check if all dependecies are installed.
Plugin documentation: https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html`),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := aws.VerifyDependencies(); err != nil {
			return err
		}
		fmt.Print("All dependencies are installed\n")
		return nil
	},
	TraverseChildren:      false,
	DisableFlagsInUseLine: true,
	DisableAutoGenTag:     true,
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
