<template>
  <b-container>
    <div
      v-if="trackStats"
      class="text-muted mt-1 mb-3"
    >
      Indexed {{ trackStats.numberOfTracks }} tracks
      including {{ trackStats.numberOfTracksWithLyrics }} songs with known lyrics.
    </div>


    <div
      v-if="latestTracks === null"
      class="text-center"
    >
      <b-spinner
        variant="primary"
      />
    </div>

    <div v-else-if="latestTracks.length > 0">
      <h1 class="h3">
        Latest songs
      </h1>

      <b-row>
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
    </div>
    <div
      v-else
      class="mt-1"
    >
      <h3>Oh no, it seems your library is currently empty.</h3>

      <div v-if="authStore.isAuthenticated">
        You might want to
        <router-link :to="{name: 'dashboard'}">
          import
        </router-link> your Spotify library first.
      </div>
      <div v-else>
        Maybe you want to sign in and change that?
      </div>
    </div>
  </b-container>
</template>

<script>
import TrackCard from '../components/TrackCard.vue';
import {TracksApi} from '@/openapi';
import {useAuthStore} from '@/stores/auth';
import {mapStores} from 'pinia';

const tracksApi = new TracksApi();

export default {
  components: {
    TrackCard,
  },
  data() {
    return {
      latestTracks: null,
      trackStats: null,
    };
  },
  computed: {
    ...mapStores(useAuthStore),
  },
  async mounted() {
    try {
      this.trackStats = await tracksApi.tracksStatsGet();
    } catch (e) {
      console.error(e);
    }

    try {
      const response = await tracksApi.tracksGet({
        limit: 8,
      });
      this.latestTracks = response.data;
    } catch (e) {
      console.error(e);
    }
  },
};
</script>

