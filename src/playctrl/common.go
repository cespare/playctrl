package playctrl

const ProtocolVersion = 1

type Message struct {
	Version int    `json:"version"`
	Value   string `json:"value"`
}

var CmdNames = map[string]bool{
	"previous":   true,
	"playpause":  true,
	"next":       true,
	"volumeup":   true,
	"volumedown": true,
}
