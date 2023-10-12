# Word of Wisdom

## This is an example of TCP server providing random phrases and having protection from DDoS using the challenge-response protocol combined with the [PoW Hashcash algorithm](https://en.wikipedia.org/wiki/Hashcash)

## Getting started

- At first you have to install dependencies, at least Golang, to run the app locally or Docker.
- Environment variables are set up in the docker-compose.yml or can be specified when you start the app locally.
- To run the app use one of the Makefile commands described below

## Dependencies:

- [Go 1.17+](https://go.dev/dl/) installed (to run tests, start server or client without Docker)
- [golangci-lint](https://github.com/golangci/golangci-lint) tool installed (to run linters)
- [Docker](https://docs.docker.com/engine/install/) installed (to run docker-compose)

## Makefile commands:

#### Start server locally (without docker):

```
SERVER_ADDR=:8099 CHALLENGE_TIMEOUT=20s CHALLENGE_COMPLEXITY=10 READ_TIMEOUT=30s make start-server
```

#### Start client locally (without docker):

```
READ_TIMEOUT=30s SERVER_ADDR=:8099 make start-client
```

#### Build docker images of server and client and run them by docker-compose:

```
make build-and-start-docker
```

#### Run server and client by docker-compose (previously built images):

```
make start-docker
```

#### Run tests:

```
make test
```

#### Run linters:

```
make lint
```

## Why Hashcash?

- Hashcash algorithm is relatively easy to implement
- The proof of work of the hashcash function is efficiently auditable compared to the cost of the work, so the server will not spend a lot of computing resources to check the solution
- It completely fits the requirements of a tcp server which needs to protect itself from DDOS attacks

## Sources

- https://en.wikipedia.org/wiki/Hashcash
- http://www.hashcash.org/hashcash.pdf
