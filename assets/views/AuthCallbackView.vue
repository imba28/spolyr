<template>
  <b-container>
    Signing in...
    <b-spinner variant="primary" />
  </b-container>
</template>

<script>
import querystring from 'querystring';
import {mapStores} from 'pinia';
import {useAuthStore} from '@/stores/auth';

export default {
  computed: {
    ...mapStores(useAuthStore),
  },
  async mounted() {
    try {
      const params = querystring.parse(window.location.search.substring(1));
      if (!params.code) {
        alert('ERROR');
      }

      await this.authStore.login(params.code);
      this.$router.push({name: 'home'});
    } catch (e) {
      console.error(e);
      this.$router.push({name: 'home'});
    }
  },
};
</script>
