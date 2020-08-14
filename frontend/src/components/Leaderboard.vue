<template>
  <div class="content">
    <div class="card game-table">
      <div
        v-if="player && player.observer && gameStarted"
        class="notification is-info"
      >You will automatically join in the next game.</div>

      <div v-if="!gameStarted" class="notification is-info">Waiting for players to join.</div>

      <div v-if="isGameRunning && !game.winner">
        <h4 class="title is-4">Game #{{game.gameCounter}}</h4>
      </div>
      <div v-if="isLoading">Loading game...</div>

      <div v-if="game" class="card-content">
        <div v-if="isGameRunning && !game.winner" class="notification is-black">
          <div>Round: {{game.roundCounter}}/{{game.rules.maxRoundsPerGame}}</div>
          <progress class="progress" :value="game.roundCounter" :max="game.rules.maxRoundsPerGame"></progress>
          <div>
            Drawn Number:
            <strong class="randomNumber">{{game.randomNumber}}</strong>
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
        </div>

        <div v-if="game.winner && secondsToNextGame !== null" class="notification is-info">
          <p>New players will be joining now.</p>
          <p>New game commencing {{secondsToNextGame}} in seconds...</p>
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
              v-for="(player,i) in game.players"
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
      <div class="card-content" v-if="overallRanking && overallRanking.length > 0">
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
              v-for="(player,i) in overallRanking"
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
</template>

<script>
import { getRankingSnapshot, getGameSnapshot } from "../api";
import { mapGetters } from "vuex";

export default {
  data: () => ({
    isLoading: false,
    game: null,
    overallRanking: [],
    secondsToNextGame: null,
  }),
  computed: {
    ...mapGetters(["socket", "player"]),
    gameStarted() {
      return this.game && this.game.gameCounter > 0;
    },
    isGameRunning() {
      return this.game && this.game.gameRunning;
    },
  },
  created() {
    this.loadRankingSnapshot();
    this.loadGameSnapshot();

    if (!this.socket.isConnected) {
      this.$connect();
    }
  },
  methods: {
    isPlayerTheWinner() {
      return this.player && this.game && this.player.id === this.game.winner.id;
    },
    async loadRankingSnapshot() {
      this.isLoading = true;
      try {
        const resp = await getRankingSnapshot();
        this.overallRanking = resp.data ? resp.data : [];
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
      } catch (error) {
        console.log(error);
      } finally {
        this.isLoading = false;
      }
    },
    highlightPlayer(id) {
      return this.player && this.player.id === id ? "is-selected" : "";
    },
    setNotObserver(player, game) {
      const meObserver = game.players.find((p) => {
        return player.id === p.id;
      });

      if (meObserver && !meObserver.observer) {
        this.$store.commit("setPlayer", meObserver);
      }
    },
    closeWebsocket() {
      this.$disconnect();
    },
  },
  mounted() {
    window.addEventListener("unload", this.closeWebsocket);
  },
  beforeDestroy() {
    window.removeEventListener("unload", this.closeWebsocket);
  },
  watch: {
    "socket.message": function (oldVal, msg) {
      switch (msg.type) {
        case "start":
          this.game = msg.data;
          break;
        case "round":
          this.secondsToNextGame = null;
          this.game = msg.data;
          if (this.player && this.player.observer && this.game && this.game.players.length > 0) {
            this.setNotObserver(this.player, this.game);
          }
          break;
        case "end":
          this.game = msg.data;
          break;
        case "overallranking":
          this.overallRanking = msg.data;
          break;
        case "intervalTicker":
          this.secondsToNextGame = msg.data;
          break;
      }
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
  min-height: 40vh;
}

.ranking-table {
  margin-bottom: 20px;
  padding: 10px;
  flex-grow: 1;
  min-width: 20vw;
  min-height: 40vh;
}

.winner-notification {
  display: flex;
  justify-content: space-between;
}

.table tr.is-selected {
  background-color: #fffe03;
  color: black;
}

.number-tag {
  margin: 1px;
  font-weight: 700;
}

.randomNumber {
  font-size: 20px;
  margin-left: 10px;
}
</style>
