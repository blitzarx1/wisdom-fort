# Wisdom-Fort
Wisdom-Fort is a TCP server that utilizes Proof of Work (PoW) challenges as a shield against DDoS attacks. It serves insightful quotes to clients, but only after they have successfully solved a PoW challenge. This project is designed to balance server-side protection mechanisms with a rewarding client-side experience

## Design Overview
The system was designed with an emphasis on simplicity, security, and resilience to DDoS attacks. A hash-based Proof of Work (PoW) system is used, which is simple to understand and implement, yet provides a strong protection against spam and flooding attacks. The PoW system generates unique challenges for each client, which prevent replay attacks, and the difficulty level is adjusted based on each client's request rate, enabling the system to handle different levels of demand and server capacity.

The process is split in 2 phases:
* get unique token and a challenge
* solve the challenge and get a quote

### Get a Challenge
This call is used to get a uniqut token and a challenge to solve.

May be considered as a public api enpoint as it is not requieres a token. Thus is not protected by PoW difficulty adjustment from DDoS attacks. The protection used on this phase is traditional rate limiter which limits the number of requests by ip.

### Solve the Challenge and Get a Quote
This call requires token to be present in the request. The token is obtained in the previous phase.

This phase is protected by PoW difficulty adjustment. The difficulty is adjusted based on the client's request rate. The higher the request rate, the higher the difficulty. This allows the server to handle different levels of demand and server capacity.

## Challenge Description
The challenge is hash-based PoW algorithm. The task is to find a sha256 hexadecimal string with a specific number of leading zero.
 
Client is given a unique token which needs to be concatenated with any number. The condition for the resulting string is for its sha256 hexadecimal representation to have a specific number of leading zeros.

The number of leading zeros is determined by the difficulty level, which is adjusted based on the client's request rate.

## Quotes
Once the PoW solution is verified, the server responds with a random wisdom quote.

## Choice of Proof of Work Algorithm
The Proof of Work (PoW) algorithm chosen for Wisdom-Fort is a hash-based PoW. This selection was driven by several key considerations:

- **Security:** The algorithm's challenge-response mechanism and the requirement for clients to produce a hash with specific properties helps guard against a range of attacks, including DDoS and replay attacks.

- **Simplicity:** Despite its robust security features, the algorithm is straightforward to implement and understand. This makes it accessible to a wide range of developers and users.

- **Scalability:** The difficulty level of the PoW challenge can be dynamically adjusted. This allows the server to respond to changes in demand and manage its resources effectively.

- **Fairness:** Every client is given a unique challenge, ensuring that the PoW system is fair. No client has an advantage over another, which promotes equal opportunity for all users to receive wisdom quotes.

- **Proven Effectiveness:** Hash-based PoW systems have been successfully used in a number of high-profile applications, such as in blockchain technology and cryptocurrency networks. This serves as a testament to their effectiveness in protecting systems against potential abuse.

## Getting Started

### Server
Server has a set of configuration parameters which can be set via environment variables. The default values are set in `.env` file. The following parameters are available:
* `PORT` - Port to listen on;
* `RPS_LIMIT_UNAUTH` - RPSLimitUnauth is ip rps limit for requests without valid token;
* `DIFF_MULT` - DiffMult is difficulty multiplier for challenges. If set to 1 the challenges trivial.Difficulty is equal to the client IPs RPS. 0 makes challenges trivial. Recommended value is 1.
* `CHALLENGE_TTL_SECONDS` - ChallengeTTLSeconds is expiration time for challenge in seconds. When the time is passed the challenge is considered invalid and the client needs to request a new one.

#### Run
To build and run the server, use the provided Dockerfile:

```sh
docker build -t wf-server -f Dockerfile.server .
docker run -p 8080:8080 wf-server
```

This will start the server id a docker container listening on port 8080 and with the default configuration taken from `.env` file.

### Client
Client is a package which provides possible realization of the client for wisdom-fort pow server.

#### Demo
Demo demonstrate the client-server interaction. 

It will demonstrate how the client-server interaction is organised and how rps limiting and difficulty adjustment works.

To run the demo:
* Start demo after the server is running locally. 
* Demo is supposed to be run with default server configuration. 
* Please check stdout for the output.

```sh
go run client/cmd/demo.go
```

You can also run demo from docker container:

```sh
docker build -t wf-client -f Dockerfile.client .
docker run --network=host wf-client
```

#### Understanding demo results
Demo will use client to send requests solving challenges and getting quoutes. Client will reach 7 rps loading server; default server config is setting difficulty equal to rps, that is why client will get challenge with `difficulty=7`. Most likely it will not be able to find the soulution before challenge ttl expires (5 seconds in default config) - client gets error `ErrInvalidSolution` because task has already been deleted and the solution is sent for the not existing task. 

If you want to see how client solves final solution before ttl expires you can set the `CHALLENGE_TTL_SECONDS` evironment variable to 120 seconds when starting server. It should be enough. So final command to start the sserver for demo to be successfull is 

```sh
docker run -p 8080:8080 --env CHALLENGE_TTL_SECONDS=120 wf-server
```

### Tests
To run tests, use the go test command:

```sh
go test ./...
```
