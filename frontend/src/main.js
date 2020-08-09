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
    connected: false,
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

      return Vue.$cookies.get("betting_game_player_id");
    },
    player(state) {
      if (state.player){
        return state.player;
      }

      const playerCookie = Vue.$cookies.get("betting_game_player");
      console.log("Playecookie", playerCookie)
      // if(playerCookie){
      //   return JSON.parse(playerCookie)
      // }

      return playerCookie

      // return null;
    },
    connected(state){
      return state.connected;
    }
  },
  mutations: {
    setConnected(state, value){
      state.connected = value;
    },
    setPlayer(state, player){
      state.player = player;
      Vue.$cookies.set("betting_game_player", player);
    },
    setPlayerId(state,id){
      state.playerId = id;
      Vue.$cookies.set("betting_game_player_id", id);
    },
    removePlayerId(state){
      state.playerId = null;
      Vue.$cookies.remove("betting_game_player_id");
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
