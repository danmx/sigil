package cmd

import (
	"fmt"

	"github.com/danmx/sigil/pkg/list"
	"github.com/danmx/sigil/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	filters  string
	listType = "instances"

	// listCmd represents the list command
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List available EC2 instances or SSM sessions",
		Long: `Show list of all EC2 instances with AWS SSM Agent running.

Supported groups of filters:
- filters that affect listing instances:
	- tags - list of tag keys with a list of values for given keys
	- instance_ids - list of instastance ids
- filters that affect listing sessions:
	- after - the timestamp, in ISO-8601 Extended format, to see sessions that started after given date
	- before - the timestamp, in ISO-8601 Extended format, to see sessions that started before given date
	- target - an instance to which session connections have been made
	- owner - an AWS user account to see a list of sessions started by that user

Filter format examples:
- Instances filters:
{
	"tags":[{"key":"Name","values":["WebApp1","WebApp2"]}],
	"instance_ids":["i-xxxxxxxxxxxxxxxx1","i-xxxxxxxxxxxxxxxx2"],
}

- Sessions filters:
{
	"after":"2018-08-29T00:00:00Z",
	"before":"2019-08-29T00:00:00Z",
	"target":"i-xxxxxxxxxxxxxxxx1",
	"owner":"user@example.com",
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
			interactive := cfg.GetBool("interactive")
			filters := cfg.GetString("filters")
			log.WithFields(log.Fields{
				"filters":       filters,
				"output-format": outputFormat,
				"region":        awsRegion,
				"profile":       awsProfile,
				"mfa":           awsMFAToken,
				"interactive":   interactive,
				"type":          listType,
			}).Debug("List inputs")
			input := &list.StartInput{
				OutputFormat: &outputFormat,
				AWSSession:   utils.StartAWSSession(awsRegion, awsProfile, awsMFAToken),
				AWSProfile:   &awsProfile,
				Filters:      &filters,
				Interactive:  &interactive,
				Type:         &listType,
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
	listCmd.Flags().StringVarP(&listType, "type", "t", listType, "specify list type: instances/sessions")
}
