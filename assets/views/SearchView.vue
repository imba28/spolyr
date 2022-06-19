<template>
  <b-container>
    <div
      v-if="loading && items === null"
      class="text-center"
    >
      <b-spinner
        variant="primary"
      />
    </div>
    <div
      v-if="items !== null && items.length === 0"
      class="text-center"
    >
      <img
        src="@/static/undraw_happy_music_g6wc.svg"
        style="width: 35%;"
      >
      <h1 class="h3 mt-4">
        Sorry, no results found
      </h1>
      <div class="text-muted">
        Please try another query.
      </div>
    </div>
    <search-results
      v-else-if="items !== null"
      :loading="loading"
      :items="items"
      :playing="playing"
      @play="startPlaying"
      @stop="stopPlaying"
    />
  </b-container>
</template>

<script>
import SearchResults from '@/components/SearchResults';
import Player from '@/track-page/player';
import {TracksApi} from '@/openapi';

const tracksApi = new TracksApi();

const player = new Player();

export default {
  components: {SearchResults},
  data: () => ({
    loading: true,
    playing: '',
    items: null,
  }),
  computed: {
    query() {
      return this.$route.params.q;
    },
  },
  watch: {
    query() {
      this.loadResults();
    },
  },
  mounted() {
    this.loadResults();
  },
  methods: {
    async startPlaying(url) {
      if (!url) {
        return;
      }
      await player.enable();
      await player.play(url);
      this.playing = url;
    },
    async stopPlaying() {
      await player.disable();
      this.playing = '';
    },
    async loadResults() {
      this.loading = true;

      try {
        const results = await tracksApi.tracksGet({
          query: this.query,
          page: 1,
          limit: 25,
        });
        this.items = results.data;
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>
