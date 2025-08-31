# TCP server with Proof of Work DDoS protection

## Task
Test task for Server Engineer  
Design and implement “Word of Wisdom” tcp server.  
- TCP server should be protected from DDOS attacks with the Proof of Work ([https://en.wikipedia.org/wiki/Proof_of_work](https://en.wikipedia.org/wiki/Proof_of_work)), the challenge-response protocol should be used.  
- The choice of the POW algorithm should be explained.  
- After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.  
- Docker file should be provided both for the server and for the client that solves the POW challenge

## Usage

### Run
```bash
make up
```

## PoW algorithm
The TCP server uses the `Hashcash` stateless algorithm with `HMAC` and `TTL` support to verify client requests. `HMAC` prevents attackers from forging easier challenges, while `TTL` limits the lifetime of valid solutions.
Compared to memory-hard algorithms like `Argon2` or `Equihash`, this design is lightweight, simpler, and faster.

The client must find a `nonce` such that the resulting hash is lower than the target provided by the server.

The challenge has the following structure: `timestamp (8 bytes) || salt (8 bytes) || target (8 bytes) || hmac (32 bytes)`.
