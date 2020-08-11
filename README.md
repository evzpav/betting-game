# Betting Game

	It is web app made with Golang server (no database) with websockets and rest API. Frontend is done in VueJs. Golang API serves the static files from frontend.
	It is a betting game, where player inputs his name and chooses 2 numbers.
	One number is drawn every 1s. Score is calculated and ranking between players is shown. Game lasts max 30 rounds.
	It is a multiplayer game, and minimum 2 players are needes to start the game.
	More info below.

 <div>
    <h2>Starting the game</h2>
    <ul>
      <li>2 players are needed to start the game (1 browser tab one player)</li>
    </ul>
    <ol>
      <li>On the initial page click "Play"</li>
      <li>Then fill your name</li>
      <li>Click on 2 numbers of your choice (1 to 10)</li>
      <li>Click "Join"</li>
      <li>Open another tab/window of this link in your browser</li>
      <li>Repeat in this other tab the steps 1 to 4</li>
      <li>Game will start</li>
    </ol>
    <ul>
      <li>When new players join in a running game they will only observe until next game starts, then they will join automatically.</li>
    </ul>
    <h2>Game Mechanics:</h2>
    <ul>
      <li>New random number generated every 1 second</li>
      <li>Maximum 30 rounds</li>
      <li>Ends if any player reaches exactly 21 points</li>
      <li>10 seconds of interval between games</li>
    </ul>
    <h2>Score:</h2>
    <ul>
      <li>
        Exact match
        <ul>
          <li>+5 points if there is exact match with one of the chosen number</li>
        </ul>
      </li>
      <li>
        Inside bounds
        <ul>
          <li>+5-(upper bound - lower bound); e.g.: 3 and 8 chosen, 7 is the round number, then 5-(8-3)=0</li>
        </ul>
      </li>
      <li>
        Out of bounds
        <ul>
          <li>-1 point; e.g.: 3 and 8 chosen, 9 is the round number</li>
        </ul>
      </li>
    </ul>
    <h2>Winner criteria (in this order):</h2>
    <ul>
      <li>Points</li>
      <li>Highest upper number chosen</li>
      <li>Highest lower number chosen</li>
      <li>Name ascending</li>
    </ul>
</div>

### Assumptions/Constraints:
- Players cannot have the same name
- Player names are saved in lowercase and cannot be edited
- Players cannot leave the game
- Numbers chosen cannot be changed
- New browser tab, means new player (player saved on session storage)


## Run project

Below are 2 options to run the project locally:

### Option 1 (recommended):
#### Prerequisites
- [GNU Make](https://www.gnu.org/software/make/)
- [Docker](http://docker.com)


### Running locally on Docker
It will build frontend and backend, generate a single docker image and run it on port 8888:
```bash
make run-docker
# Open in the browser: http://localhost:8888
```
### Option 2:
#### Prerequisites
- [GNU Make](https://www.gnu.org/software/make/)
- [Golang](http://golang.org/)(>11.0)
- [NodeJs](http://nodejs.org)

### Running locally
```bash
# Run server locally on port 8787
make run-local

# Install Vue frontend dependencies and run it locally
make run-frontend
# Open in the browser: http://localhost:8080
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

### Environment variables

```bash
HOST
PORT
LOGGER_LEVEL
```

### Deployment
The command below will build frontend Vue app for production(generate static files) and will build Go backend and generate the server binary.
Afterwards it generates a small Alpine docker image that it is deployable
```bash
make image
```