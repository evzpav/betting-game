<template>
  <div>
    <div class="content">
      <div class="card game-table">
        <div v-if="player && player.observer">You will be able to join when new game starts.</div>
        <div v-if="!game || !game.gameRunning">Waiting for players to join.</div>

        <div v-if="game && game.gameRunning">
          <h4 class="title is-4">Current Game #{{game.gameCounter}}</h4>
        </div>
        <div v-if="isLoading">Loading game...</div>

        <div v-if="game" class="card-content">
          <div v-if="game.gameRunning && !game.winner" class="notification">
            <progress
              class="progress is-info"
              :value="game.roundCounter"
              :max="game.rules.maxRoundsPerGame"
            ></progress>
            <div>Round: {{game.roundCounter}}/{{game.rules.maxRoundsPerGame}}</div>
            <div>
              Number:
              <strong>{{game.randomNumber}}</strong>
            </div>
          </div>

          <div v-if="game.winner">
            <p>New players can join now. New game commencing soon.</p>
            <progress class="progress is-small is-warning" max="100"></progress>
            <br />
          </div>

          <div v-if="game.winner" class="notification is-link winner-notification">
            <div>
              <div>Game winner is:</div>
              <div>
                <strong>{{game.winner.name}}!</strong>
              </div>
            </div>
            <img class="trophy" :src="trophy" alt="trophy" width="20" height="40" />
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
import trophy from "../assets/images/trophy1.svg";
import { newWebsocket, getRankingSnapshot, getGameSnapshot } from "../api";
import { mapGetters } from "vuex";

export default {
  data: () => ({
    leaderboard: [],
    overallranking: [],
    game: null,
    trophy: trophy,
    isLoading: false,
  }),
  computed: {
    ...mapGetters(["player", "gameStarted"]),
  },
  created() {
    this.loadRankingSnapshot();
    this.loadGameSnapshot();

    const ws = newWebsocket();

    ws.onopen = (evt) => {
      console.log("Connected to WS");
      const url = evt.target.url;

      if (url) {
        const id = url.split("id=");
        this.$store.commit("setConnected", id[1]);
        console.log(id);
      }
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
          console.log(parsedMsg);
          throw "failed to parse message";
        }

        switch (parsedMsg.type) {
          case "start":
            console.log("start");
            this.$store.commit("setGameStarted");
            break;
          case "restart":
            break;
          case "round":
            if (parsedMsg.data) {
              console.log(parsedMsg.data);
              this.game = parsedMsg.data;
              this.leaderboard = parsedMsg.data.players;

              const meObserver = this.leaderboard.find((player) => {
                return this.player.id === player.id && this.player.observer;
              });

              if (meObserver && !meObserver.observer) {
                this.$store.commit("setPlayer", meObserver);
              }
            }

            break;
          case "end":
            this.game = parsedMsg.data;
            this.leaderboard = this.game.players;
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

.number-tag {
  margin: 1px;
}
</style>
