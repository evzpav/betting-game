<template>
  <div class="hello">
    Websocker example
    <button @click="observe()">Observe</button>
    <button @click="play()">Play</button>
    <div v-if="isPlayer">
      <div>
        <label for="name">Name:</label>
        <input type="text" name="name" id="name" v-model="name" required>
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

    <!-- <div v-for="(round,i) in rounds" :key="i">
        {{round}}
    </div>-->
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
    name: null
  }),
  computed: {
  },
  created() {
    this.ws = new WebSocket("ws://" + this.backendUrl + "/api/ws");

    this.ws.onopen = function () {
      console.log("Connected to WS");
    };

    this.ws.onmessage = (evt) => {
      console.log(evt);
      this.rounds.unshift(evt.data);
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
