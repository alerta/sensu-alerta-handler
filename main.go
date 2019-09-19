package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"net/http"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugins-go-library/sensu"
)
// HandlerConfig represents plugin configuration settings.
type HandlerConfig struct {
	sensu.PluginConfig
	AlertaEndpoint       string
	AlertaAPIKey         string
	Environment          string
}

const (
	endpointURL = "endpoint-url"
	apiKey      = "api-key"
	environment = "environment"

	defaultEndpointURL string = "http://localhost:8080"
)

var (
	config = HandlerConfig{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-alerta-handler",
			Short:    "The Sensu Go Alerta handler for event forwarding",
			Timeout:  10,
			Keyspace: "sensu.io/plugins/alerta/config",
		},
	}

	alertaConfigOptions = []*sensu.PluginConfigOption{
		{
			Path:      endpointURL,
			Env:       "ALERTA_ENDPOINT",
			Argument:  endpointURL,
			Shorthand: "",
			Default:   defaultEndpointURL,
			Usage:     "API endpoint URL",
			Value:     &config.AlertaEndpoint,
		},
		{
			Path:      apiKey,
			Env:       "ALERTA_API_KEY",
			Argument:  apiKey,
			Shorthand: "K",
			Default:   "",
			Usage:     "API key for authenticated access",
			Value:     &config.AlertaAPIKey,
		},
		{
			Path:      environment,
			Argument:  environment,
			Shorthand: "E",
			Default:   "Entity Namespace",
			Usage:     "Environment eg. Production, Development",
			Value:     &config.Environment,
		},
	}
)

func main() {
  goHandler := sensu.NewGoHandler(&config.PluginConfig, alertaConfigOptions, checkArgs, sendAlert)
  goHandler.Execute()
}

func checkArgs(event *corev2.Event) error {
	if !event.HasCheck() {
		return fmt.Errorf("event does not contain check")
	}
	return nil
}

// Alert represents an event to be sent to Alerta.
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

// TODO(satterly) could make these severity lookups configurable

func eventSeverity(event *corev2.Event) string {
	switch event.Check.Status {
	case 0:
		return "normal"
	case 2:
		return "critical"
	default:
		return "warning"
	}
}

func sendAlert(event *corev2.Event) error {
	environment := config.Environment
	if environment == "" {
		environment = event.Entity.Namespace
	}

	hostname, _ := os.Hostname()
	data := &Alert{
		Resource: event.Entity.Name,
		Event: event.Check.Name,
		Environment: environment,
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
	url := fmt.Sprintf("%s/alert?api-key=%s", config.AlertaEndpoint, config.AlertaAPIKey)

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
