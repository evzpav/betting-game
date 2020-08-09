<template>
  <div>
    <nav class="navbar" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <a class="navbar-item" href="https://bulma.io">
          <h1 class="title">Betting Game</h1>
        </a>
      </div>
      <div class="navbar-end">
        <div class="navbar-item">
          <div class="buttons">
            <router-link to="/" tag="a" class="button is-warning">
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
      <div class="content has-text-centered">
        <p>
          <a href="https://github.com/evzpav">Evandro Pavei</a> - Florianópolis/Brazil - 2020
        </p>
      </div>
    </footer>
  </div>
</template>

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
  bottom: 10px;
  width: 100%;
  height: 20px;
}

#footer p {
  font-size: 12px;
}

#main-title {
  padding: 10px;
  margin-bottom: 10vh;
}
</style>
