<template>
  <div class="container">
    <div v-if="!isPlayer" class="action-buttons">
      <div id="play-btn" class="action-button play" @click="play()">Play</div>
      <div id="observe-btn" class="action-button observe" @click="observe()">Observe</div>
    </div>
    <div v-if="isPlayer" class="field">
      <label class="label">Name</label>
      <div class="control mb">
        <input
          id="name-input"
          v-model="name"
          class="input"
          :class="{'is-danger': nameError }"
          type="text"
          placeholder="john"
          maxlength="30"
          @input="onInputChange()"
        />
        <p class="help is-danger" v-if="nameError">{{nameError}}</p>
      </div>

      <label class="label">Choose your lucky numbers. Select 2:</label>
      <div class="numbers mb">
        <button
          :id="genButtonId(i)"
          type="button"
          class="button is-black number-item"
          v-for="i in 10"
          :key="i"
          @click="selectNumber(i)"
          :disabled="setNumberButtonDisabled(i)"
        >{{i}}</button>
      </div>
      <p class="help is-danger" v-if="numberError">{{numberError}}</p>
      <div class="buttons">
        <button id="join-btn" class="button is-success" :class="{'is-loading': isLoading}" @click="joinGame">JOIN</button>
        <button id="cancel-btn" class="button is-secondary" @click="cancel()">CANCEL</button>
      </div>
      <p class="help is-danger" v-if="error">Could not proceed. Please contact support.</p>
    </div>
  </div>
</template>

<script>
import { postJoinGame } from "../api";

export default {
  data: () => ({
    name: "",
    isPlayer: false,
    firstNumber: null,
    secondNumber: null,
    isLoading: false,
    error: false,
    nameError: "",
    numberError: "",
  }),
  created() {
    if (this.player) {
      this.$router.push("leaderboard");
    }
  },
  methods: {
    genButtonId(i){
      return `number-btn-${i}`
    },
    play() {
      this.isPlayer = true;
    },
    cancel() {
      this.isPlayer = false;
      this.name = "";
      this.firstNumber = null;
      this.secondNumber = null;
      this.clearError();
    },
    observe() {
      this.$store.commit("removePlayer");
      this.$router.push("leaderboard");
    },
    selectNumber(number) {
      this.numberError = "";
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
      if (!this.name || this.name.length < 3) {
        this.nameError = "Invalid name. Minimum 3 characters";
        return false;
      }

      if (!this.firstNumber || !this.secondNumber) {
        this.numberError = "Please select 2 numbers";
        return false;
      }

      return true;
    },
    onInputChange() {
      this.name = this.name.toLowerCase();
      this.clearError();
    },
    clearError() {
      this.nameError = "";
      this.numberError = "";
    },
    async joinGame() {
      if (!this.validateInputs()) {
        return false;
      }

      const numbers = [this.firstNumber, this.secondNumber];
      numbers.sort();

      this.name = this.name.toLowerCase();

      const payload = {
        name: this.name,
        numbers: numbers,
      };

      try {
        this.isLoading = true;
        const resp = await postJoinGame(payload);

        if (resp && resp.data) {
          const player = resp.data;
          this.$store.commit("setPlayer", player);
        }

        this.$router.push("leaderboard");
      } catch (error) {
        if (error && error.response && error.response.status === 400) {
          this.nameError = error.response.data.message;
          return;
        }

        this.error = true;
        console.log(error);
      } finally {
        this.isLoading = false;
      }
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
  height: 60vh;
  padding: 15px;
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
