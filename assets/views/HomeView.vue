<template>
  <b-container>
    <div class="text-muted mt-1 mb-3">
      Indexed 4159 tracks including 3600 tracks with known lyrics.
    </div>
    <h1 class="h3">
      Latest songs
    </h1>

    <div
      v-if="latestTracks === null"
      class="text-center"
    >
      <b-spinner
        variant="primary"
      />
    </div>

    <b-row v-else>
      <b-col
        v-for="track in latestTracks"
        :key="track.id"
        sm="6"
        md="4"
        lg="3"
        class="mb-4"
      >
        <track-card
          :id="track.id"
          :title="track.title"
          :cover="track.cover"
          :spotify-id="track.spotifyId"
          :artists="track.artists"
        />
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import TrackCard from '../components/TrackCard.vue';
export default {
  components: {
    TrackCard,
  },
  data() {
    return {
      latestTracks: null,
    };
  },
  mounted() {
    setTimeout(() => {
      this.latestTracks = [
        {
          id: 1,
          title: 'Tresor',
          spotifyId: '0k5UCvZ1BDP28db9y7OymD',
          artists: [
            'JURI', 'AK AUSSERKONTROLLE',
          ],
          album: 'Tresor',
          cover: 'https://i.scdn.co/image/ab67616d0000b2733310f6638fa218e34d5be19e',
        },
      ];
    }, 1500);
  },
};
</script>

