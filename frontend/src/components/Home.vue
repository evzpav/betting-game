<template>
  <div class="container">
    <div v-if="!isPlayer" class="action-buttons">
      <div class="action-button play" @click="play()">Play</div>
      <div class="action-button observe" @click="observe()">Observe</div>
    </div>
    <div v-if="isPlayer" class="field">
      <label class="label">Name</label>
      <div class="control mb">
        <input v-model="name" class="input" type="text" placeholder="John Smith" />
      </div>

      <label class="label">Choose your lucky numbers. Select 2:</label>
      <div class="numbers mb">
        <button
          type="button"
          class="button is-black number-item"
          v-for="i in 10"
          :key="i"
          @click="selectNumber(i)"
          :disabled="setNumberButtonDisabled(i)"
        >{{i}}</button>
      </div>
      <button class="button is-primary" :class="{'is-loading': isLoading}" @click="joinGame">Join</button>
      <div v-if="error">Could not proceed. Please contact support</div>
    </div>
  </div>
</template>

<script>
import { postJoinGame } from "../api";

export default {
  data: () => ({
    name: "",
    playerId: "",
    isPlayer: false,
    firstNumber: null,
    secondNumber: null,
    isLoading: false,
    error: false,
  }),
  computed: {},

  methods: {
    play() {
      this.isPlayer = true;
    },
    observe() {
      this.removeCookiePlayerId();
      this.$router.push("leaderboard");
    },
    selectNumber(number) {
      if (!this.firstNumber) {
        this.firstNumber = number;
        return;
      }

      if (!this.secondNumber) {
        this.secondNumber = number;
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
    validateInputs() {
      this.error = false;
      if (!this.name || this.name.length < 4) {
        return false;
      }

      if (!this.firstNumber || !this.secondNumber) {
        return false;
      }

      return true;
    },
    async joinGame() {
      if (!this.validateInputs()) {
        alert("invaldi");
        return false;
      }

      const numbers = [this.firstNumber, this.secondNumber];
      numbers.sort();

      const payload = {
        name: this.name,
        numbers: numbers,
      };

      try {
        this.isLoading = true;
        const resp = await postJoinGame(payload);

        if (resp && resp.data) {
          this.playerId = resp.data.id;
          this.setCookiePlayerId(this.playerId);
        }

        this.$router.push("leaderboard");
      } catch (e) {
        this.error = true;
        console.error(e);
      } finally {
        this.isLoading = false;
      }
    },
    getCookiePlayerId() {
      return this.$cookies.get("betting_game_player");
    },
    setCookiePlayerId(id) {
      return this.$cookies.set("betting_game_player", id);
    },
    removeCookiePlayerId() {
      this.$cookies.remove("betting_game_player");
    },
  },
};
</script>


<style scoped>
.mb {
  margin-bottom: 20px;
}

.container {
  display: flex;
  margin-bottom: 20px;
}

.action-buttons {
  display: flex;
  margin-top: 10vh;
}

.action-button {
  color: white;
  text-transform: uppercase;
  font-size: 1rem;
  font-weight: 500;
  width: 10vw;
  min-width: 100px;
  height: 10vh;
  padding: 1.25rem;
  display: flex;
  justify-content: center;
  align-items: center;
  margin: 15px;
  cursor: pointer;
  border-radius: 6px;
  box-shadow: 0 0.5em 1em -0.125em rgba(10, 10, 10, 0.1),
    0 0 0 1px rgba(10, 10, 10, 0.02);
}

.play {
  background-color: #fffe03;
  color: black;
}

.observe {
  background-color: black;
  color: #fffe03;
}

.numbers {
  display: flex;
  flex-wrap: wrap;
}

.number-item {
  margin: 2px;
  width: 10px;
}
</style>  
