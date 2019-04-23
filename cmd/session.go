package cmd

import (
	"fmt"

	"github.com/danmx/sigil/pkg/session"
	"github.com/danmx/sigil/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:     "session",
	Short:   "Start a session",
	Long:    `Start a session in chosen EC2 instance.`,
	Aliases: []string{"sess", "s"},
	Example: fmt.Sprintf("%s session --type instance-id --target i-xxxxxxxxxxxxxxxxx", AppName),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Config bindings
		if err := cfg.BindPFlag("target", cmd.Flags().Lookup("target")); err != nil {
			log.Fatal(err)
		}
		if err := cfg.BindPFlag("type", cmd.Flags().Lookup("type")); err != nil {
			log.Fatal(err)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		target := cfg.GetString("target")
		targetType := cfg.GetString("type")
		awsProfile := cfg.GetString("profile")
		awsRegion := cfg.GetString("region")
		log.WithFields(log.Fields{
			"target":  target,
			"type":    targetType,
			"region":  awsRegion,
			"profile": awsProfile,
			"mfa":     awsMFAToken,
		}).Debug("Session inputs")
		input := &session.StartInput{
			Target:     &target,
			TargetType: &targetType,
			AWSSession: utils.StartAWSSession(awsRegion, awsProfile, awsMFAToken),
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

	sessionCmd.Flags().String("target", "", "specify the target depedning on the type")
	sessionCmd.Flags().String("type", "instance-id", "specify target type: instance-id/priv-dns/name-tag")
}
