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
        <search-form
          :value="searchStore.keywords"
          @input="search"
        />
      </b-navbar-nav>

      <b-navbar-nav class="ml-auto">
        <b-nav-item-dropdown
          v-if="authStore.isAuthenticated"
          right
        >
          <template #button-content>
            <span class="mr-1 text-white font-weight-bold">
              {{ authStore.displayName }}
            </span>
            <b-avatar
              :src="authStore.avatarUrl"
              variant="dark"
            />
          </template>

          <b-dropdown-item @click="logout">
            <i class="fa fa-sign-out" /> Sign out
          </b-dropdown-item>
        </b-nav-item-dropdown>
        <b-nav-item
          v-else
          @click="login"
        >
          <div class="btn btn-dark">
            <i class="fa fa-sign-in" /> Sign in
          </div>
        </b-nav-item>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script>
import SearchForm from './SearchForm.vue';
import {AuthApi} from '@/openapi';
import querystring from 'querystring';
import {useAuthStore, useSearchStore} from '@/stores';
import {mapStores} from 'pinia';

const authClient = new AuthApi();

export default {
  components: {
    SearchForm,
  },
  computed: {
    ...mapStores(useAuthStore, useSearchStore),
  },
  methods: {
    search(query) {
      this.searchStore.keywords = query;
      this.$router.push({name: 'search', params: {q: query}}).catch(() => {});
    },
    async logout() {
      try {
        await this.authStore.logout();
      } catch (e) {
        this.$toast.error('Something went wrong while trying to sign you out!');
      }
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
