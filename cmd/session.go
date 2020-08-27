package cmd

import (
	"fmt"

	"github.com/danmx/sigil/pkg/aws"
	"github.com/danmx/sigil/pkg/session"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:                   "session [--type TYPE] ... TARGET",
	DisableFlagsInUseLine: true,
	Short:                 "Start a session",
	Long:                  `Start a new session in chosen EC2 instance.`,
	Aliases:               []string{"sess", "s"},
	Example:               fmt.Sprintf("%s session --type instance-id i-xxxxxxxxxxxxxxxxx", appName),
	Args:                  cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 && cmd.Flags().Lookup("target").Changed {
			return fmt.Errorf("can't use both target argument (%s) and deprecated flag (%s)", args[0], cmd.Flags().Lookup("target").Value.String())
		}
		// Config bindings
		for _, flag := range []string{"target", "type"} {
			if err := cfg.BindPFlag(flag, cmd.Flags().Lookup(flag)); err != nil {
				log.WithFields(log.Fields{
					"flag": flag,
				}).Error(err)
				return err
			}
		}
		if len(args) > 0 {
			cfg.Set("target", args[0])
		}
		// returns err
		return aws.VerifyDependencies()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		target := cfg.GetString("target")
		targetType := cfg.GetString("type")
		profile := cfg.GetString("profile")
		region := cfg.GetString("region")
		mfaToken := cfg.GetString("mfa")
		trace := log.IsLevelEnabled(log.TraceLevel)
		log.WithFields(log.Fields{
			"target":  target,
			"type":    targetType,
			"region":  region,
			"profile": profile,
			"mfa":     mfaToken,
			"trace":   trace,
		}).Debug("Session inputs")
		input := &session.StartInput{
			Target:     &target,
			TargetType: &targetType,
			Region:     &region,
			Profile:    &profile,
			MFAToken:   &mfaToken,
			Trace:      &trace,
		}
		// returns err
		return session.Start(input)
	},
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(sessionCmd)

	sessionCmd.Flags().String("target", "", "specify the target depending on the type")
	// Deprecating the target flag
	err := sessionCmd.Flags().MarkDeprecated("target", "this flag will be deprecated in future releases, use args instead")
	if err != nil {
		log.WithField("flag", sessionCmd.Flags().Lookup("target")).Error(err)
	}
	sessionCmd.Flags().String("type", aws.TargetTypeInstanceID, fmt.Sprintf("specify target type: %s/%s/%s", aws.TargetTypeInstanceID, aws.TargetTypePrivateDNS, aws.TargetTypeName))
}
