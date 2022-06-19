<template>
  <b-navbar
    toggleable="lg"
    type="dark"
    variant="primary"
    class="py-2"
  >
    <b-navbar-brand href="#">
      <i class="fas fa-music" /> <strong>Spo</strong>lyr
    </b-navbar-brand>

    <b-navbar-toggle target="nav-collapse" />

    <b-collapse
      id="nav-collapse"
      is-nav
    >
      <b-navbar-nav class="ml-auto">
        <search-form @search="search" />
      </b-navbar-nav>

      <b-navbar-nav class="ml-auto">
        <b-nav-item-dropdown
          v-if="loggedIn"
          right
        >
          <template #button-content>
            <b-avatar
              src="https://i.scdn.co/image/ab6775700000ee85c83eb00882a56e93684d7ccc"
              variant="dark"
            />
          </template>

          <b-dropdown-item href="#">
            Sign Out
          </b-dropdown-item>
        </b-nav-item-dropdown>
        <b-nav-item
          v-else
          @click="login"
        >
          Link
        </b-nav-item>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script>
import SearchForm from './SearchForm.vue';
import {AuthApi} from '@/openapi';
import querystring from 'querystring';
const authClient = new AuthApi();

export default {
  components: {
    SearchForm,
  },
  data: () => ({
    loggedIn: false,
  }),
  methods: {
    search(query) {
      this.$router.push({name: 'search', params: {q: query}}).catch(() => {});
    },
    async login() {
      const config = await authClient.authConfigurationGet();
      window.location = 'https://accounts.spotify.com/authorize?' + querystring.stringify({
        response_type: 'code',
        client_id: config.clientId,
        scope: config.scope,
        redirect_uri: config.redirectUrl,
      });
    },
  },
};
</script>
