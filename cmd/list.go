package cmd

import (
	"fmt"

	"github.com/danmx/sigil/pkg/list"
	"github.com/danmx/sigil/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	tags map[string]string

	// listCmd represents the list command
	listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List available EC2 instances",
		Long:    `Show list of all EC2 instances with AWS SSM Agent running.`,
		Aliases: []string{"ls", "l"},
		Example: fmt.Sprintf("%s list --output-format wide -t Name=webapp1", AppName),
		PreRun: func(cmd *cobra.Command, args []string) {
			// Config bindings
			if err := cfg.BindPFlag("output-format", cmd.Flags().Lookup("output-format")); err != nil {
				log.Fatal(err)
			}
			if err := cfg.BindPFlag("interactive", cmd.Flags().Lookup("interactive")); err != nil {
				log.Fatal(err)
			}
			// For now not suppoting tags in config because
			// Viper is lowercasing all keys https://github.com/spf13/viper/pull/635
			//if err := cfg.BindPFlag("tags", cmd.Flags().Lookup("tags")); err != nil {
			//	log.Fatal(err)
			//}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFormat := cfg.GetString("output-format")
			awsProfile := cfg.GetString("profile")
			awsRegion := cfg.GetString("region")
			startSession := cfg.GetBool("interactive")
			log.WithFields(log.Fields{
				"tags":          tags,
				"output-format": outputFormat,
				"region":        awsRegion,
				"profile":       awsProfile,
				"mfa":           awsMFAToken,
				"interactive":   startSession,
			}).Debug("List inputs")
			input := &list.StartInput{
				OutputFormat: &outputFormat,
				AWSSession:   utils.StartAWSSession(awsRegion, awsProfile, awsMFAToken),
				TagFilter:    &tags,
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
	listCmd.Flags().StringToStringVarP(&tags, "tags", "t", tags, "specify tags to filter out results, e.g.: key1=value1,key2=value2")
}
