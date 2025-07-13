package kvservice

// command type the KVService submits to Raft log
type Command struct {
	Kind         CommandKind
	Key, Value   string
	CompareValue string
	ResultValue  string
	ResultFound  bool

	// Raft ID of the server submitting this command
	Id int
}

type CommandKind int

const (
	CommandInvalid CommandKind = iota
	CommandGet
	CommandPut
	CommandCAS
)

var commandName = map[CommandKind]string{
	CommandInvalid: "invalid",
	CommandGet:     "get",
	CommandPut:     "put",
	CommandCAS:     "cas",
}

func (ck CommandKind) String() string {
	return commandName[ck]
}
