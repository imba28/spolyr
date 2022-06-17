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
    loadResults() {
      this.loading = true;

      setTimeout(() => {
        this.loading = false;
        this.items = [
          {
            title: 'Tresor',
            spotifyId: '0k5UCvZ1BDP28db9y7OymD',
            artists: [
              'JURI', 'AK AUSSERKONTROLLE',
            ],
            album: 'Tresor',
            hasLyrics: true,
            // eslint-disable-next-line
            previewURL: 'https://p.scdn.co/mp3-preview/a602c8aedcb0eec26fcb64289a1cc15861e89350?cid=fa338f8a902f43dfb275be2eb27d96e5',
          },
          {
            title: 'Drug Test',
            spotifyId: '0k5UCvZ1BDP28db9y7OymD',
            artists: [
              'JURI', 'AK AUSSERKONTROLLE',
            ],
            album: 'The R.E.D Album',
            hasLyrics: true,
          },
          {
            title: 'Praise The Lord',
            spotifyId: '0k5UCvZ1BDP28db9y7OymD',
            artists: [
              'JURI',
            ],
            album: 'aaa',
            hasLyrics: true,
            // eslint-disable-next-line
            previewURL: 'https://p.scdn.co/mp3-preview/a515fa169153249641255877d5a36af261cf0d7c?cid=fa338f8a902f43dfb275be2eb27d96e5'
          },
        ];
      }, 500);
    },
  },
};
</script>
