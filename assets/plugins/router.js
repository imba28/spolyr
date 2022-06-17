import Vue from 'vue';
import VueRouter from 'vue-router';
import HomeView from '../views/HomeView.vue';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
  },
  {
    path: '/search/:q',
    name: 'search',
    component: () => import('../views/SearchView.vue'),
  },
  {
    path: '/tracks/:id',
    name: 'track-detail',
    component: () => import('../views/TrackDetailView.vue'),
  },
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

export default router;
