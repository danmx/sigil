package cmd

import (
	"fmt"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// AppName is a name of this tool (added at compile time)
	AppName string
	// Version is the semantic version (added at compile time)
	Version string
	// Revision is the git commit id (added at compile time)
	Revision string
	// LogLevel level is setting loging level (added at compile time)
	LogLevel string
	// Debug is turning a debug mode (added at compile time)
	Debug string

	workDir     string
	cfgFile     string
	awsMFAToken string
	cfg         *viper.Viper

	cfgFileName = "config"
	cfgProfile  = "default"
	cfgType     = "toml"
	workDirName = "." + AppName

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:               AppName,
		Short:             "AWS SSM Session manager client",
		Long:              `A tool for establishing a session in EC2 instances with AWS SSM Agent installed`,
		Version:           fmt.Sprintf("%s version %s (build %s)\n", AppName, Version, Revision),
		DisableAutoGenTag: true,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Set debug
	if Debug == "true" {
		log.SetReportCaller(true)
	}
	// Set startup Log level
	if err := setLogLevel(LogLevel); err != nil {
		log.WithFields(log.Fields{
			"LogLevel": LogLevel,
		}).Fatal(err)
	}

	cobra.OnInitialize(initConfig)

	// Config file
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFile, "full config file path")
	rootCmd.PersistentFlags().StringVar(&cfgType, "config-type", cfgType, "specify the type of a config file: json, yaml, toml, hcl, props")
	rootCmd.PersistentFlags().StringVarP(&cfgProfile, "config-profile", "p", cfgProfile, "pick the config profile")
	// Log level
	rootCmd.PersistentFlags().StringVar(&LogLevel, "log-level", LogLevel, "specify the log level: trace/debug/info/warn/error/fatal/panic")
	// AWS
	rootCmd.PersistentFlags().StringP("region", "r", "", "specify AWS region")
	rootCmd.PersistentFlags().String("profile", "", "specify AWS profile")
	rootCmd.PersistentFlags().StringVarP(&awsMFAToken, "mfa", "m", awsMFAToken, "specify MFA token")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	// Set Log level
	if err = setLogLevel(LogLevel); err != nil {
		fmt.Fprintln(os.Stderr, err)
		log.WithFields(log.Fields{
			"LogLevel": LogLevel,
		}).Fatal(err)
	}
	cfg = viper.GetViper()
	if cfgFile != "" {
		// Use config file from the flag.
		cfg.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			log.Fatal(err)
		}
		workDir = path.Join(home, workDirName)
		stat, err := os.Stat(workDir)
		if !(err == nil && stat.IsDir()) {
			if err = os.MkdirAll(workDir, 0750); err != nil {
				fmt.Fprintln(os.Stderr, err)
				log.Fatal(err)
			}
		}
		// Search config in home directory with name from cfgFileName (without extension).
		cfg.AddConfigPath(workDir)
		cfg.SetConfigName(cfgFileName)
		cfg.SetConfigType(cfgType)
	}

	// If a config file is found, read it in.
	if err = cfg.ReadInConfig(); err == nil {
		log.WithFields(log.Fields{
			"config": cfg.ConfigFileUsed(),
		}).Debug("Using config file")
		cfg, err = safeSub(cfg, cfgProfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			log.WithFields(log.Fields{
				"config-profile": cfgProfile,
			}).Fatal(err)
		}
	}

	// Config bindings
	if err := cfg.BindPFlag("region", rootCmd.PersistentFlags().Lookup("region")); err != nil {
		log.Fatal(err)
	}
	if err := cfg.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile")); err != nil {
		log.Fatal(err)
	}
}

// because of https://github.com/spf13/viper/issues/616
func safeSub(v *viper.Viper, profile string) (*viper.Viper, error) {
	subConfig := v.Sub(profile)
	if subConfig == nil {
		return nil, fmt.Errorf("Config profile doesn't exist. Profile: %s", profile)
	}
	return subConfig, nil
}

func setLogLevel(level string) error {
	// Log level
	switch LogLevel {
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
		return fmt.Errorf("Unsupported log level: %s", LogLevel)
	}
	return nil
}
