import Vue from "vue"
import Vuex from 'vuex'
Vue.use(Vuex)

const store = new Vuex.Store({
  state: {
    showStartButton: true,
    player: null,
    socket: {
      isConnected: false,
      message: '',
      reconnectError: false,
    }
  },
  getters: {
    showStartButton(state) {
      return state.showStartButton;
    },
    player(state) {
      return  state.player || Vue.$session.get("betting_game_player");
    },
    socket(state) {
      return state.socket;
    },
  },
  mutations: {
    setPlayer(state, player) {
      state.player = player;
      Vue.$session.set("betting_game_player", player);
    },
    setNotObserver(state) {
      state.player = { ...state.player, observer: false }
      Vue.$session.set("betting_game_player", state.player);
    },
    removePlayer(state) {
      state.player = null;
      Vue.$session.remove("betting_game_player");
    },
    SOCKET_ONOPEN(state, event) {
      Vue.prototype.$socket = event.currentTarget
      state.socket.isConnected = true
    },
    SOCKET_ONCLOSE(state) {
      state.socket.isConnected = false

      state.player = null;
      Vue.$session.remove("betting_game_player", null);
    },
    SOCKET_ONERROR(state, event) {
      console.error(state, event)
      state.player = null;
      Vue.$session.remove("betting_game_player", null);
    },
    // default handler called for all methods
    SOCKET_ONMESSAGE(state, msg) {
      state.socket.message = msg;
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT(state, count) {
      console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR(state) {
      state.socket.reconnectError = true;
    },
  },

})



export default store;