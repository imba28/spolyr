import Vue from 'vue';
import App from './App.vue';

import './scss/main.scss';

import './icons';
import './plugins/bootstrap-vue';
import router from './plugins/router';

Vue.config.productionTip = false;


new Vue({
  router,

  render: (h) => h(App),
}).$mount('#app');
