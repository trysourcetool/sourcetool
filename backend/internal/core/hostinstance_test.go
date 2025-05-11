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
		{HostInstanceStatus(100), "unknown"}, // 範囲外
	}
	for _, c := range cases {
		if got := c.status.String(); got != c.expected {
			t.Errorf("HostInstanceStatus(%d).String() = %q, want %q", c.status, got, c.expected)
		}
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
		{"invalid", HostInstanceStatusUnknown}, // 未知の文字列
	}
	for _, c := range cases {
		if got := HostInstanceStatusFromString(c.input); got != c.expected {
			t.Errorf("HostInstanceStatusFromString(%q) = %v, want %v", c.input, got, c.expected)
		}
	}
}
