package cmd

import (
	"fmt"
	"path"

	"github.com/danmx/sigil/pkg/aws"
	"github.com/danmx/sigil/pkg/ssh"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const tempKeyName = "temp_key"

var (
	portNum uint64 = 22
	sshCmd         = &cobra.Command{
		Use:   "ssh",
		Short: "Start ssh session",
		Long:  `Start a new ssh for chosen EC2 instance.`,
		//nolint:dupl
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Config bindings
			if err := cfg.BindPFlag("target", cmd.Flags().Lookup("target")); err != nil {
				log.Error(err)
				return err
			}
			if err := cfg.BindPFlag("type", cmd.Flags().Lookup("type")); err != nil {
				log.Error(err)
				return err
			}
			if err := cfg.BindPFlag("pub-key", cmd.Flags().Lookup("pub-key")); err != nil {
				log.Error(err)
				return err
			}
			if err := cfg.BindPFlag("os-user", cmd.Flags().Lookup("os-user")); err != nil {
				log.Error(err)
				return err
			}
			if err := cfg.BindPFlag("gen-key-pair", cmd.Flags().Lookup("gen-key-pair")); err != nil {
				log.Error(err)
				return err
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
			pubKey := cfg.GetString("pub-key")
			OSUser := cfg.GetString("os-user")
			genKeyPair := cfg.GetBool("gen-key-pair")
			if genKeyPair {
				pubKey = path.Join(workDir, tempKeyName+".pub")
			}
			log.WithFields(log.Fields{
				"target":       target,
				"type":         targetType,
				"region":       region,
				"profile":      profile,
				"mfa":          mfaToken,
				"pub-key":      pubKey,
				"port":         portNum,
				"os-user":      OSUser,
				"gen-key-pair": genKeyPair,
			}).Debug("ssh inputs")
			input := &ssh.StartInput{
				Target:     &target,
				TargetType: &targetType,
				PortNumber: &portNum,
				PublicKey:  &pubKey,
				OSUser:     &OSUser,
				GenKeyPair: &genKeyPair,
				Region:     &region,
				Profile:    &profile,
			}
			err := ssh.Start(input)
			if err != nil {
				return err
			}
			return nil
		},
		DisableAutoGenTag: true,
	}
)

func init() {
	rootCmd.AddCommand(sshCmd)

	sshCmd.Flags().String("target", "", "specify the target depending on the type")
	sshCmd.Flags().String("type", aws.TargetTypeInstanceID, fmt.Sprintf("specify target type: %s/%s/%s", aws.TargetTypeInstanceID, aws.TargetTypePrivateDNS, aws.TargetTypeName))
	sshCmd.Flags().Bool("gen-key-pair", false, fmt.Sprintf("generate a temporary key pair that will be send and used. Use %s as an identity file", path.Join("${HOME}", ".sigil", tempKeyName)))
	sshCmd.Flags().String("os-user", "ec2-user", "specify an instance OS user which will be using sent public key")
	sshCmd.Flags().String("pub-key", "", "local public key that will be send to the instance, ignored when gen-key-pair is true")
	sshCmd.Flags().Uint64Var(&portNum, "port", portNum, "specify ssh port")
}
