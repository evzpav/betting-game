# Betting Game

It is web app made with Golang backend and Vue frontend.
It is a betting game, where player inputs his name and chooses 2 numbers.

## Getting Started

### Prerequisites

- [Golang](http://golang.org/)(>11.0)
- [GNU Make](https://www.gnu.org/software/make/)
- [Docker](http://docker.com)

### Environment variables

```bash
	HOST
	PORT
	LOGGER_LEVEL
```

### Installing and running locally

```bash

#Run server locally on port 8787
make run-local

#Run vue frontend locally 
make run-frontenv
# open in the browser: http://localhost:8080

# or

#Run on docker - it will build frontend and backend and generate a single docker image
make run-docker
# open in the browser: http://localhost:8888
```

## Starting the game
- 2 players are needed to start the game
1. Open the url of the project on Chrome click on "Play"
2. Then fill your name and choose 2 numbers and click "Join"
3. You will be waiting until another player joins.
4. Repeat steps 1 and 2 (choosing another name, same name of another player is no allowed)
5. Game will start

When new players join when a game is running they will only observe it until next game starts, then they will join automatically



