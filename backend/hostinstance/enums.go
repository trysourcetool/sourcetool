package hostinstance

type HostInstanceStatus int

const (
	HostInstanceStatusUnknown HostInstanceStatus = iota
	HostInstanceStatusOnline
	HostInstanceStatusUnreachable
	HostInstanceStatusOffline
	HostInstanceStatusShuttingDown

	hostInstanceStatusUnknown      = "unknown"
	hostInstanceStatusOnline       = "online"
	hostInstanceStatusUnreachable  = "unreachable"
	hostInstanceStatusOffline      = "offline"
	hostInstanceStatusShuttingDown = "shuttingDown"
)

func (s HostInstanceStatus) String() string {
	switch s {
	case HostInstanceStatusOnline:
		return hostInstanceStatusOnline
	case HostInstanceStatusUnreachable:
		return hostInstanceStatusUnreachable
	case HostInstanceStatusOffline:
		return hostInstanceStatusOffline
	case HostInstanceStatusShuttingDown:
		return hostInstanceStatusShuttingDown
	default:
		return hostInstanceStatusUnknown
	}
}

func HostInstanceStatusFromString(s string) HostInstanceStatus {
	switch s {
	case hostInstanceStatusOnline:
		return HostInstanceStatusOnline
	case hostInstanceStatusUnreachable:
		return HostInstanceStatusUnreachable
	case hostInstanceStatusOffline:
		return HostInstanceStatusOffline
	case hostInstanceStatusShuttingDown:
		return HostInstanceStatusShuttingDown
	default:
		return HostInstanceStatusUnknown
	}
}
