<template>
  <div id="main-app">
    <nav class="navbar" role="navigation" aria-label="main navigation">
      <!-- <div class="navbar-brand"> -->
      <div class="navbar-start">
        <router-link to="/" tag="a" class="navbar-item">
          <h1 class="title">Betting Game</h1>
        </router-link>
      </div>
      <div class="navbar-end">
        <div class="navbar-item" v-if="player">
          <p>
            <strong>{{player.name}}</strong>
          </p>
          <p v-if="player.observer">
            <small>(observing)</small>
          </p>
        </div>
        <div class="navbar-item">
          <div class="buttons">
            <router-link id="start-btn" v-if="showStartButton && !player" to="/" tag="a" class="button">
              <strong>Start</strong>
            </router-link>

            <router-link id="rules-btn" to="/rules" class="button is-light">Rules</router-link>
          </div>
        </div>
        <div class="navbar-item">
          <div id="connection-status">
            <p v-if="connected" class="help is-success">Connected</p>
            <p v-else class="help is-danger">Not connected</p>
          </div>
        </div>
      </div>
    </nav>
    <div class="app">
      <div>
        <router-view></router-view>
      </div>
    </div>
    <footer id="footer">
      <div id="author-link">
        <p>
          <a href="https://github.com/evzpav">Evandro Pavei</a> - Florianópolis/Brazil - 2020
        </p>
      </div>
    </footer>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
export default {
  computed: {
    ...mapGetters(["player", "connected","showStartButton"]),
  },
};
</script>

<style>
#main-app {
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

.app {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-content: center;
  align-items: center;
  margin: auto;
}

#footer {
  width: 100%;
  text-align: center;
  padding: 10px;
  display: flex;
  justify-content: center;
  justify-items: center;
  margin-top: 30px;
}

#footer p {
  font-size: 12px;
}

#connection-status p {
  font-size: 10px;
}

#start-btn {
  background-color: #fffe03;
  border: none;
}
</style>
