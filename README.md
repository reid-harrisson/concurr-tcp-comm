# Concurrent TCP Communication

## Requirements

There are two parts to this coding challenge. The intent is to take 30-60 minutes for an initial rough implementation. Please note any areas where future robustness could be considered. The goal is to write a concurrent server and client that speaks a protocol, and then test that the client and server are correct under concurrent loaded conditions. Please write the solution in Go.

1. The Client and Server

   The client connects to the server with your choice of TCP, UDP, or websockets, and sends a packet that is the string message X, where X can be any string. The client listens on the
   connection for one or messages of the format message Y or ok X. When the client receives ok X for the same X that was originally sent, it should close the connection. When the client receives message Y, it should immediately reply with ok Y. If the client does not receive ok X after some time, it should disconnect. If the client disconnects before receiving ok X, it should reconnect with message X. The server is a single instance that listens for your choice of TCP, UDP, or websockets, and maintains a collection of a active connections. When a connection is closed it should be removed from the active connections. When the server receives message X it should forward that message to a randonly selected connection, including the connection that originally sent the message. When the server receives ok X it shold forward the message to the connections that originally sent message X.

2. Load Test

   Write a test to simulate N clients connecting in parallel with some random initial delay over a time of 60s. Test for large values of N such as 100k. Test that all sent message X eventually receive a corresponding ok X. The test should make sure that the client connections close naturally as part of the protocol.The test should make sure that the server ends in a zero state when all client connections are closed. Try testing with various delays between receiving message X and sending ok X.

## Prequiresites

- Install latest go version

## Installation

- Clone the project from <https://github.com/princecharming0115/bringyour-test> and navigate.
- Install dependencies using this command.

```bash
go mod tidy
```

- Rename `.env.example` to `.env`

- Run client and server at the same computer using this command.

```bash
go run ./cmd
```

- Run server using this command.

```bash
go run ./cmd/server
```

- Run client using this command.

```bash
go run ./cmd/client
```

Server is running on 8000 and 10 clients are running on different ports.

- You can test this app using this command.

```bash
go test ./test
```
