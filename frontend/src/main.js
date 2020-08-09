import Vue from 'vue'
import App from './views/Main.vue'

import VueRouter from "vue-router";
Vue.use(VueRouter);
import routes from "./routes";

Vue.config.productionTip = false

import "bulma/css/bulma.min.css"

import VueCookies from 'vue-cookies'
Vue.use(VueCookies)

import Vuex from 'vuex'
Vue.use(Vuex)

const store = new Vuex.Store({
  state: {
    player: null,
    playerId: null,
    gameRunning: false,
    gameStarted: false,
  },
  getters:{
    playerId(state) {
      if (state.playerId){
        return state.playerId;
      }

      return Vue.$cookies.get("betting_game_player");
    },
  },
  mutations: {
    setPlayer(state, player){
      state.player = player;
    },
    setPlayerId(state,id){
      state.playerId = id;
      Vue.$cookies.set("betting_game_player", id);
    },
    removePlayerId(state){
      state.playerId = null;
      Vue.$cookies.remove("betting_game_player");
    },
    setGameStarted(state){
      state.gameStarted = true;
    },
  }
})


const router = new VueRouter({
  routes,
  mode: "history",
  base: "/",
});

Vue.config.productionTip = false;

new Vue({
  store,
  router,
  render: h => h(App),
}).$mount('#app')
