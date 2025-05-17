package core

import "testing"

func TestHostInstanceStatus_String(t *testing.T) {
	cases := []struct {
		status   HostInstanceStatus
		expected string
	}{
		{HostInstanceStatusUnknown, "unknown"},
		{HostInstanceStatusOnline, "online"},
		{HostInstanceStatusUnreachable, "unreachable"},
		{HostInstanceStatus(100), "unknown"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.expected, func(t *testing.T) {
			t.Parallel()
			if got := c.status.String(); got != c.expected {
				t.Errorf("HostInstanceStatus(%d).String() = %q, want %q", c.status, got, c.expected)
			}
		})
	}
}

func TestHostInstanceStatusFromString(t *testing.T) {
	cases := []struct {
		input    string
		expected HostInstanceStatus
	}{
		{"unknown", HostInstanceStatusUnknown},
		{"online", HostInstanceStatusOnline},
		{"unreachable", HostInstanceStatusUnreachable},
		{"invalid", HostInstanceStatusUnknown},
	}
	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			t.Parallel()
			if got := HostInstanceStatusFromString(c.input); got != c.expected {
				t.Errorf("HostInstanceStatusFromString(%q) = %v, want %v", c.input, got, c.expected)
			}
		})
	}
}
