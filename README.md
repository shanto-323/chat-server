# Chat-Server
<h4><b>A scalable chat server with three services that maintainable, and horizontally scalable.</h4></b>

### Chat Server Architecture 

- **Gateway** – Handles connections from clients. It can scale horizontally, and clients can connect to any instance dynamically. The gateway simply routes requests and responses; it doesn’t store data.

- **Client Service** – Keeps track of all connected clients, their sessions, and status. It acts as the source of truth for client information.

- **Message Service** – Stores and manages all chat messages. It serves messages to clients on request.

<br>
<img src="asset/diagram.gif" width="500" height="450"/>

</br></br>

# Connection Process
## **Step-1 :** Connect to the server</h1> 
```json
  ws://localhost:8080/ws
```
<br>

## **Step-2 :** Authentication
<h6>REQUEST</h6>

```json
{
    "method": "signup/signin",
    "credentials": {
        "username": "username",
        "password": "password"
    }
}
```
<h6>RESPONSE</h6>

```json
{
    "type": "auth",
    "payload": {
        "status": "true/false",
        "uid": "empty/value"
    }
}
```

## **Step-3 :** Select Peer to Chat
<h4>When the user clicks a client to chat, they send a peerId to the system. This triggers the Message Service to return chat history and subscribe the client to real-time updates for that peer, so any new messages are delivered immediately.</h4>

```json
{
    "type": "peer",
    "payload": {
        "remote_peer_id": "peer_id"
    }
}
```