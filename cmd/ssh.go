package cmd

import (
	"fmt"
	"path"

	"github.com/danmx/sigil/pkg/ssh"
	"github.com/danmx/sigil/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const tempKeyName = "temp_key"

var (
	portNum int
	sshCmd  = &cobra.Command{
		Use:   "ssh",
		Short: "Start ssh session",
		Long:  `Start a new ssh for chosen EC2 instance.`,
		//Example: fmt.Sprintf("%s session --type instance-id --target i-xxxxxxxxxxxxxxxxx", AppName),
		PreRun: func(cmd *cobra.Command, args []string) {
			// Config bindings
			if err := cfg.BindPFlag("target", cmd.Flags().Lookup("target")); err != nil {
				log.Fatal(err)
			}
			if err := cfg.BindPFlag("type", cmd.Flags().Lookup("type")); err != nil {
				log.Fatal(err)
			}
			if err := cfg.BindPFlag("pub-key", cmd.Flags().Lookup("pub-key")); err != nil {
				log.Fatal(err)
			}
			if err := cfg.BindPFlag("os-user", cmd.Flags().Lookup("os-user")); err != nil {
				log.Fatal(err)
			}
			if err := cfg.BindPFlag("gen-key-pair", cmd.Flags().Lookup("gen-key-pair")); err != nil {
				log.Fatal(err)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			target := cfg.GetString("target")
			targetType := cfg.GetString("type")
			awsProfile := cfg.GetString("profile")
			awsRegion := cfg.GetString("region")
			pubKey := cfg.GetString("pub-key")
			OSUser := cfg.GetString("os-user")
			genKeyPair := cfg.GetBool("gen-key-pair")
			if genKeyPair {
				pubKey = path.Join(workDir, tempKeyName+".pub")
			}
			log.WithFields(log.Fields{
				"target":       target,
				"type":         targetType,
				"region":       awsRegion,
				"profile":      awsProfile,
				"mfa":          awsMFAToken,
				"pub-key":      pubKey,
				"port":         portNum,
				"os-user":      OSUser,
				"gen-key-pair": genKeyPair,
			}).Debug("ssh inputs")
			input := &ssh.StartInput{
				Target:     &target,
				TargetType: &targetType,
				AWSSession: utils.StartAWSSession(awsRegion, awsProfile, awsMFAToken),
				AWSProfile: &awsProfile,
				PortNumber: &portNum,
				PublicKey:  &pubKey,
				OSUser:     &OSUser,
				GenKeyPair: &genKeyPair,
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

	sshCmd.Flags().String("target", "", "specify the target depedning on the type")
	sshCmd.Flags().String("type", "instance-id", "specify target type: instance-id/private-dns/name-tag")
	sshCmd.Flags().Bool("gen-key-pair", false, fmt.Sprintf("generate a temporary key pair that will be send and used. Use %s as an identity file", path.Join("${HOME}", ".sigil", tempKeyName)))
	sshCmd.Flags().String("os-user", "ec2-user", "specify an instance OS user which will be using sent public key")
	sshCmd.Flags().String("pub-key", "", "local public key that will be send to the instance, ignored when gen-key-pair is true")
	sshCmd.Flags().IntVar(&portNum, "port", 22, "specify ssh port")
}
