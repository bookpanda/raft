package raft

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"os"
	"sync"
	"time"
)

// simluate delays/drops in the network
type RPCProxy struct {
	mu sync.Mutex
	cm *ConsensusModule

	// -1: not dropping any calls
	//  0: dropping all calls now
	// >0: start dropping calls after this number is made
	numCallsBeforeDrop int
}

func NewProxy(cm *ConsensusModule) *RPCProxy {
	return &RPCProxy{cm: cm, numCallsBeforeDrop: -1}
}

type RequestVoteArgs struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

func (rpp *RPCProxy) RequestVote(args RequestVoteArgs, reply *RequestVoteReply) error {
	if len(os.Getenv("RAFT_UNRELIABLE_RPC")) > 0 {
		dice := rand.Intn(10)
		if dice == 9 { // 10% chance to drop the request
			rpp.cm.dlog("drop RequestVote")
			return fmt.Errorf("RPC failed")
		} else if dice == 8 { // 10% chance to delay the request
			rpp.cm.dlog("delay RequestVote")
			time.Sleep(75 * time.Millisecond)
		}
	} else {
		time.Sleep(time.Duration(1+rand.Intn(5)) * time.Millisecond)
	}
	return rpp.cm.RequestVote(args, reply)
}

type AppendEntriesArgs struct {
	Term     int
	LeaderId int

	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool

	// faster conflict resolution
	ConflictIndex int
	ConflictTerm  int
}

func (rpp *RPCProxy) AppendEntries(args AppendEntriesArgs, reply *AppendEntriesReply) error {
	if len(os.Getenv("RAFT_UNRELIABLE_RPC")) > 0 {
		dice := rand.Intn(10)
		if dice == 9 { // 10% chance to drop the request
			rpp.cm.dlog("drop AppendEntries")
			return fmt.Errorf("RPC failed")
		} else if dice == 8 { // 10% chance to delay the request
			rpp.cm.dlog("delay AppendEntries")
			time.Sleep(75 * time.Millisecond)
		}
	} else {
		time.Sleep(time.Duration(1+rand.Intn(5)) * time.Millisecond)
	}
	return rpp.cm.AppendEntries(args, reply)
}

func (rpp *RPCProxy) Call(peer *rpc.Client, method string, args any, reply any) error {
	rpp.mu.Lock()
	if rpp.numCallsBeforeDrop == 0 {
		rpp.mu.Unlock()
		rpp.cm.dlog("drop Call %s: %v", method, args)
		return fmt.Errorf("RPC failed")
	}

	if rpp.numCallsBeforeDrop > 0 {
		rpp.numCallsBeforeDrop--
	}
	rpp.mu.Unlock()
	return peer.Call(method, args, reply)
}

func (rpp *RPCProxy) DropCallsAfterN(n int) {
	rpp.mu.Lock()
	defer rpp.mu.Unlock()

	rpp.numCallsBeforeDrop = n
}

func (rpp *RPCProxy) DontDropCalls() {
	rpp.mu.Lock()
	defer rpp.mu.Unlock()

	rpp.numCallsBeforeDrop = -1
}
