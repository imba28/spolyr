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
    <div v-else-if="items !== null">
      <div class="text-muted mt-1 mb-3">
        {{ totalRows }} tracks found
      </div>

      <b-pagination
        v-if="totalRows > perPage"
        v-model="currentPage"
        class="justify-content-center"
        :total-rows="totalRows"
        :per-page="perPage"
        aria-controls="my-table"
        :ellipsis="true"
      />
      <search-results
        :loading="loading"
        :items="items"
        :playing="playing"
        @play="startPlaying"
        @stop="stopPlaying"
      />
      <b-pagination
        v-if="totalRows > perPage"
        v-model="currentPage"
        class="justify-content-center"
        :total-rows="totalRows"
        :per-page="perPage"
        aria-controls="my-table"
      />
    </div>
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
    currentPage: 1,
    perPage: 12,
    totalRows: null,
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
    currentPage() {
      this.loadResults();
    },
  },
  mounted() {
    this.loadResults();
  },
  async beforeDestroy() {
    await this.stopPlaying();
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
          page: this.currentPage,
          limit: this.perPage,
        });

        this.totalRows = results.meta.total ?? 0;
        this.items = results.data ?? [];
      } catch (e) {
        this.items = [];
        this.totalRows = 0;
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>
