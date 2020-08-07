<template>
  <div class="hello">
    Websocker example
    <button @click="observe()">Observe</button>
    <button @click="play()">Play</button>
    <div v-if="isPlayer">
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
  name: "HelloWorld",
  data: () => ({
    ws: null, // Our websocket
    // newMsg: "", // Holds new messages to be sent to the server
    // chatContent: "", // A running list of chat messages displayed on the screen
    // email: null, // Email address used for grabbing an avatar
    // username: null, // Our username
    // joined: false, // True if email and username have been filled in
    players: [],
    rounds: [],
    backendUrl: "localhost:8787",
    numbersSelected: [],
    firstNumber: null,
    secondNumber: null,
    isPlayer: false,
  }),
  computed: {
    // setNumberButtonDisabled(number) {
    //   if (this.firstNumber == number){
    //       return true
    //   }
    //   return false;
    // },
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
        //start game
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
      await postJoinGame({ numbers: [this.firstNumber, this.secondNumber] });
    },
  },
};
</script>


<style scoped>
</style>
