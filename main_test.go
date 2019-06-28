package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sensu/sensu-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestSendAlert(t *testing.T) {
	assert := assert.New(t)
	event := types.FixtureEvent("entity1", "check1")

	var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		expectedBody := `{"resource":"entity1","event":"check1","environment":"default","severity":"normal","status":"","service":["Sensu"],"group":"default","value":"","text":"","origin":"sensu-go/`
		assert.Contains(string(body), expectedBody)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status": "ok"}`))
		require.NoError(t, err)
	}))

	config.AlertaEndpoint.Value = apiStub.URL
	config.AlertaApiKey.Value = "demo-key"
	err := sendAlert(event)
	assert.NoError(err)
}

func TestMain(t *testing.T) {
	assert := assert.New(t)
	file, _ := ioutil.TempFile(os.TempDir(), "sensu-handler-alerta-")
	defer func() {
		_ = os.Remove(file.Name())
	}()

	event := types.FixtureEvent("entity1", "check1")
	eventJSON, _ := json.Marshal(event)
	_, err := file.WriteString(string(eventJSON))
	require.NoError(t, err)
	require.NoError(t, file.Sync())
	_, err = file.Seek(0, 0)
	require.NoError(t, err)
	stdin = file
	requestReceived := false

	var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived = true
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status": "ok"}`))
		require.NoError(t, err)
	}))

	oldArgs := os.Args
	os.Args = []string{"alerta", "--endpoint-url", apiStub.URL}
	defer func() { os.Args = oldArgs }()

	main()
	assert.True(requestReceived)
}
