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
        this.$toast.error('Authentication with Spotify failed. No code was provided!');
        this.$router.push({name: 'home'});
        return;
      }

      await this.authStore.login(params.code);
      this.$router.push({name: 'home'});
    } catch (e) {
      console.error(e);
      this.$toast.error('Something went wrong when trying to sign you in!');
      this.$router.push({name: 'home'});
    }
  },
};
</script>
