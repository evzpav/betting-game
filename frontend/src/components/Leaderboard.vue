<template>
  <div>
    <div class="content">
      <div class="card game-table">
        <div v-if="game">
          <h4 class="title is-4">Current Game #{{game.gameCounter}}</h4>
        </div>
        <div v-if="isLoading">Loading game...</div>
        <div v-if="game" class="card-content">
          <div v-if="!game.gameRunning">Game stopped. Waiting for players to join.</div>

          <div v-if="game.gameRunning && !winner" class="notification is-info">
            <div>Round: {{game.roundCounter}}</div>
            <div>Number: {{game.randomNumber}}</div>
          </div>

          <div v-if="winner" class="notification is-info winner-notification">
            <div>Winner is: {{winner.name}}!</div>
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
import { mapState, mapGetters } from "vuex";

export default {
  data: () => ({
    leaderboard: [],
    overallranking: [],
    winner: null,
    game: null,
    trophy: trophy,
    isLoading: false,
  }),
  computed: {
    ...mapState(["player",  "gameStarted"]),
    ...mapGetters(["playerId"])
  },
  created() {
    console.log(this.playerId)
    this.loadRankingSnapshot();
    this.loadGameSnapshot();

    const ws = newWebsocket();

    ws.onopen = () => {
      console.log("Connected to WS");
    };

    ws.onerror = function () {
      console.log("Cannot connect to WS");
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
            this.$store.commit("setGameStarted");
            break;
          case "round":
            this.winner = null;
            if (parsedMsg.data) {
              this.game = parsedMsg.data;
              this.leaderboard = parsedMsg.data.players;
            }

            break;
          case "end":
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
    async loadRankingSnapshot() {
      this.isLoading = true;
      try {
        const resp = await getRankingSnapshot();
        console.log(resp.data);
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
        console.log(resp.data);
        this.game = resp.data ? resp.data : [];
      } catch (error) {
        console.log(error);
      } finally {
        this.isLoading = false;
      }
    },
    highlightPlayer(id) {
      return this.playerId === id ? "is-selected" : "";
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
