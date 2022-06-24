<template>
  <b-container>
    <b-card>
      <b-card-text>
        <b-button
          variant="primary"
          class="mr-1"
          :disabled="library.importing"
          @click="importLibrary"
        >
          <div v-if="library.importing">
            <b-spinner
              small
              label="Busy"
            /> Importing library...
          </div>
          <span v-else>
            <i class="fab fa-spotify" /> Import tracks from your library
          </span>
        </b-button>
        <b-button variant="primary">
          <i class="fa fa-quote-right" /> Import lyrics
        </b-button>
      </b-card-text>
    </b-card>

    <h3 class="mt-4">
      Your Playlists <span v-if="playlists.total">({{ playlists.total }})</span>
    </h3>
    <div
      v-if="playlists.loading"
      class="text-center mb-4"
    >
      <b-spinner
        variant="primary"
      />
    </div>
    <b-row v-else>
      <b-col
        v-for="playlist in playlists.items"
        :key="`playlist-${playlist.spotifyId}`"
        md="3"
        sm="6"
        class="mb-4"
      >
        <playlist-card
          :name="playlist.name"
          :spotify-id="playlist.spotifyId"
          :cover="playlist.coverImage"
          :is-collaborative="playlist.isCollaborative"
          :is-public="playlist.isPublic"
          :owner="playlist.owner"
          :track-count="playlist.trackCount"
          :is-importing="playlists.importingIds.includes(playlist.spotifyId)"
          @button-click="importPlaylist(playlist)"
        />
      </b-col>
    </b-row>
    <b-pagination
      v-if="playlists.total > playlists.perPage"
      v-model="playlists.currentPage"
      class="justify-content-center"
      :hide-ellipsis="true"
      :total-rows="playlists.total"
      :per-page="playlists.perPage"
    />
  </b-container>
</template>

<script>
import {PlaylistsApi, ImportApi} from '@/openapi';
import PlaylistCard from '@/components/PlaylistCard';
const api = new PlaylistsApi();
const importApi = new ImportApi();

export default {
  components: {
    PlaylistCard,
  },
  data: () => ({
    library: {
      importing: false,
    },
    playlists: {
      importingIds: [],
      items: [],
      loading: true,
      currentPage: 1,
      perPage: 8,
      total: null,
    },
  }),
  computed: {
    playlistCurrentPage() {
      return this.playlists.currentPage;
    },
  },
  watch: {
    playlistCurrentPage() {
      this.loadPlaylists();
    },
  },
  mounted() {
    this.loadPlaylists();
  },
  methods: {
    async importLibrary() {
      this.library.importing = true;
      try {
        await importApi.importLibraryPost();
        this.$toast.success('Successfully imported your library 🎶', {
          timeout: 2000,
        });
      } catch (e) {
        console.log(e);
        this.$toast.error('An error occurred while importing your library!');
      } finally {
        this.library.importing = false;
      }
    },
    async importPlaylist({spotifyId, name}) {
      if (this.playlists.importingIds.includes(spotifyId)) {
        return;
      }

      try {
        this.playlists.importingIds.push(spotifyId);
        await importApi.importPlaylistIdPost(spotifyId);
        this.$toast.success('Successfully imported your playlist ' + name, {
          timeout: 2000,
        });
      } catch (e) {
        console.log(e);
        this.$toast.error(e.message);
      } finally {
        this.playlists.importingIds.splice(this.playlists.importingIds.indexOf(spotifyId), 1);
      }
    },
    async loadPlaylists() {
      this.playlists.loading = true;

      try {
        const response = await api.playlistsGet({
          page: this.playlists.currentPage,
          limit: this.playlists.perPage,
        });
        this.playlists.items = response.data;
        this.playlists.total = response.meta.total;
      } catch (e) {
        console.error(e);
        this.$toast.error(e.message);
      } finally {
        this.playlists.loading = false;
      }
    },
  },
};
</script>
