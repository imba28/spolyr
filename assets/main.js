import Vue from 'vue';
import App from './App.vue';

import './scss/main.scss';

import './icons';
import './plugins/bootstrap-vue';
import router from './plugins/router';

Vue.config.productionTip = false;

import {ApiClient} from './openapi';
ApiClient.instance.enableCookies = true;

new Vue({
  router,

  render: (h) => h(App),
}).$mount('#app');
