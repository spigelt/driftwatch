package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"driftwatch/internal/drift"
)

func TestWebhookSink_SendsPayload(t *testing.T) {
	var received WebhookPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&received))
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	reports := []drift.Report{makeReport(true)}
	sink := NewWebhookSink(srv.URL, zerolog.Nop())
	err := sink.Send(reports)

	require.NoError(t, err)
	assert.True(t, received.HasDrift)
	assert.NotEmpty(t, received.Timestamp)
	assert.Len(t, received.Reports, 1)
}

func TestWebhookSink_NoDrift_HasDriftFalse(t *testing.T) {
	var received WebhookPayload

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&received))
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	reports := []drift.Report{makeReport(false)}
	sink := NewWebhookSink(srv.URL, zerolog.Nop())
	err := sink.Send(reports)

	require.NoError(t, err)
	assert.False(t, received.HasDrift)
}

func TestWebhookSink_ServerError_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	sink := NewWebhookSink(srv.URL, zerolog.Nop())
	err := sink.Send([]drift.Report{makeReport(false)})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestWebhookSink_InvalidURL_ReturnsError(t *testing.T) {
	sink := NewWebhookSink("http://127.0.0.1:0/no-server", zerolog.Nop())
	err := sink.Send([]drift.Report{makeReport(false)})
	require.Error(t, err)
}
