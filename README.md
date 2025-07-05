# raft
```bash
./dotest.sh TestElectionFollowerComesBack
```

## 1. Elections
- higher term = higher priority

### Follower
- does not send any messages
- only replies to messages from leader or candidate
### Candidate
- starts election by incrementing term and voting for itself
- sends RequestVote RPC to all other servers
- waits for votes from majority of servers
- if it receives votes from majority, it becomes leader
- if it receives RequestVote RPC from a higher term, it becomes follower
### Leader
- sends AppendEntries RPC to all followers (heartbeat or log replication)

## 2. Commands and log replication
- client sends command to leader
- leader appends command to its log
- leader sends AppendEntries RPC to followers with new log entry
- followers append entry to their log and send back success response
- leader waits for majority of followers to respond
- if majority respond, leader commits the entry and sends CommitEntry to commitChan

### 2 RPC round-trips to commit a command
1. leader sends next log entries to followers
2. leader sends updated commit index to followers, who will then mark these entries as committed and will send them on the commit channel

### Election Safety
- prevent a candidate from winning an election unless its log is at least as up-to-date as a majority of peers in the cluster
- RV args: `lastLogIndex`, `lastLogTerm`, followers compare these fields to their own and decide whether the candidate is sufficiently up-to-date to be elected
- this also prevents `runaway` leaders from being elected, i.e. leaders that have no log entries or have stale log entries but higher term (e.g. they are separated due to network partition then rejoin)

## 3. Persistence, Optimizations
### Persistence
- needs only to persist currentTerm, votedFor, and log entries
### Delivery
- Raft = `at least once` delivery
- commands should have unique IDs to prevent duplicates