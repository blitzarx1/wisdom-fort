# Wisdom-Fort

Wisdom-Fort is a TCP server that utilizes Proof of Work (PoW) challenges as a shield against DDoS attacks. It serves insightful quotes to clients, but only after they have successfully solved a PoW challenge. This project is designed to balance server-side protection mechanisms with a rewarding client-side experience

## Design Overview

The system was designed with an emphasis on simplicity, security, and resilience to DDoS attacks. A hash-based Proof of Work (PoW) system is used, which is simple to understand and implement, yet provides a strong protection against spam and flooding attacks. The PoW system generates unique challenges for each client, which prevent replay attacks, and the difficulty level is adjusted based on each client's request rate, enabling the system to handle different levels of demand and server capacity.

The design also includes a unique token assigned to each client upon their first successful PoW solution. This token helps track and limit the number of requests from each client, thus providing additional protection against DDoS attacks.

## Choice of Proof of Work Algorithm

The Proof of Work (PoW) algorithm chosen for Wisdom-Fort is a Hash-based PoW. This selection was driven by several key considerations:

- **Security:** The algorithm's challenge-response mechanism and the requirement for clients to produce a hash with specific properties helps guard against a range of attacks, including DDoS and replay attacks.

- **Simplicity:** Despite its robust security features, the algorithm is straightforward to implement and understand. This makes it accessible to a wide range of developers and users.

- **Scalability:** The difficulty level of the PoW challenge can be dynamically adjusted. This allows the server to respond to changes in demand and manage its resources effectively.

- **Fairness:** Every client is given a unique challenge, ensuring that the PoW system is fair. No client has an advantage over another, which promotes equal opportunity for all users to receive wisdom quotes.

- **Proven Effectiveness:** Hash-based PoW systems have been successfully used in a number of high-profile applications, such as in blockchain technology and cryptocurrency networks. This serves as a testament to their effectiveness in protecting systems against potential abuse.

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

Once the PoW solution is verified and the client's request rate is checked, the server responds with a random wisdom quote.

## Testing

To run tests, use the go test command:

```sh
go test ./...
```

## TODO

Server:
- [ ] Tests
- [ ] Graceful shutdown with ctx

Client:
- [ ] Extract client from main