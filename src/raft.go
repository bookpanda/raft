package raft

import (
	"sync"
	"time"
)

type LogEntry struct {
	Command any
	Term    int
}

type CMState int

const (
	Follower CMState = iota
	Candidate
	Leader
	Dead
)

func (s CMState) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	case Dead:
		return "Dead"
	default:
		panic("unreachable")
	}
}

type ConsensusModule struct {
	mu      sync.Mutex
	id      int // same as server.serverId
	peerIds []int
	server  *Server

	currentTerm int
	votedFor    int
	log         []LogEntry

	state              CMState
	electionResetEvent time.Time
}

func NewConsensusModule(id int, peerIds []int, server *Server, ready <-chan any) *ConsensusModule {
	cm := &ConsensusModule{
		id:       id,
		peerIds:  peerIds,
		server:   server,
		state:    Follower,
		votedFor: -1,
	}

	go func() {
		<-ready
		cm.mu.Lock()
		cm.electionResetEvent = time.Now()
		cm.mu.Unlock()
		// cm.RunElectionTimer()
	}()

	return cm
}
