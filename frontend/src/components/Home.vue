<template>
  <div class="hello">
    Websocker example
    <button @click="observe()">Observe</button>
    <button @click="play()">Play</button>
    <div v-if="isPlayer">
      <div>
        <label for="name">Name:</label>
        <input type="text" name="name" id="name" v-model="name" required />
      </div>
      <div v-for="i in 10" :key="i">
        <button @click="selectNumber(i)" :disabled="setNumberButtonDisabled(i)">{{i}}</button>
      </div>
      <div>
        Numbers selected:
        <div>{{firstNumber}}</div>
        <div>{{secondNumber}}</div>
      </div>
    </div>

    <div v-for="(round,i) in rounds" :key="i">
        {{round}}
    </div>

    <div v-for="player in leaderboard" :key="player.id">
        {{player.name}} {{player.points}}
    </div>
  </div>
</template>

<script>
import { getTest, postJoinGame } from "../api";

export default {
  data: () => ({
    ws: null,
    players: [],
    rounds: [],
    backendUrl: "localhost:8787",
    numbersSelected: [],
    firstNumber: null,
    secondNumber: null,
    isPlayer: false,
    name: null,
    message: null,
    leaderboard : []
  }),
  computed: {},
  created() {
    this.ws = new WebSocket("ws://" + this.backendUrl + "/api/ws");

    this.ws.onopen = function () {
      console.log("Connected to WS");
    };

    this.ws.onmessage = (evt) => {
      // console.log(evt);
      const msg = evt.data;
      // console.log(msg);
      try {
        const parsedMsg = JSON.parse(msg);
        if (!parsedMsg.type) {
          return;
        }

        switch (parsedMsg.type) {
          case "round":
            this.rounds.unshift(parsedMsg.data.roundCounter);
            this.leaderboard = parsedMsg.data.players
            break;
          case "end":
            this.rounds = [];
            console.log(parsedMsg.data);
            break;
        }
      } catch {
        console.log("failed to parse json");
      }
    };

    getTest();
  },
  methods: {
    play() {
      this.isPlayer = true;
      console.log("play");
    },
    observe() {},
    selectNumber(number) {
      if (!this.firstNumber) {
        this.firstNumber = number;
        return;
      }

      if (!this.secondNumber) {
        this.secondNumber = number;
        this.joinGame();
        return;
      }
    },
    setNumberButtonDisabled(number) {
      if (this.firstNumber == number) {
        return true;
      }

      if (this.secondNumber == number) {
        return true;
      }

      return false;
    },
    async joinGame() {
      const payload = {
        name: this.name,
        numbers: [this.firstNumber, this.secondNumber],
      };
      await postJoinGame(payload);
    },
  },
};
</script>


<style scoped>
</style>
