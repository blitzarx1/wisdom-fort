# wisdom-fort

This project implements a TCP server that provides a quote from the "Word of Wisdom" book after a client successfully solves a Proof of Work (PoW) challenge. The server is designed to be resistant to DDoS attacks by requiring a PoW solution before serving a quote.

## Design Overview

The system was designed with an emphasis on simplicity, security, and resilience to DDoS attacks. A hash-based Proof of Work (PoW) system is used, which is simple to understand and implement, yet provides a strong protection against spam and flooding attacks. The PoW system generates unique challenges for each client, which prevent replay attacks, and the difficulty level is adjusted based on each client's request rate, enabling the system to handle different levels of demand and server capacity.

The design also includes a unique token assigned to each client upon their first successful PoW solution. This token helps track and limit the number of requests from each client, thus providing additional protection against DDoS attacks.

## Choice of Hash Function

The SHA-256 cryptographic hash function is chosen for the PoW system. This choice was based on several factors:

- **Security:** SHA-256 is currently considered to be very secure. It produces a 256-bit (32-byte) hash, which provides a large enough output size to be resistant to collisions (two different inputs producing the same output).

- **Performance:** SHA-256 strikes a good balance between security and performance. It's not the fastest hash function, but it's fast enough for many applications, and its security is well-regarded.

- **Availability and Use:** SHA-256 is widely available and used. It's included in the standard libraries of most programming languages, including Go, which is used in this project.

- **Standardization:** SHA-256 is standardized by the National Institute of Standards and Technology (NIST) in the United States, which means it has undergone extensive scrutiny and testing.

- **Use in PoW Systems:** SHA-256 is used in the PoW system of Bitcoin, the most well-known application of PoW. This demonstrates its effectiveness in this kind of application.

## Getting Started

### Server

To build and run the server, use the provided Dockerfile:

```sh
docker build -t wow-server -f Dockerfile.server .
docker run -p 8080:8080 wow-server
```

This will start the server on port 8080.

### Client

To build and run the client, use the provided Dockerfile:

```sh
docker build -t wow-client -f Dockerfile.client .
docker run wow-client
```

By default, the client will try to connect to a server running on `localhost:8080`. This can be configured by setting the `SERVER_HOST` and `SERVER_PORT` environment variables.

## Quotes

Once the PoW solution is verified and the client's request rate is checked, the server responds with a random wisdom.

## Testing

To run tests, use the go test command:

```sh
go test ./...
```
