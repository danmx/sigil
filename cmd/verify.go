package cmd

import (
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var pluginName = "session-manager-plugin"

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify if all external dependencies are available",
	Long: fmt.Sprintf(`This command will check if %s is installed.
Plugin documentation: https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html`,
		pluginName),
	RunE: func(cmd *cobra.Command, args []string) error {
		o, err := exec.LookPath(pluginName)
		if err != nil {
			log.Error(err)
			return err
		}
		fmt.Printf("%s is installed successfully in %s\n", pluginName, o)
		return nil
	},
	TraverseChildren:      false,
	DisableFlagsInUseLine: true,
	DisableAutoGenTag:     true,
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
