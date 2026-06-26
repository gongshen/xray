package middleware

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

func TestSysOperationRecordSupportsTargets(t *testing.T) {
	record := system.SysOperationRecord{Targets: `{"id":1}`}
	if record.Targets != `{"id":1}` {
		t.Fatalf("Targets = %q, want %q", record.Targets, `{"id":1}`)
	}
}

func TestExtractOperationTargetsCapturesBodyTargets(t *testing.T) {
	got := extractOperationTargets(url.Values{}, []byte(`{"targets":["google.com","1.1.1.1"],"page":1}`))
	targets := decodeOperationTargets(t, got)

	want := []interface{}{"google.com", "1.1.1.1"}
	if !reflect.DeepEqual(targets["targets"], want) {
		t.Fatalf("targets = %#v, want %#v; full = %#v", targets["targets"], want, targets)
	}
	if _, ok := targets["page"]; ok {
		t.Fatalf("page should not be captured as an operation target: %#v", targets)
	}
}

func TestExtractOperationTargetsCapturesIdentifiersFromBodyAndQuery(t *testing.T) {
	query := url.Values{"id": {"5"}, "page": {"1"}}
	got := extractOperationTargets(query, []byte(`{"ID":12,"server_id":3,"name":"hk-node","pageSize":10}`))
	targets := decodeOperationTargets(t, got)

	want := map[string]interface{}{
		"ID":        json.Number("12"),
		"id":        "5",
		"name":      "hk-node",
		"server_id": json.Number("3"),
	}
	if !reflect.DeepEqual(targets, want) {
		t.Fatalf("operation targets = %#v, want %#v", targets, want)
	}
}

func TestExtractOperationTargetsReturnsEmptyWhenNoTargetFieldsExist(t *testing.T) {
	got := extractOperationTargets(url.Values{"page": {"1"}}, []byte(`{"page":1,"pageSize":10}`))
	if got != "" {
		t.Fatalf("operation targets = %q, want empty string", got)
	}
}

func decodeOperationTargets(t *testing.T, raw string) map[string]interface{} {
	t.Helper()
	if raw == "" {
		t.Fatal("operation targets should not be empty")
	}

	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()

	var targets map[string]interface{}
	if err := decoder.Decode(&targets); err != nil {
		t.Fatalf("decode operation targets %q: %v", raw, err)
	}
	return targets
}
