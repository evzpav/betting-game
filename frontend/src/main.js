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

import store from "./store"
import {resolveWebsocketUrl} from "./api"

import VueNativeSock from 'vue-native-websocket'
Vue.use(VueNativeSock, resolveWebsocketUrl(), { 
  store: store, 
  format: 'json',
  connectManually: true
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
