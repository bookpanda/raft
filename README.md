# raft
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
```bash
./dotest.sh TestElectionFollowerComesBack
```