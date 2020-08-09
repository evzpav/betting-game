<template>
  <div>
    <nav class="navbar" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <router-link to="/" tag="a" class="navbar-item">
          <h1 class="title">Betting Game</h1>
        </router-link>
      </div>
      <div class="navbar-end">
        <div class="navbar-item" v-if="player">{{player.name}}</div>
        <div class="navbar-item">
          <div class="buttons">
            <router-link v-if="!player" to="/" tag="a" class="button is-warning">
              <strong>Start</strong>
            </router-link>

            <router-link to="/rules" class="button is-light">Rules</router-link>
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
      <div id="connection-status" class="has-text-right">
        <p v-if="connected" class="help is-success">Connected</p>
        <p v-else class="help is-danger">Not connected</p>
      </div>
      <div id="author-link" class="content has-text-centered">
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
    ...mapGetters(["player", "connected"]),
  },
};
</script>

<style>
.app {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-content: center;
  align-items: center;
}

#footer {
  position: fixed;
  bottom: 0;
  width: 100%;
  height: 40px;
}

#footer p {
  font-size: 12px;
}

#main-title {
  padding: 10px;
  margin-bottom: 10vh;
}

#connection-status {
}
</style>
