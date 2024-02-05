## Idea
This is a basic implementation of the raft consensus algorithm on top of a Key-Value store to create a distributed
key-value state machine.

### steps:
- In-memory Key-Value store
- Http endpoint for interractions
    - Get value
    - Set Value
    - Join Cluster

- Config loading from args/config file
- Tests
