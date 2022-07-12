import Vue from 'vue';
import VueRouter from 'vue-router';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

import HomeView from '../views/HomeView.vue';
import {useAuthStore} from '@/stores/auth';

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
  {
    path: '/auth/callback',
    name: 'auth-callback',
    component: () => import('../views/AuthCallbackView.vue'),
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: () => import('../views/DashboardView.vue'),
    meta: {
      authRequired: false,
    },
  },
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

router.beforeEach((to, _, next) => {
  const authStore = useAuthStore();
  if (to.meta.authRequired && !authStore.isAuthenticated) {
    router.app.$toast.info('This pages requires you to be signed in!', {
      timeout: 2000,
    });
    return {name: 'home'};
  }

  next();
});

NProgress.configure({showSpinner: false});
router.beforeResolve((to, _, next) => {
  if (to.name) {
    NProgress.start();
  }
  next();
});
router.afterEach(() => {
  NProgress.done();
});

export default router;
