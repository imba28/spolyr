<template>
  <aside
    class="position-sticky"
    style="stop: 1em"
  >
    <ul class="nav flex-column">
      <li class="nav-item">
        <router-link
          :to="{name:'home'}"
          class="nav-link"
        >
          <i class="fas fa-home" /> Home
        </router-link>
      </li>
      <li
        v-if="authStore.isAuthenticated"
        class="nav-item"
      >
        <router-link
          :to="{name:'dashboard'}"
          class="nav-link"
        >
          <i class="fas fa-music" /> Dashboard
        </router-link>
      </li>
    </ul>
  </aside>
</template>

<script>
import {mapStores} from 'pinia';
import {useAuthStore} from '@/stores/auth';

export default {
  computed: {
    ...mapStores(useAuthStore),
  },
};

</script>
