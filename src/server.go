package raft

import (
	"net"
	"net/rpc"
	"sync"
)

type Server struct {
	mu sync.Mutex

	serverId int
	peerIds  []int

	cm       *ConsensusModule
	rpcProxy *RPCProxy

	rpcServer *rpc.Server
	listener  net.Listener

	peerClients map[int]*rpc.Client

	ready <-chan any // receive-only channel
	quit  chan any   // bidirectional (normal) channel
	wg    sync.WaitGroup
}

func NewServer(serverId int, peerIds []int, ready <-chan any) *Server {
	s := new(Server)
	s.serverId = serverId
	s.peerIds = peerIds
	s.peerClients = make(map[int]*rpc.Client)
	s.ready = ready
	s.quit = make(chan any)
	return s
}

func (s *Server) Serve() {
	s.mu.Lock()
	s.cm = NewConsensusModule(s.serverId, s.peerIds, s, s.ready)

	s.rpcServer = rpc.NewServer()
	s.rpcProxy = &RPCProxy{cm: s.cm}
}

// simluate delays/drops in the network
type RPCProxy struct {
	cm *ConsensusModule
}
