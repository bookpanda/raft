package kvservice

// command type the KVService submits to Raft log
type Command struct {
	Kind         CommandKind
	Key, Value   string
	CompareValue string
	ResultValue  string
	ResultFound  bool

	// Raft ID of the service submitting this command
	ServiceID int

	// for command deduplication
	ClientID, RequestID int64

	IsDuplicate bool
}

type CommandKind int

const (
	CommandInvalid CommandKind = iota
	CommandGet
	CommandPut
	CommandAppend
	CommandCAS
)

var commandName = map[CommandKind]string{
	CommandInvalid: "invalid",
	CommandGet:     "get",
	CommandPut:     "put",
	CommandAppend:  "append",
	CommandCAS:     "cas",
}

func (ck CommandKind) String() string {
	return commandName[ck]
}
