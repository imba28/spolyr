import Vue from 'vue';
import App from './App.vue';

import './scss/main.scss';

import './icons';
import {router, pinia} from './plugins';

Vue.config.productionTip = false;

import {ApiClient} from './openapi';
ApiClient.instance.enableCookies = true;

new Vue({
  router,
  pinia,

  render: (h) => h(App),
}).$mount('#app');
