<template>
  <b-container>
    <b-card>
      <b-card-text>
        <loading-button
          :loading="library.importing"
          variant="primary"
          class="mr-1"
          :disabled="library.importing"
          @click="importLibrary"
        >
          <template #loading>
            Importing library...
          </template>
          <i class="fab fa-spotify" /> Import tracks from your library
        </loading-button>

        <loading-button
          :loading="lyrics.importing"
          variant="primary"
          class="mr-1"
          :disabled="lyrics.importing"
          @click="importLyrics"
        >
          <template #loading>
            Importing lyrics...
          </template>
          <i class="fab fa-spotify" /> Import lyrics
        </loading-button>
      </b-card-text>
    </b-card>

    <div v-if="lyrics.importing">
      <h3 class="mt-4">
        Lyrics
      </h3>
      <div>
        <b-progress
          class="mt-2 mb-2"
          :max="lyrics.total"
          show-value
        >
          <b-progress-bar
            :value="lyrics.done"
            variant="success"
          >
            {{ lyrics.done }} / {{ lyrics.total }}
          </b-progress-bar>
        </b-progress>
        <b-row
          class="mb-2 justify-content-start"
          no-gutters
        >
          <b-col
            v-if="lyrics.tracksSuccessful > 0"
            cols="3"
            md="4"
            sm="6"
          >
            <status-card>
              <template #text>
                Lyrics found
              </template>
              <template #number>
                {{ lyrics.tracksSuccessful }}
              </template>
              <template #icon>
                <i class="fa fa-3x fa-check fa-bounce" />
              </template>
            </status-card>
          </b-col>
          <b-col
            v-if="lyrics.tracksError > 0"
            cols="3"
            md="4"
            sm="6"
          >
            <status-card variant="danger">
              <template #text>
                Lyrics not found
              </template>
              <template #number>
                {{ lyrics.tracksError }}
              </template>
              <template #icon>
                <i class="fa fa-3x fa-exclamation-triangle" />
              </template>
            </status-card>
          </b-col>
        </b-row>
        <log-file-list :text="lyrics.log" />
      </div>
    </div>

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
import LoadingButton from '@/components/LoadingButton';
import LogFileList from '@/components/LogFileList';
import StatusCard from '@/components/StatusCard';
const api = new PlaylistsApi();
const importApi = new ImportApi();

export default {
  components: {
    StatusCard,
    LogFileList,
    LoadingButton,
    PlaylistCard,
  },
  data: () => ({
    library: {
      importing: false,
    },
    lyrics: {
      importing: false,
      interval: null,
      total: null,
      done: null,
      log: '',
      tracksError: 0,
      tracksSuccessful: 0,
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
    this.loadLyricsImportStatus();
  },
  beforeDestroy() {
    if (this.lyrics.interval) {
      clearInterval(this.lyrics.interval);
    }
  },
  methods: {
    async loadLyricsImportStatus() {
      try {
        const response = await importApi.importLyricsGet();

        if (response.running) {
          if (!this.lyrics.interval) {
            this.lyrics.interval = setInterval(this.loadLyricsImportStatus, 2000);
          }
        } else if (this.lyrics.interval) {
          if (this.lyrics.tracksSuccessful > 0) {
            this.$toast.success(`Successfully imported lyrics of ${this.lyrics.tracksSuccessful} tracks!`, {
              timeout: false,
            });
          } else {
            this.$toast.warning(`Import was finished without retrieving any lyrics. That's strange... 🤔`, {
              timeout: false,
            });
          }

          clearInterval(this.lyrics.interval);
        }

        this.lyrics.importing = response.running;
        this.lyrics.total = response.tracksTotal;
        this.lyrics.done = response.tracksCompleted;
        this.lyrics.log = response.log;
        this.lyrics.tracksError = response.tracksError;
        this.lyrics.tracksSuccessful = response.tracksSuccessful;
      } catch (e) {
        console.log(e);
      }
    },
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
    async importLyrics() {
      try {
        await importApi.importLyricsPost();
        this.lyrics.importing = true;
        this.lyrics.interval = setInterval(this.loadLyricsImportStatus, 2000);
      } catch (e) {
        // todo
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
        this.$toast.error(e.body);
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
        this.$toast.error(e.body);
      } finally {
        this.playlists.loading = false;
      }
    },
  },
};
</script>
