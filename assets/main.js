import Vue from 'vue';
import App from './App.vue';

import {ApiClient} from './openapi';
import './icons';
import {router, pinia} from './plugins';
import jwtRefreshPlugin from './plugins/superagent';

import './scss/main.scss';

if (process.env.NODE_ENV === 'production') {
  ApiClient.instance.basePath = '/api';
}
ApiClient.instance.enableCookies = true;
ApiClient.instance.plugins = [jwtRefreshPlugin];

Vue.config.productionTip = false;
new Vue({
  router,
  pinia,

  render: (h) => h(App),
}).$mount('#app');
