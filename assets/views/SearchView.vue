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

    <error-box
      v-if="items !== null && items.length === 0"
      message="Sorry, no results found"
      text="Please try another query."
    />
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
import Player from '@/audio-player';
import {TracksApi} from '@/openapi';
import ErrorBox from '@/components/ErrorBox';

const tracksApi = new TracksApi();

const player = new Player();

export default {
  components: {ErrorBox, SearchResults},
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
