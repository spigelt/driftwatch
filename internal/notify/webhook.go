package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"driftwatch/internal/drift"
)

// WebhookPayload is the JSON body sent to the webhook endpoint.
type WebhookPayload struct {
	Timestamp string         `json:"timestamp"`
	HasDrift  bool           `json:"has_drift"`
	Summary   string         `json:"summary"`
	Reports   []drift.Report `json:"reports"`
}

// WebhookSink delivers drift reports to an HTTP endpoint.
type WebhookSink struct {
	url    string
	client *http.Client
	log    zerolog.Logger
}

// NewWebhookSink creates a WebhookSink that POSTs to the given URL.
func NewWebhookSink(url string, log zerolog.Logger) *WebhookSink {
	return &WebhookSink{
		url: url,
		client: &http.Client{Timeout: 10 * time.Second},
		log: log,
	}
}

// Send implements Sink. It serialises reports and POSTs them to the webhook URL.
func (w *WebhookSink) Send(reports []drift.Report) error {
	hasDrift := false
	for _, r := range reports {
		if r.HasDrift() {
			hasDrift = true
			break
		}
	}

	payload := WebhookPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		HasDrift:  hasDrift,
		Reports:   reports,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}

	w.log.Debug().Str("url", w.url).Int("status", resp.StatusCode).Msg("webhook delivered")
	return nil
}
