import {defineStore /* acceptHMRUpdate*/} from 'pinia';
import {AuthApi, AuthLoginPostRequest} from '@/openapi';
const authApi = new AuthApi();

export const useAuthStore = defineStore({
  id: 'auth',
  persist: true,
  state: () => ({
    avatarUrl: null,
    displayName: null,
  }),
  getters: {
    isAuthenticated(state) {
      return state.avatarUrl !== null || state.displayName !== null;
    },
  },

  actions: {
    async logout() {
      await authApi.authLogoutGet();

      this.$patch({
        avatarUrl: null,
        displayName: null,
      });
    },

    async login(code) {
      const body = AuthLoginPostRequest.constructFromObject({
        code,
      });
      const response = await authApi.authLoginPost(body);

      this.$patch({
        avatarUrl: response.avatarUrl,
        displayName: response.displayName,
      });
    },
  },
});

/* if (import.meta.webpackHot) {
  import.meta.webpackHot.accept(acceptHMRUpdate(useAuthStore, import.meta.webpackHot));
}*/
