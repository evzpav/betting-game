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

### Running locally

```bash

# Run server locally on port 8787
make run-local

# Run vue frontend locally 
make run-frontend
# Open in the browser: http://localhost:8080

```

### Running locally on Docker
It will build frontend and backend, generate a single docker image and run it on port 8888:
```bash
make run-docker
# Open in the browser: http://localhost:8888
```

### Unit Tests
```bash
# It will run tests locally
make test

# or

# It will run the same tests on docker
make test-docker
```

### Lint
```bash
# It will run linting on docker
make lint-docker
```


### Deployment
The command below will build frontend Vue app for production(generate static files) and will build Go backend and generate the server binary.
Afterwards it generates a small Alpine docker image that it is deployable
```bash
make image
```


## Starting the game
- 2 players are needed to start the game
1. Open the url of the project on Chrome click on "Play"
2. Then fill your name and choose 2 numbers and click "Join"
3. You will be waiting until another player joins.
4. Repeat steps 1 and 2 (choose another name, because same name of another player is no allowed)
5. Game will start

When new players join when a game is running they will only observe it until next game starts, then they will join automatically.

### Assumptions/Constraints:
- Players cannot have the same name
- Player names are saved in lowercase
- Players cannot leave the game
- Numbers chosen cannot be changed
- 

