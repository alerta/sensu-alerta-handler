package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"net/http"

	"github.com/sensu/sensu-go/types"
	"github.com/spf13/cobra"
)

type HandlerConfigOption struct {
	Value string
	Path  string
	Env   string
}

type HandlerConfig struct {
	AlertaEndpoint  HandlerConfigOption
	AlertaApiKey    HandlerConfigOption
	Timeout         int
	Keyspace        string
}

var (
	stdin  *os.File
	config = HandlerConfig{
		// default values
		AlertaEndpoint: HandlerConfigOption{Path: "endpoint-url", Env: "SENSU_ALERTA_ENDPOINT"},
		AlertaApiKey:    HandlerConfigOption{Path: "api-key", Env: "SENSU_ALERTA_API_KEY"},
		Timeout:         10,
		Keyspace:        "sensu.io/plugins/alerta/config",
	}
	options = []*HandlerConfigOption{
		// iterable slice of user-overridable configuration options
		&config.AlertaEndpoint,
		&config.AlertaApiKey,
	}
)

func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func configureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sensu-alerta-handler",
		Short: "The Sensu Go Alerta handler for event forwarding",
		RunE:  run,
	}

	cmd.Flags().StringVarP(&config.AlertaEndpoint.Value,
		"endpoint-url",
		"",
		os.Getenv("ALERTA_ENDPOINT"),
		"API endpoint URL.")

	cmd.Flags().StringVarP(&config.AlertaApiKey.Value,
		"api-key",
		"K",
		os.Getenv("ALERTA_API_KEY"),
		"API key for authenticated access.")

	_ = cmd.MarkFlagRequired("endpoint-url")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		_ = cmd.Help()
		return fmt.Errorf("invalid argument(s) received")
	}

	if stdin == nil {
		stdin = os.Stdin
	}

	eventJSON, err := ioutil.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("failed to read stdin: %s", err)
	}

	event := &types.Event{}
	err = json.Unmarshal(eventJSON, event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal stdin data: %s", err)
	}

	if config.AlertaEndpoint.Value == "" {
		_ = cmd.Help()
		return fmt.Errorf("API endpoint URL is empty")
	}

	configurationOverrides(&config, options, event)

	if err = event.Validate(); err != nil {
		return fmt.Errorf("failed to validate event: %s", err)
	}

	if !event.HasCheck() {
		return fmt.Errorf("event does not contain check")
	}

	return sendAlert(event)
}

func configurationOverrides(config *HandlerConfig, options []*HandlerConfigOption, event *types.Event) {
	if config.Keyspace == "" {
		return
	}
	for _, opt := range options {
		if opt.Path != "" {
			// compile the Annotation keyspace to look for configuration overrides
			k := path.Join(config.Keyspace, opt.Path)
			switch {
			case event.Check.Annotations[k] != "":
				opt.Value = event.Check.Annotations[k]
				log.Printf("Overriding default handler configuration with value of \"Check.Annotations.%s\" (\"%s\")\n", k, event.Check.Annotations[k])
			case event.Entity.Annotations[k] != "":
				opt.Value = event.Entity.Annotations[k]
				log.Printf("Overriding default handler configuration with value of \"Entity.Annotations.%s\" (\"%s\")\n", k, event.Entity.Annotations[k])
			}
		}
	}
}

type Alert struct {
	Resource    string `json:"resource"`
	Event       string `json:"event"`
	Environment string `json:"environment"`
	Severity    string `json:"severity"`
	Correlate   []string `json:"correlate,omitempty"`
	Status      string `json:"status"`
	Service     []string `json:"service"`
	Group       string `json:"group"`
	Value       string `json:"value"`
	Text        string `json:"text"`
	Tags        []string `json:"tags,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	Origin      string `json:"origin"`
	Type        string `json:"type"`
	CreateTime  string `json:"createTime,omitempty"`
	Timeout     int `json:"timeout"`
	RawData     string `json:"rawData"`
}

func eventSeverity(event *types.Event) string {
	switch event.Check.Status {
	case 0:
		return "normal"
	case 2:
		return "critical"
	default:
		return "warning"
	}
}

func sendAlert(event *types.Event) error {
	hostname, _ := os.Hostname()
	data := &Alert{
		Resource: event.Entity.Name,
		Event: event.Check.Name,
		Environment: event.Entity.Namespace,
		Severity: eventSeverity(event),
		Service: []string{"Sensu"},
		Group: event.Check.Namespace,
		Value: event.Check.State,
		Text: event.Check.Output,
		Attributes: event.Entity.Labels,
		Origin: fmt.Sprintf("sensu-go/%s", hostname),
		Type: "sensuAlert",
		RawData: event.Entity.String(),
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/alert?api-key=%s", config.AlertaEndpoint.Value, config.AlertaApiKey.Value )

	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	_, _ = fmt.Fprintf(os.Stdout, "response: %s\n", body)

	var r map[string]interface{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return err
	}
	return nil
}
