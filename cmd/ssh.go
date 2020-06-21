package cmd

import (
	"fmt"
	"os"
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
		//nolint:dupl // deduplicating it wouldn't provide much value
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Config bindings
			for flag, lookup := range map[string]string{
				"target":       "target",
				"type":         "type",
				"pub-key":      "pub-key",
				"os-user":      "os-user",
				"gen-key-pair": "gen-key-pair",
				"gen-key-dir":  "gen-key-dir",
			} {
				if err := cfg.BindPFlag(flag, cmd.Flags().Lookup(lookup)); err != nil {
					log.WithFields(log.Fields{
						"flag":   flag,
						"lookup": lookup,
					}).Error(err)
					return err
				}
			}
			// returns err
			return aws.VerifyDependencies()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			target := cfg.GetString("target")
			targetType := cfg.GetString("type")
			profile := cfg.GetString("profile")
			region := cfg.GetString("region")
			pubKey := cfg.GetString("pub-key")
			OSUser := cfg.GetString("os-user")
			genKeyPair := cfg.GetBool("gen-key-pair")
			genKeyDir := cfg.GetString("gen-key-dir")
			mfaToken := cfg.GetString("mfa")
			if genKeyPair {
				stat, err := os.Stat(genKeyDir)
				if !(err == nil && stat.IsDir()) {
					if err = os.MkdirAll(genKeyDir, 0750); err != nil {
						return err
					}
				}
				if err != nil {
					err = fmt.Errorf("failed creating directory for temporary keys: %e", err)
					log.WithFields(log.Fields{
						"genKeyDir": genKeyDir,
					}).Error(err)
					return err
				}
				pubKey = path.Join(genKeyDir, tempKeyName+".pub")
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
				"gen-key-dir":  genKeyDir,
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
				MFAToken:   &mfaToken,
			}
			// returns err
			return ssh.Start(input)
		},
		DisableAutoGenTag: true,
	}
)

func init() {
	rootCmd.AddCommand(sshCmd)

	sshCmd.Flags().String("target", "", "specify the target depending on the type")
	sshCmd.Flags().String("type", aws.TargetTypeInstanceID, fmt.Sprintf("specify target type: %s/%s/%s/%s (deprecated)", aws.TargetTypeInstanceID, aws.TargetTypePrivateDNS, aws.TargetTypeName, aws.DeprecatedTargetTypeName))
	sshCmd.Flags().Bool("gen-key-pair", false, fmt.Sprintf("generate a temporary key pair that will be send and used. By default use %s as an identity file", path.Join(workDir, tempKeyName)))
	sshCmd.Flags().String("gen-key-dir", workDir, "the directory where temporary keys will be generated")
	sshCmd.Flags().String("os-user", "ec2-user", "specify an instance OS user which will be using sent public key")
	sshCmd.Flags().String("pub-key", "", "local public key that will be send to the instance, ignored when gen-key-pair is true")
	sshCmd.Flags().Uint64Var(&portNum, "port", portNum, "specify ssh port")
}
