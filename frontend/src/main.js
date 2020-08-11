import Vue from 'vue'
import App from './views/Main.vue'

import VueRouter from "vue-router";
Vue.use(VueRouter);
import routes from "./routes";


import "bulma/css/bulma.min.css"

import Storage from 'vue-ls';

const storageOptions = {
  namespace: 'vuejs__', // key prefix
  name: '$session', // name variable Vue.[ls] or this.[$ls],
  storage: 'session', // storage name session, local, memory
};

Vue.use(Storage, storageOptions);

import Vuex from 'vuex'
Vue.use(Vuex)

const store = new Vuex.Store({
  state: {
    connected: false,
    player: null,
    showStartButton: true
  },
  getters:{
    player(state) {
      return state.player || Vue.$session.get("betting_game_player");
    },
    connected(state){
      return state.connected;
    },
    showStartButton(state){
      return state.showStartButton;
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
    },
    setShowStartButton(state, value){
      state.showStartButton = value;
    },
  }
})

Vue.config.productionTip = false

const router = new VueRouter({
  routes,
  mode: "history",
  base: "/",
});

router.beforeEach((to, from, next) => {
  if(to.fullPath === "/"){
    store.state.showStartButton = false;
  } else {
    store.state.showStartButton = true;
  }
  next()
})

Vue.config.productionTip = false;

new Vue({
  store,
  router,
  render: h => h(App),
}).$mount('#app')
