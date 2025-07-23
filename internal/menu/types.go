package menu

// Action represents a menu action
type Action struct {
	Name        string
	Description string
	Handler     func() error
}

// ActionType represents different types of actions
type ActionType int

const (
	ListDevicesAction ActionType = iota
	ScanAction
	ConnectAction
	DisconnectAction
	QuitAction
)
