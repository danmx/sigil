package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/danmx/sigil/pkg/aws"
	"github.com/danmx/sigil/pkg/list"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	sessionFilters = map[string]string{
		"after":  "",
		"before": "",
		"target": "",
		"owner":  "",
	}
	// listCmd represents the list command
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List available EC2 instances or SSM sessions",
		Long: `Show list of all EC2 instances with AWS SSM Agent running or active SSM sessions.

Supported groups of filters:
- instances:
	- tags - list of tag keys with a list of values for given keys
	- ids - list of instastance ids
- sessions:
	- after - the timestamp, in ISO-8601 Extended format, to see sessions that started after given date
	- before - the timestamp, in ISO-8601 Extended format, to see sessions that started before given date
	- target - an instance to which session connections have been made
	- owner - an AWS user account to see a list of sessions started by that user

Filter format examples:
[default.filters.session]
  after="2018-08-29T00:00:00Z"
  before="2019-08-29T00:00:00Z"
  target="i-xxxxxxxxxxxxxxxx1"
  owner="user@example.com"
[default.filters.instance]
  ids=["i-xxxxxxxxxxxxxxxx1","i-xxxxxxxxxxxxxxxx2"]
  tags=[{key="Name",values=["WebApp1","WebApp2"]}]
`,
		Aliases: []string{"ls", "l"},
		Example: fmt.Sprintf(`%s list --output-format wide --instance-tags '[{"key":"Name","values":["Web","DB"]}]'`, AppName),
		//nolint:dupl
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Config bindings
			if err := cfg.BindPFlag("output-format", cmd.Flags().Lookup("output-format")); err != nil {
				log.Error(err)
				return err
			}
			if err := cfg.BindPFlag("interactive", cmd.Flags().Lookup("interactive")); err != nil {
				log.Error(err)
				return err
			}
			if err := cfg.BindPFlag("filters.session", cmd.Flags().Lookup("session-filters")); err != nil {
				log.Error(err)
				return err
			}
			if err := cfg.BindPFlag("filters.instance.ids", cmd.Flags().Lookup("instance-ids")); err != nil {
				log.Error(err)
				return err
			}
			if err := cfg.BindPFlag("filters.instance.tags", cmd.Flags().Lookup("instance-tags")); err != nil {
				log.Error(err)
				return err
			}
			if err := cfg.BindPFlag("list-type", cmd.Flags().Lookup("type")); err != nil {
				log.Error(err)
				return err
			}
			if err := aws.VerifyDependencies(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var filters aws.Filters
			if err := cfg.UnmarshalKey("filters", &filters); err != nil {
				log.Error("failed unmarshaling filters")
				return err
			}
			outputFormat := cfg.GetString("output-format")
			profile := cfg.GetString("profile")
			region := cfg.GetString("region")
			interactive := cfg.GetBool("interactive")
			listType := cfg.GetString("list-type")
			instanceIDs := cfg.GetStringSlice("filters.instance.ids")
			// hack to get map[string]string from args
			// https://github.com/spf13/viper/issues/608
			if cfg.IsSet("filters.session") {
				filters.Session = aws.SessionFilters{
					After:  sessionFilters["after"],
					Before: sessionFilters["before"],
					Target: sessionFilters["target"],
					Owner:  sessionFilters["owner"],
				}
			}
			if cfg.IsSet("filters.instance.ids") {
				filters.Instance.IDs = instanceIDs
			}
			var tags []aws.TagValues
			if cfg.IsSet("filters.instance.tags") {

				if err := json.Unmarshal([]byte(cfg.GetString("filters.instance.tags")), &tags); err != nil {
					log.WithField("tags", cfg.GetString("filters.instance.tags")).Error("failed unmarshaling tags")
					return err
				}
				filters.Instance.Tags = tags
			}
			log.WithFields(log.Fields{
				"filters":        filters,
				"output-format":  outputFormat,
				"region":         region,
				"profile":        profile,
				"mfa":            mfaToken,
				"interactive":    interactive,
				"type":           listType,
				"instanceIDs":    instanceIDs,
				"sessionFilters": sessionFilters,
				"tags":           tags,
			}).Debug("List inputs")
			input := &list.StartInput{
				OutputFormat: &outputFormat,
				MFAToken:     &mfaToken,
				Region:       &region,
				Profile:      &profile,
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

	listCmd.Flags().String("output-format", list.FormatText, fmt.Sprintf("specify output format: %s/%s/%s/%s", list.FormatText, list.FormatWide, list.FormatJSON, list.FormatYAML))
	listCmd.Flags().BoolP("interactive", "i", false, "pick an instance from a list and start the session")
	listCmd.Flags().StringP("type", "t", list.TypeListInstances, fmt.Sprintf("specify list type: %s/%s", list.TypeListInstances, list.TypeListSessions))
	listCmd.Flags().StringToStringVar(&sessionFilters, "session-filters", sessionFilters, "specify session filters to limit results")
	listCmd.Flags().StringSlice("instance-ids", []string{}, "specify instance ids to limit results")
	listCmd.Flags().String("instance-tags", "", "specify instance tags, in JSON format, to limit results")
}
