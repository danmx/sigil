package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	"github.com/danmx/sigil/pkg/list"
	"github.com/danmx/sigil/pkg/session"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

type stringMapStringType struct {
	Map map[string]string
}

func (m *stringMapStringType) Set(value string) error {
	tagsMap, err := stringTagsToMap(value)
	if err != nil {
		return err
	}
	m.Map = tagsMap
	return nil
}

func (m *stringMapStringType) String() string {
	list := make([]string, 0, len(m.Map))
	for key, value := range m.Map {
		list = append(list, key+"="+value)
	}
	return strings.Join(list, ",")
}

func stringTagsToMap(value string) (map[string]string, error) {
	tagsMap := make(map[string]string)
	keyValuePairs := strings.Split(value, ",")
	for _, pair := range keyValuePairs {
		splittedPair := strings.Split(pair, "=")
		if len(splittedPair) != 2 {
			log.WithFields(log.Fields{
				"keyValuePairs": keyValuePairs,
				"pair":          pair,
				"splittedPair":  splittedPair,
			}).Error("wrong format of a key-value pair")
			return nil, fmt.Errorf("wrong format of a key-value pair: %s", pair)
		}
		tagsMap[splittedPair[0]] = splittedPair[1]
	}
	log.WithFields(log.Fields{
		"Tags": tagsMap,
	}).Debug("Tags Map")
	return tagsMap, nil
}

var (
	// AppName is a name of this tool (added at compile time)
	AppName string
	// Version is the semantic version (added at compile time)
	Version string
	// Revision is the git commit id (added at compile time)
	Revision string
	// LogLevel level is setting loging level (added at compile time)
	LogLevel string

	workDir       string
	cfgFilePath   string
	awsRegion     string
	target        string
	targetType    string
	outputFormat  string
	configProfile string
	startSession  bool
	cfgFile       = "config.yaml"
	workDirName   = "." + AppName
	pluginName    = "session-manager-plugin"
)

func init() {
	// Set logging
	log.SetReportCaller(true)
	switch LogLevel {
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.PanicLevel)
	}
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	workDir = path.Join(home, workDirName)

	stat, err := os.Stat(workDir)
	if !(err == nil && stat.IsDir()) {
		if err = os.MkdirAll(workDir, 0750); err != nil {
			log.Fatal(err)
		}
	}

	// CLI init config
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "show version",
	}
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%s version %s (build %s)\n", c.App.Name, c.App.Version, Revision)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = AppName
	app.Version = Version
	app.Usage = "AWS SSM Session manager client"
	app.Description = `A tool for establishing a session in EC2 instances with AWS SSM Agent installed`

	rootFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Value:       cfgFile,
			Usage:       "load configuration from YAML `file`",
			Destination: &cfgFile,
		},
		cli.StringFlag{
			Name:        "work-dir, w",
			Value:       workDir,
			Usage:       "`path` to work directory with config files",
			Destination: &workDir,
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "region",
			Usage:       "specify AWS `region`",
			Destination: &awsRegion,
		}),
	}

	listFlags := []cli.Flag{
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "output-format",
			Value:       "text",
			Usage:       "output formats `text`/json/yaml/wide",
			Destination: &outputFormat,
		}),
		cli.GenericFlag{
			Name:  "tags, t",
			Usage: "specify tags to filter out results, e.g.: `key1=value1,key2=value2`",
			Value: &stringMapStringType{},
		},
		cli.BoolFlag{
			Name:        "interactive, i",
			Usage:       "allows to pick an instance and start the session",
			Destination: &startSession,
		},
	}

	sessionFlags := []cli.Flag{
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "target",
			Usage:       "specify the `target` depedning on the type",
			Destination: &target,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "type",
			Usage:       "specify target `type`: instance-id/priv-dns/name-tag",
			Value:       "instance-id",
			Destination: &targetType,
		}),
	}

	app.Commands = []cli.Command{
		{
			Name:        "list",
			Aliases:     []string{"ls", "l"},
			Usage:       "List available EC2 instances",
			Description: `Show list of all EC2 instances with AWS SSM Agent running.`,
			Flags:       listFlags,
			Action: func(c *cli.Context) error {
				tagFilter := c.Generic("tags").(*stringMapStringType)
				log.WithFields(log.Fields{
					"tags":          tagFilter.String(),
					"output-format": outputFormat,
					"region":        awsRegion,
				}).Debug("List inputs")
				input := &list.StartInput{
					OutputFormat: &outputFormat,
					AWSRegion:    &awsRegion,
					TagFilter:    &tagFilter.Map,
					StartSession: &startSession,
				}
				err := list.Start(input)
				if err != nil {
					return err
				}
				return nil
			},
			Before: func(c *cli.Context) error {
				for _, flag := range c.Command.Flags {
					log.WithFields(log.Fields{
						"FlagName": flag.GetName(),
						"IsSet":    c.IsSet(flag.GetName()),
					}).Debug("List: Flags")
				}
				inputSource, err := altsrc.NewYamlSourceFromFile(path.Join(workDir, cfgFile))
				if err != nil {
					log.Error(err)
					return nil
				}

				return altsrc.ApplyInputSourceValues(c, inputSource, c.Command.Flags)
			},
		},
		{
			Name:        "session",
			Aliases:     []string{"sess", "s"},
			Usage:       "Start a session",
			Description: `Start a session in chosen EC2 instance.`,
			Flags:       sessionFlags,
			Action: func(c *cli.Context) error {
				log.WithFields(log.Fields{
					"target": target,
					"type":   targetType,
					"region": awsRegion,
				}).Debug("Session inputs")
				input := &session.StartInput{
					Target:     &target,
					TargetType: &targetType,
					AWSRegion:  &awsRegion,
				}
				err := session.Start(input)
				if err != nil {
					return err
				}
				return nil
			},
			Before: func(c *cli.Context) error {
				for _, flag := range c.Command.Flags {
					log.WithFields(log.Fields{
						"FlagName": flag.GetName(),
						"IsSet":    c.IsSet(flag.GetName()),
					}).Debug("Session: Flags")
				}
				inputSource, err := altsrc.NewYamlSourceFromFile(path.Join(workDir, cfgFile))
				if err != nil {
					log.Error(err)
					return nil
				}

				return altsrc.ApplyInputSourceValues(c, inputSource, c.Command.Flags)
			},
		},
		{
			Name:  "verify",
			Usage: "Verify if all external dependencies are available",
			Description: fmt.Sprintf(`This command will check if %s is installed.
			Plugin documentation: https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html`,
				pluginName),
			Flags: sessionFlags,
			Action: func(c *cli.Context) error {
				o, err := exec.LookPath(pluginName)
				if err != nil {
					return err
				}
				fmt.Printf("%s is installed successfully in %s\n", pluginName, o)
				return nil
			},
		},
	}
	app.Before = func(c *cli.Context) error {
		for _, flag := range c.App.Flags {
			log.WithFields(log.Fields{
				"FlagName": flag.GetName(),
				"IsSet":    c.IsSet(flag.GetName()),
			}).Debug("Root: Flags")
		}
		inputSource, err := altsrc.NewYamlSourceFromFile(path.Join(workDir, cfgFile))
		if err != nil {
			log.Error(err)
			return nil
		}

		return altsrc.ApplyInputSourceValues(c, inputSource, c.App.Flags)
	}
	app.Flags = rootFlags

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
