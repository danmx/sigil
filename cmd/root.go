package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/danmx/sigil/pkg/aws"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	appName string = "sigil"
	// appVersion is the semantic appVersion (added at compile time)
	appVersion string
	// gitCommit is the git commit id (added at compile time)
	gitCommit string
	// logLevel level is setting loging level (added at compile time)
	logLevel string = "panic"
	// dev is turning a debug mode (added at compile time)
	dev string = "false"

	workDir string
	cfg     *viper.Viper

	cfgFileName = "config"
	cfgType     = "toml"
	workDirName = "." + appName

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:               appName,
		Short:             "AWS SSM Session manager client",
		Long:              `A tool for establishing a session in EC2 instances with AWS SSM Agent installed`,
		Version:           fmt.Sprintf("%s (build %s)", appVersion, gitCommit),
		DisableAutoGenTag: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := aws.AppendUserAgent(appName + "/" + appVersion); err != nil {
				return err
			}
			return nil
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	// Set debug
	if dev == "true" {
		log.SetReportCaller(true)
	}
	// Set startup Log level
	if err := setLogLevel(logLevel); err != nil {
		log.WithFields(log.Fields{
			"logLevel": logLevel,
		}).Fatal(err)
	}
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		log.Fatal(err)
	}
	workDir = path.Join(home, workDirName)
	stat, err := os.Stat(workDir)
	if !(err == nil && stat.IsDir()) {
		if err := os.MkdirAll(workDir, 0750); err != nil {
			fmt.Fprintln(os.Stderr, err)
			log.Fatal(err)
		}
	}

	// init config and env vars
	cobra.OnInitialize(func() {
		if err := initConfig(rootCmd); err != nil {
			fmt.Fprintln(os.Stderr, err)
			log.Fatal(err)
		}
	})

	// Config file
	rootCmd.PersistentFlags().StringP("config", "c", "", "full config file path, supported formats: json/yaml/toml/hcl/props")
	rootCmd.PersistentFlags().StringP("config-profile", "p", "default", "pick the config profile")
	// Log level
	rootCmd.PersistentFlags().String("log-level", logLevel, "specify the log level: trace/debug/info/warn/error/fatal/panic")
	// AWS
	rootCmd.PersistentFlags().StringP("region", "r", "", "specify AWS region")
	rootCmd.PersistentFlags().String("profile", "", "specify AWS profile")
	rootCmd.PersistentFlags().StringP("mfa", "m", "", "specify MFA token")
}

// initConfig reads in config file and ENV variables if set
func initConfig(cmd *cobra.Command) error {
	cfg = viper.New()
	// Environment variables
	cfg.SetEnvPrefix(appName)
	cfg.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	// Root config bindings
	for _, key := range []string{"config", "config-profile", "log-level"} {
		if err := cfg.BindEnv(key); err != nil {
			log.WithFields(log.Fields{
				"env": key,
			}).Error(err)
			return err
		}
		if err := cfg.BindPFlag(key, cmd.PersistentFlags().Lookup(key)); err != nil {
			log.WithFields(log.Fields{
				"flag": key,
			}).Error(err)
			return err
		}
	}
	cfgFile := cfg.GetString("config")
	cfgProfile := cfg.GetString("config-profile")
	logLevel := cfg.GetString("log-level")

	// Set Log level
	if err := setLogLevel(logLevel); err != nil {
		fmt.Fprintln(os.Stderr, err)
		log.WithFields(log.Fields{
			"logLevel": logLevel,
		}).Error(err)
		return err
	}
	if cfgFile != "" {
		// Use config file from the flag.
		cfg.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name from cfgFileName (without extension).
		cfg.AddConfigPath(workDir)
		cfg.SetConfigName(cfgFileName)
		cfg.SetConfigType(cfgType)
	}

	// If a config file is found, read it in.
	if err := cfg.ReadInConfig(); err == nil {
		log.WithFields(log.Fields{
			"config": cfg.ConfigFileUsed(),
		}).Debug("Using config file")
		cfg, err = safeSub(cfg, cfgProfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			log.WithFields(log.Fields{
				"config-profile": cfgProfile,
			}).Error(err)
			return err
		}
	}

	// Rebinding config bindings that will be propagated to subcommands because of the subconfig (config profile)
	cfg.SetEnvPrefix(appName)
	if err := cfg.BindEnv("mfa"); err != nil {
		log.WithFields(log.Fields{
			"env": "mfa",
		}).Error(err)
		return err
	}
	for _, key := range []string{"region", "config-profile", "region"} {
		if err := cfg.BindPFlag(key, cmd.PersistentFlags().Lookup(key)); err != nil {
			log.WithFields(log.Fields{
				"flag": key,
			}).Error(err)
			return err
		}
	}
	return nil
}

// because of https://github.com/spf13/viper/issues/616
func safeSub(v *viper.Viper, profile string) (*viper.Viper, error) {
	subConfig := v.Sub(profile)
	if subConfig == nil {
		return nil, fmt.Errorf("config profile doesn't exist. Profile: %s", profile)
	}
	return subConfig, nil
}

// setLogLevel sets the log level
func setLogLevel(level string) error {
	// Log level
	switch level {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	default:
		return fmt.Errorf("unsupported log level: %s", level)
	}
	return nil
}
