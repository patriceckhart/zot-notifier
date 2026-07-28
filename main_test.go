package main

import (
	"testing"

	"github.com/patriceckhart/zot/packages/agent/ext"
)

func TestShouldNotify(t *testing.T) {
	tests := []struct {
		name  string
		event ext.Event
		want  bool
	}{
		{name: "confirmation requested", event: ext.Event{Name: "tool_confirmation_requested"}, want: true},
		{name: "final reply", event: ext.Event{Name: "turn_end", Stop: "end"}, want: true},
		{name: "intermediate tool turn", event: ext.Event{Name: "turn_end", Stop: "tool_use"}, want: false},
		{name: "aborted turn", event: ext.Event{Name: "turn_end", Stop: "aborted"}, want: false},
		{name: "ordinary tool call", event: ext.Event{Name: "tool_call"}, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := shouldNotify(test.event); got != test.want {
				t.Fatalf("shouldNotify(%+v) = %v, want %v", test.event, got, test.want)
			}
		})
	}
}
