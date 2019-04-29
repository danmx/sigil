package cmd

import (
	"fmt"

	"github.com/danmx/sigil/pkg/list"
	"github.com/danmx/sigil/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	filters string

	// listCmd represents the list command
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List available EC2 instances",
		Long: `Show list of all EC2 instances with AWS SSM Agent running.

Supported filters:
- tags
- instance_ids

Filter format example:
{
	"tags":[{"key":"Name","values":["WebApp1","WebApp2"]}],
	"instance_id":["i-xxxxxxxxxxxxxxxx1","i-xxxxxxxxxxxxxxxx2"]
}`,
		Aliases: []string{"ls", "l"},
		Example: fmt.Sprintf("%s list --output-format wide --filters=\"{\\\"tags\\\":[{\\\"key\\\":\\\"Name\\\",\\\"values\\\":[\\\"WebApp\\\"]}]}\"", AppName),
		PreRun: func(cmd *cobra.Command, args []string) {
			// Config bindings
			if err := cfg.BindPFlag("output-format", cmd.Flags().Lookup("output-format")); err != nil {
				log.Fatal(err)
			}
			if err := cfg.BindPFlag("interactive", cmd.Flags().Lookup("interactive")); err != nil {
				log.Fatal(err)
			}
			if err := cfg.BindPFlag("filters", cmd.Flags().Lookup("filters")); err != nil {
				log.Fatal(err)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFormat := cfg.GetString("output-format")
			awsProfile := cfg.GetString("profile")
			awsRegion := cfg.GetString("region")
			startSession := cfg.GetBool("interactive")
			filters := cfg.GetString("filters")
			log.WithFields(log.Fields{
				"filters":       filters,
				"output-format": outputFormat,
				"region":        awsRegion,
				"profile":       awsProfile,
				"mfa":           awsMFAToken,
				"interactive":   startSession,
			}).Debug("List inputs")
			input := &list.StartInput{
				OutputFormat: &outputFormat,
				AWSSession:   utils.StartAWSSession(awsRegion, awsProfile, awsMFAToken),
				Filters:      &filters,
				StartSession: &startSession,
			}
			err := list.Start(input)
			if err != nil {
				log.Error(err)
				return err
			}
			return nil
		},
		DisableAutoGenTag: true,
	}
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().String("output-format", "text", "specify output format: text/json/yaml/wide")
	listCmd.Flags().BoolP("interactive", "i", false, "pick an instance from a list and start the session")
	listCmd.Flags().String("filters", "", "specify filters, in JSON format, to limit results")
}
