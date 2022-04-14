package list

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/danmx/sigil/pkg/aws"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// List wraps methods used from the pkg/list package
type List interface {
	Start(input *StartInput) error
}

const (
	// FormatText points to the text format type
	FormatText = "text"
	// FormatWide points to the wide format type
	FormatWide = "wide"
	// FormatJSON points to the json format type
	FormatJSON = "json"
	// FormatYAML points to the yaml format type
	FormatYAML = "yaml"
	// TypeListInstances points to the instances list type
	TypeListInstances = "instances"
	// TypeListSessions points to the sessions list type
	TypeListSessions = "sessions"
	tabPadding       = 2
)

// StartInput struct contains all input data
type StartInput struct {
	// Define output format
	OutputFormat *string
	Interactive  *bool
	Type         *string
	MFAToken     *string
	Region       *string
	Profile      *string
	Filters      *aws.Filters
	Trace        *bool
}

// Start will output a ist of all available EC2 instances
func Start(input *StartInput) error {
	provider := aws.Provider{}
	err := provider.NewWithConfig(&aws.Config{
		Filters:  *input.Filters,
		Region:   *input.Region,
		Profile:  *input.Profile,
		MFAToken: *input.MFAToken,
		Trace:    *input.Trace,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	switch *input.Type {
	case TypeListInstances:
		err := input.listInstances(&provider)
		if err != nil {
			return err
		}
	case TypeListSessions:
		err := input.listSessions(&provider)
		if err != nil {
			return err
		}
	default:
		err := fmt.Errorf("unsupported list type: %s", *input.Type)
		log.WithField("type", *input.Type).Error(err)
		return err
	}
	return nil
}

func sessionsToString(format string, sessions []*aws.Session) (string, error) {
	switch format {
	case FormatText:
		buf := bytes.NewBufferString("")
		w := new(tabwriter.Writer)
		w.Init(buf, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Index\tSession ID\tTarget\tStart Date")
		for i, session := range sessions {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
				(i + 1), session.SessionID, session.Target, session.StartDate)
		}
		err := w.Flush()
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	case FormatJSON:
		data, err := json.Marshal(sessions)
		if err != nil {
			return "", err
		}
		// JSON output was missing new line
		return string(data) + "\n", nil
	case FormatYAML:
		data, err := yaml.Marshal(sessions)
		if err != nil {
			return "", err
		}
		return string(data), nil
	case FormatWide:
		buf := new(bytes.Buffer)
		w := new(tabwriter.Writer)
		w.Init(buf, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Index\tSession ID\tTarget\tStart Date\tOwner\tStatus")
		for i, session := range sessions {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
				(i + 1), session.SessionID, session.Target, session.StartDate,
				session.Owner, session.Status)
		}
		err := w.Flush()
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
}

func instancesToString(format string, instances []*aws.Instance) (string, error) {
	switch format {
	case FormatText:
		buf := bytes.NewBufferString("")
		w := new(tabwriter.Writer)
		w.Init(buf, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Index\tName\tInstance ID\tIP Address\tPrivate DNS Name")
		for i, instance := range instances {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				(i + 1), instance.Name, instance.ID, instance.IPAddress, instance.PrivateDNSName)
		}
		err := w.Flush()
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	case FormatJSON:
		data, err := json.Marshal(instances)
		if err != nil {
			return "", err
		}
		// JSON output was missing new line
		return string(data) + "\n", nil
	case FormatYAML:
		data, err := yaml.Marshal(instances)
		if err != nil {
			return "", err
		}
		return string(data), nil
	case FormatWide:
		buf := new(bytes.Buffer)
		w := new(tabwriter.Writer)
		w.Init(buf, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Index\tName\tInstance ID\tIP Address\tPrivate DNS Name\tHostname\tOS Name\tOS Version\tOS Type")
		for i, instance := range instances {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				(i + 1), instance.Name, instance.ID, instance.IPAddress, instance.PrivateDNSName,
				instance.Hostname, instance.OSName, instance.OSVersion, instance.OSType)
		}
		err := w.Flush()
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
}

func (input *StartInput) listInstances(provider aws.CloudInstances) error {
	instances, err := provider.ListInstances()
	if err != nil {
		log.Error("Failed listing instances")
		return err
	}
	outString, err := instancesToString(*input.OutputFormat, instances)
	if err != nil {
		log.Error("Failed stringifying instances")
		return err
	}
	// TODO Mock stdout
	fmt.Fprint(os.Stdout, outString)
	if *input.Interactive && len(instances) > 0 {
		// TODO Mock stdin and stderr
		reader := bufio.NewReader(os.Stdin)
		fmt.Fprintf(os.Stderr, "Choose an instance to connect to [1 - %d]: ", len(instances))
		textInput, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		i, err := strconv.Atoi(strings.ReplaceAll(textInput, "\n", ""))
		if err != nil {
			return err
		}
		log.WithField("index", i).Debug("Picked EC2 Instance")
		if i < 1 || i > len(instances) {
			return fmt.Errorf("instance index out of range: %d", i)
		}
		instance := instances[i-1]
		err = provider.StartSession(aws.TargetTypeInstanceID, instance.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (input *StartInput) listSessions(provider aws.CloudSessions) error {
	sessions, err := provider.ListSessions()
	if err != nil {
		log.Error("Failed listing instances")
		return err
	}
	outString, err := sessionsToString(*input.OutputFormat, sessions)
	if err != nil {
		log.Error("Failed stringifying instances")
		return err
	}
	// TODO Mock stdout
	fmt.Fprint(os.Stdout, outString)
	if *input.Interactive && len(sessions) > 0 {
		// TODO Mock stdin and stderr
		reader := bufio.NewReader(os.Stdin)
		fmt.Fprintf(os.Stderr, "Terminate session [1 - %d]: ", len(sessions))
		textInput, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		i, err := strconv.Atoi(strings.ReplaceAll(textInput, "\n", ""))
		if err != nil {
			return err
		}
		log.WithField("index", i).Debug("Picked session")
		if i < 1 || i > len(sessions) {
			return fmt.Errorf("session index out of range: %d", i)
		}
		chosenSession := sessions[i-1]
		err = provider.TerminateSession(chosenSession.SessionID)
		if err != nil {
			return err
		}
	}
	return nil
}
