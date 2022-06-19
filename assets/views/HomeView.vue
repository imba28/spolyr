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
          :title="track.title"
          :cover="track.coverImage"
          :spotify-id="track.spotifyId"
          :artists="track.artists"
        />
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import TrackCard from '../components/TrackCard.vue';
import {TracksApi} from '@/openapi';

const tracksApi = new TracksApi();

export default {
  components: {
    TrackCard,
  },
  data() {
    return {
      latestTracks: null,
    };
  },
  async mounted() {
    try {
      const response = await tracksApi.tracksGet({
        limit: 10,
      });
      this.latestTracks = response.data;
    } catch (e) {
      console.error(e);
    }
  },
};
</script>

