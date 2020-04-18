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
	Use:     "session",
	Short:   "Start a session",
	Long:    `Start a new session in chosen EC2 instance.`,
	Aliases: []string{"sess", "s"},
	Example: fmt.Sprintf("%s session --type instance-id --target i-xxxxxxxxxxxxxxxxx", AppName),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Config bindings
		for _, flag := range []string{"target", "type"} {
			if err := cfg.BindPFlag(flag, cmd.Flags().Lookup(flag)); err != nil {
				log.WithFields(log.Fields{
					"flag": flag,
				}).Error(err)
				return err
			}
		}
		if err := aws.VerifyDependencies(); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		target := cfg.GetString("target")
		targetType := cfg.GetString("type")
		profile := cfg.GetString("profile")
		region := cfg.GetString("region")
		mfaToken := cfg.GetString("mfa")
		log.WithFields(log.Fields{
			"target":  target,
			"type":    targetType,
			"region":  region,
			"profile": profile,
			"mfa":     mfaToken,
		}).Debug("Session inputs")
		input := &session.StartInput{
			Target:     &target,
			TargetType: &targetType,
			Region:     &region,
			Profile:    &profile,
			MFAToken:   &mfaToken,
		}
		err := session.Start(input)
		if err != nil {
			return err
		}
		return nil
	},
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(sessionCmd)

	sessionCmd.Flags().String("target", "", "specify the target depending on the type")
	sessionCmd.Flags().String("type", aws.TargetTypeInstanceID, fmt.Sprintf("specify target type: %s/%s/%s", aws.TargetTypeInstanceID, aws.TargetTypePrivateDNS, aws.TargetTypeName))
}
