import Vue from 'vue'
import App from './views/Main.vue'

import VueRouter from "vue-router";
Vue.use(VueRouter);
import routes from "./routes";

Vue.config.productionTip = false

import "bulma/css/bulma.min.css"

import VueCookies from 'vue-cookies'
Vue.use(VueCookies)

import Storage from 'vue-ls';
 
const options = {
  namespace: 'vuejs__', // key prefix
  name: '$session', // name variable Vue.[ls] or this.[$ls],
  storage: 'session', // storage name session, local, memory
};
 
Vue.use(Storage, options);

import Vuex from 'vuex'
Vue.use(Vuex)

const store = new Vuex.Store({
  state: {
    connected: false,
    player: null,
    gameStarted: false,
  },
  getters:{
    gameStarted(state){
      return state.gameStarted;
    },
    player(state) {
      if (state.player){
        return state.player;
      }
      console.log(Vue.$session.get("betting_game_player"))
      return Vue.$session.get("betting_game_player");
     
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
      Vue.$session.set("betting_game_player", player);
    },
    removePlayer(state){
      state.player = null;
      Vue.$session.remove("betting_game_player");
      // Vue.$session.clear()
      // Vue.$session.destroy()

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
