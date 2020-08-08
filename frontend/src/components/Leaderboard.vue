<template>
  <div>
    <h2 class="subtitle">Leaderboard</h2>
   
    <div v-for="(round,i) in rounds" :key="i">{{round}}</div>

    <div v-for="player in leaderboard" :key="player.id">{{player.name}} {{player.points}}</div>
  </div>
</template>

<script>

export default {
  data: () => ({
    ws: null,
    players: [],
    rounds: [],
    backendUrl: "localhost:8787",
    numbersSelected: [],
    message: null,
    leaderboard: [],
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
            this.leaderboard = parsedMsg.data.players;
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
  },
  methods: {
 
  },
};
</script>


<style scoped>
</style>
