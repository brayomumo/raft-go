## Networking
This describes the base network required for the node to interract with outside world.
Required entry and exit points include:
    - Sockets - for communication with other nodes
    - Http - For communication with clients

## Sockets
This is the main communication mechanism for Node-Node communication. We can use raw sockets(berkely sockets) for communication or advanced socket libraries like ZeroMQ
Using ZeroMq gives us connection handling, incase connection to other nodes


