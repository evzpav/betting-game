<template>
  <div>
    <h2 class="subtitle">Leaderboard</h2>

    <div class="card mb">
      <div class="card-header">
        <h4 class="title is-4">Current Game</h4>
      </div>
      <div class="card-content">
        <div v-if="!gameRunning">Waiting for another player to join</div>

        <div v-if="gameRunning" class="notification is-warning">
          <div>Round: {{game.roundCounter}}</div>
          <div>Number: {{game.randomNumber}}</div>
        </div>

        <div
          v-if="winner && !gameRunning"
          class="notification is-success"
        >Winner is: {{winner.name}}</div>

        <table class="table">
          <thead>
            <tr>
              <th>
                <abbr title="Position">Pos</abbr>
              </th>
              <th>Player</th>
              <th>Numbers</th>
              <th>Points</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="(player,i) in leaderboard" :key="player.id">
              <td>{{i+1}}</td>
              <td>{{player.name}}</td>
              <td>{{player.numbers[0]}}, {{player.numbers[1]}}</td>
              <td>{{player.points}}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="card mb">
      <div class="card-header">
        <h4 class="title is-4">Overall Ranking</h4>
      </div>
      <div class="card-content">
        <table class="table">
          <thead>
            <tr>
              <th>
                <abbr title="Position">Pos</abbr>
              </th>
              <th>Player</th>
              <th>Won</th>
              <th>Played</th>
            </tr>
          </thead>

          <tbody>
            <tr v-for="(player,i) in overallranking" :key="player.id">
              <td>{{i+1}}</td>
              <td>{{player.name}}</td>
              <td>{{player.winners}}</td>
              <td>{{player.gamesPlayed}}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data: () => ({
    players: [],
    rounds: [],
    backendUrl: "localhost:8787",
    numbersSelected: [],
    message: null,
    leaderboard: [],
    overallranking: [],
    gameRunning: false,
    winner: null,
    game: null,
  }),
  computed: {},
  created() {
    const ws = new WebSocket("ws://" + this.backendUrl + "/api/ws");

    ws.onopen = function () {
      console.log("Connected to WS");
    };

    ws.onmessage = (evt) => {
      // console.log(evt);
      const msg = evt.data;
      // console.log(msg);
      try {
        const parsedMsg = JSON.parse(msg);
        if (!parsedMsg.type) {
          console.log("no type");
          return;
        }

        switch (parsedMsg.type) {
          case "start":
            console.log("start");
            this.gameRunning = true;
            break;
          case "round":
            this.winner = null;

            if (parsedMsg.data) {
              this.game = parsedMsg.data;
              this.gameRunning = parsedMsg.data.gameRunning;
              this.leaderboard = parsedMsg.data.players;
            }

            break;
          case "end":
            // this.rounds = [];
            // console.log(parsedMsg.data);
            this.gameRunning = false;
            this.winner = parsedMsg.data;
            break;
          case "overallranking":
            this.overallranking = parsedMsg.data;
            break;
        }
      } catch {
        console.log("failed to parse json");
      }
    };
  },
  methods: {},
};
</script>


<style scoped>
.mb {
  margin-bottom: 20px;
}
</style>
