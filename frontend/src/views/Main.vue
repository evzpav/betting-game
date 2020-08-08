<template>
  <div>
    <div class="app">
      <div id="main-title">
        <h1 class="title">Betting Game</h1>
      </div>
      <div>
        <router-view></router-view>
      </div>
    </div>
    <footer id="footer">
      <div class="content has-text-centered">
        <p>
          <!-- <a href="https://github.com/evzpav">Evandro Pavei</a> - Florianópolis/Brazil - 2020 -->
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

#main-title {
  padding: 10px;
  margin-bottom: 10vh;
}
</style>
