import {defineStore} from 'pinia';

export const useSearchStore = defineStore({
  id: 'search',
  persist: true,
  state: () => ({
    keywords: '',
  }),
});
