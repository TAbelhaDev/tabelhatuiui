package tuiui

import (
	"testing"
)

func TestParseIPCArgs(t *testing.T) {
	got, err := ParseIPCArgs([]string{"jobs.list", "dirty=true", "--json"})
	if err != nil {
		t.Fatalf("ParseIPCArgs: %v", err)
	}
	if got.Method != "jobs.list" {
		t.Fatalf("Method = %q, want jobs.list", got.Method)
	}
	if !got.JSON {
		t.Fatal("JSON = false, want true")
	}
	if v, ok := got.Filters["dirty"]; !ok || v != "true" {
		t.Fatalf("Filters = %v, want dirty=true", got.Filters)
	}
}

func TestParseIPCArgsErrors(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"no args", nil},
		{"missing --json", []string{"jobs.list"}},
		{"invalid arg", []string{"jobs.list", "not-a-kv", "--json"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ParseIPCArgs(tc.args); err == nil {
				t.Fatalf("ParseIPCArgs(%v): expected error, got nil", tc.args)
			}
		})
	}
}
