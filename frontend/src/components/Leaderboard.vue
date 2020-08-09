<template>
  <div>

    <div class="content">
      <div class="card game-table">
        <div v-if="game">
          <h4 class="title is-4">Current Game #{{game.gameCounter}}</h4>
        </div>
        <div class="card-content">
          <div v-if="!started">Game stopped. Waiting for players to join.</div>

          <div v-if="gameRunning" class="notification is-warning">
            <div>Round: {{game.roundCounter}}</div>
            <div>Number: {{game.randomNumber}}</div>
          </div>

          <div v-if="winner" class="notification is-success winner-notification">
            <div>Winner is: {{winner.name}}</div>
            <img class="trophy" :src="trophy" alt="trophy" width="20" height="40" />
          </div>

          <table class="table" v-if="gameRunning">
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
              <tr v-for="(player,i) in leaderboard" :key="player.id" :class="highlightPlayer(player.id)">
                <td>{{i+1}}</td>
                <td>{{player.name}}</td>
                <td>
                  <span class="tag is-black">{{player.numbers[0]}}</span>
                  <span class="tag is-black">{{player.numbers[1]}}</span>
                </td>
                <td>{{player.points}}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="card ranking-table">
        <div>
          <h4 class="title is-4">Overall Ranking</h4>
        </div>
        <div class="card-content" v-if="overallranking.length > 0">
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
              <tr
                v-for="(player,i) in overallranking"
                :key="player.id"
                :class="highlightPlayer(player.id)"
              >
                <td>{{i+1}}</td>
                <td>{{player.name}}</td>
                <td>{{player.winners}}</td>
                <td>{{player.gamesPlayed}}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else>No games played yet</div>
      </div>
    </div>
  </div>
</template>

<script>
import trophy from "../assets/images/trophy1.svg";

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
    started: false,
    winner: null,
    game: null,
    trophy: trophy,
    playerId: null,
  }),
  computed: {},
  created() {
    this.playerId = this.getCookiePlayerId();

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
        if (!parsedMsg.type || !parsedMsg.data) {
          console.log(parsedMsg);
          throw "failed to parse message";
        }

        switch (parsedMsg.type) {
          case "start":
            console.log("start");
            this.started = true;
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
            this.gameRunning = false;
            this.winner = parsedMsg.data;
            break;
          case "overallranking":
            this.overallranking = parsedMsg.data;
            break;
        }
      } catch {
        console.log("catch");
      }
    };
  },
  methods: {
    highlightPlayer(id) {
      return this.playerId === id ? "is-selected" : "";
    },
    getCookiePlayerId() {
      return this.$cookies.get("betting_game_player");
    },
  },
};
</script>


<style scoped>
.content {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
}

.game-table {
  margin-bottom: 20px;
  padding: 10px;
  flex-grow: 1;
  margin-right: 20px;
  min-width: 35vw;
}

.ranking-table {
  margin-bottom: 20px;
  padding: 10px;
  flex-grow: 1;
  min-width: 20vw;
}

.winner-notification {
  display: flex;
  justify-content: space-between;
}

.trophy {
  color: #fffe03;
}
</style>
