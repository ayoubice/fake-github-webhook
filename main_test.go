package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var payloadSample = payload{
	"ref":      "refs/heads/changes",
	"before":   "9049f1265b7d61be4a8904a9a27120d2064dab3b",
	"after":    "0d1a26e67d8f5eaf1f6ba5c57fc3c7d91ac0fd1c",
	"created":  false,
	"deleted":  false,
	"forced":   false,
	"base_ref": nil,
	"compare":  "https://github.com/baxterthehacker/public-repo/compare/9049f1265b7d...0d1a26e67d8f",
}

const dataDir = "./data"

func TestSendPayload(t *testing.T) {

	echoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var p payload
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			t.Fatal("[echo Server] enable to decode the request body as an event payload")
		}

		if !cmp.Equal(p, payloadSample) {
			t.Errorf("mismatch between payloads (sent, received) %s", cmp.Diff(payloadSample, p))
		}
	}))

	defer echoSrv.Close()

	err := sendPayload(echoSrv.URL, payloadSample)
	if err != nil {
		t.Fatalf("Error sending payload to echoServer: %v", err)
	}
}

func TestLoadSequences(t *testing.T) {
	s, err := loadSequences(dataDir)
	if err != nil {
		t.Fatalf("error loading sequences %s", err.Error())
	}

	if len(s) < 1 {
		t.Errorf("no sequence loaded")
	}
}
