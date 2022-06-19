<template>
  <b-container>
    Signing in...
    <b-spinner variant="primary" />
  </b-container>
</template>

<script>
import querystring from 'querystring';
import {AuthApi, AuthLoginPostRequest} from '@/openapi';

const authApi = new AuthApi();
export default {
  async mounted() {
    try {
      const params = querystring.parse(window.location.search.substring(1));
      if (!params.code) {
        alert('ERROR');
      }

      const body = AuthLoginPostRequest.constructFromObject({
        code: params.code,
      });
      const response = await authApi.authLoginPost(body);

      sessionStorage.setItem('displayName', response.displayName);
      sessionStorage.setItem('avatarUrl', response.avatarUrl);

      this.$router.push({name: 'home'});
    } catch (e) {
      console.error(e);

      this.$router.push({name: 'home'});
    }
  },
};
</script>
