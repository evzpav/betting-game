<template>
  <div>
    <div class="content">
      <div class="card game-table">
        <div
          v-if="player && player.observer && gameStarted"
          class="notification is-info"
        >You will automatically join in the next game.</div>

        <div
          v-if="!gameStarted"
          class="notification is-info"
        >Waiting for players to join.</div>

        <div v-if="isGameRunning && !game.winner">
          <h4 class="title is-4">Game #{{game.gameCounter}}</h4>
        </div>
        <div v-if="isLoading">Loading game...</div>

        <div v-if="game" class="card-content">
          <div v-if="isGameRunning && !game.winner" class="notification is-black">
            <progress
              class="progress"
              :value="game.roundCounter"
              :max="game.rules.maxRoundsPerGame"
            ></progress>
            <div>Round: {{game.roundCounter}}/{{game.rules.maxRoundsPerGame}}</div>
            <div>
              Number:
              <strong>{{game.randomNumber}}</strong>
            </div>
          </div>

          <div v-if="game.winner" class="notification is-black winner-notification">
            <div>
              <div v-if="isPlayerTheWinner()">
                <p>
                  Congratulations
                  <strong>{{game.winner.name}}</strong>!
                </p>
                <p>YOU won game #{{game.gameCounter}} with {{game.winner.points}} points!</p>
              </div>
              <div v-else>
                <div>
                  <strong>{{game.winner.name}}</strong>
                  is the winner of game #{{game.gameCounter}} with {{game.winner.points}} points.
                </div>
              </div>
            </div>
            <!-- <img class="trophy" :src="trophy" alt="trophy" width="20" height="40" /> -->
          </div>

          <div v-if="game.winner" class="notification is-info">
            <p>New players will be joining now. New game commencing in few seconds...</p>
            <progress class="progress is-small is-info" max="100"></progress>
          </div>

          <table class="table" v-if="game.gameRunning">
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
              <tr
                v-for="(player,i) in leaderboard"
                :key="player.id"
                :class="highlightPlayer(player.id)"
              >
                <td>{{i+1}}</td>
                <td>{{player.name}}</td>
                <td>
                  <span class="tag is-black number-tag">{{player.numbers[0]}}</span>
                  <span class="tag is-black number-tag">{{player.numbers[1]}}</span>
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
        <div v-if="isLoading">Loading ranking...</div>
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
// import trophy from "../assets/images/trophy1.svg";
import { newWebsocket, getRankingSnapshot, getGameSnapshot } from "../api";
import { mapGetters } from "vuex";

export default {
  data: () => ({
    leaderboard: [],
    overallranking: [],
    game: null,
    newGameStarting: false,
    // trophy: trophy,
    isLoading: false,
  }),
  computed: {
    ...mapGetters(["player"]),
    gameStarted(){
      return this.game && this.game.gameCounter > 0;
    },
    isGameRunning(){
      return this.game && this.game.gameRunning
    }
  },
  created() {
    this.loadRankingSnapshot();
    this.loadGameSnapshot();

    const ws = newWebsocket();

    ws.onopen = () => {
      console.log("Connected to WS");
      this.$store.commit("setConnected", true);
    };

    ws.onerror = () => {
      console.log("Cannot connect to WS");
      this.clearData();
    };

    ws.onclose = () => {
      this.clearData();
    };

    ws.onmessage = (evt) => {
      const msg = evt.data;

      try {
        const parsedMsg = JSON.parse(msg);
        if (!parsedMsg.type || !parsedMsg.data) {
          throw "failed to parse message";
        }

        switch (parsedMsg.type) {
          case "start":
            this.game = parsedMsg.data;
            this.leaderboard = this.game.players;
            break;
          case "round":
            this.game = parsedMsg.data;
            this.leaderboard = parsedMsg.data.players;
            this.isNotObserver(this.leaderboard);

            break;
          // case "restart":
          //   console.log("restart");
          //   this.game = parsedMsg.data;
          //   console.log(this.game)
          //   this.leaderboard = this.game.players;
          //   break;
          case "end":
            this.game = parsedMsg.data;
            this.leaderboard = this.game.players;
            console.log(this.game)
            break;
          case "overallranking":
            this.overallranking = parsedMsg.data;
            break;
        }
      } catch (e) {
        console.log("catch: ", e);
      }
    };
  },
  methods: {
    isNotObserver(players) {
      const meObserver = players.find((player) => {
        return (
          this.player && this.player.id === player.id && this.player.observer
        );
      });

      if (meObserver && !meObserver.observer) {
        this.$store.commit("setPlayer", meObserver);
      }
    },
    isPlayerTheWinner() {
      return this.player && this.game && this.player.id === this.game.winner.id;
    },
    async loadRankingSnapshot() {
      this.isLoading = true;
      try {
        const resp = await getRankingSnapshot();
        this.overallranking = resp.data ? resp.data : [];
      } catch (error) {
        console.log(error);
      } finally {
        this.isLoading = false;
      }
    },

    async loadGameSnapshot() {
      this.isLoading = true;
      try {
        const resp = await getGameSnapshot();
        this.game = resp && resp.data ? resp.data : [];
        this.leaderboard = this.game.players;
      } catch (error) {
        console.log(error);
      } finally {
        this.isLoading = false;
      }
    },
    highlightPlayer(id) {
      return this.player && this.player.id === id ? "is-selected" : "";
    },
    clearData() {
      this.$store.commit("setConnected", false);
      this.$store.commit("setPlayer", null);
      this.leaderboard = [];
      this.overallranking = [];
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

.table tr.is-selected {
  background-color: #fffe03;
  color: black;
}

.number-tag {
  margin: 1px;
}
</style>
