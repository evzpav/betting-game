import Vue from 'vue'
import App from './views/Main.vue'

import VueRouter from "vue-router";
Vue.use(VueRouter);
import routes from "./routes";

Vue.config.productionTip = false

import "bulma/css/bulma.min.css"

import VueCookies from 'vue-cookies'
Vue.use(VueCookies)

const router = new VueRouter({
  routes,
  mode: "history",
  base: "/",
});

Vue.config.productionTip = false;

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
