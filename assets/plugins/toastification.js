import Vue from 'vue';
import Toast from 'vue-toastification';
import 'vue-toastification/dist/index.css';

const options = {
  position: 'top-right',
  maxToasts: 5,
  hideProgressBar: true,
};

Vue.use(Toast, options);
