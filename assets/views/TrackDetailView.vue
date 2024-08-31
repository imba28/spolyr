<template>
  <b-container>
    <error-box
      v-if="notFound"
      message="Sorry, the requested track was not found."
      text="Please try another one"
    />
    <div v-else-if="track !== null">
      <div class="mt-2 mb-2">
        <b-button
          variant="primary"
          @click="$router.back()"
        >
          <i class="fas fa-arrow-left" /> Back
        </b-button>
      </div>

      <b-row>
        <b-col
          md="4"
          sm="12"
        >
          <b-card
            class="position-sticky"
            style="top: 1em"
          >
            <img
              v-if="track.coverImage"
              :src="track.coverImage"
              class="img-fluid"
            >

            <b-card-body>
              <h5 class="card-title">
                {{ track.title }}
              </h5>
              <p class="card-text">
                <span v-if="track.artists">
                  <router-link
                    v-for="artist in track.artists"
                    :key="`artist-${artist}`"
                    :to="{name:'search', params: {q:artist}}"
                  >
                    {{ artist }}
                  </router-link> -
                </span>
                {{ track.album }}
              </p>
              <div class="card-text">
                <audio
                  v-if="track.previewURL"
                  data-testid="audio-player"
                  controls
                  preload="none"
                  style="width: 100%"
                >
                  <source
                    :src="track.previewURL"
                    type="audio/mpeg"
                  >
                </audio>
                <a
                  :href="`https://open.spotify.com/track/${ track.spotifyId }`"
                  class="d-inline-block mt-1 btn btn-primary"
                  target="_blank"
                  rel="noopener"
                >
                  <i class="fab fa-spotify" /> to Spotify
                </a>
              </div>
            </b-card-body>
          </b-card>
        </b-col>
        <b-col
          md="8"
          sm="12"
        >
          <b-card header-bg-variant="secondary">
            <template #header>
              <div class="d-flex align-items-center justify-content-between">
                <span>
                  Lyrics <span
                    v-if="track.language"
                    class="text-muted"
                  >({{ track.language }})</span>
                </span>

                <div v-if="authStore.isAuthenticated">
                  <b-button
                    v-if="editMode"
                    variant="primary"
                    aria-label="save lyrics"
                    :disabled="importingLyrics"
                    @click="save"
                  >
                    <i class="fa fa-save" /> Save
                  </b-button>
                  <b-button
                    v-else
                    variant="primary"
                    aria-label="edit lyrics"
                    :disabled="importingLyrics"
                    @click="edit"
                  >
                    <i class="fa fa-edit" /> Edit
                  </b-button>
                  <loading-button
                    v-if="!track.hasLyrics"
                    :loading="importingLyrics"
                    aria-label="import lyrics"
                    :disabled="importingLyrics"
                    @click="importLyrics"
                  >
                    <i class="fa fa-quote-left" /> Download lyrics
                  </loading-button>
                </div>
              </div>
            </template>
            <b-card-body>
              <b-card-text>
                <textarea
                  v-if="editMode"
                  v-model="track.lyrics"
                  class="form-control"
                  :style="textareaStyle"
                />
                <div v-else-if="track.lyrics">
                  <HighlightWords
                    :keywords="searchStore.keywords"
                    class="lyrics-text"
                  >
                    {{ track.lyrics }}
                  </HighlightWords>
                </div>
                <div v-else-if="track.lyricsImportErrorCount > maxImportErrorCount">
                  <small class="text-warning">Lyrics not found. Import
                    failed {{ track.lyricsImportErrorCount }} times.</small>
                </div>
                <div v-else-if="track.lyricsImportErrorCount > 0">
                  <small class="text-warning">Lyrics not found. Import
                    failed {{ track.lyricsImportErrorCount }} times.</small>
                </div>
                <div v-else>
                  <small class="text-muted">Lyrics not imported yet.</small>
                </div>
              </b-card-text>
            </b-card-body>
          </b-card>
        </b-col>
      </b-row>
    </div>
  </b-container>
</template>

<script>
import {ImportApi, Lyrics, TracksApi} from '@/openapi';
import ErrorBox from '@/components/ErrorBox';
import {mapStores} from 'pinia';
import {useAuthStore, useSearchStore} from '@/stores';
import LoadingButton from '@/components/LoadingButton';
import HighlightWords from '@/components/HighlightWords';

const api = new TracksApi();
const importApi = new ImportApi();

export default {
  components: {LoadingButton, ErrorBox, HighlightWords},
  data: () => ({
    track: null,
    notFound: false,
    maxImportErrorCount: 3,
    editMode: false,
    importingLyrics: false,
  }),
  computed: {
    ...mapStores(useAuthStore, useSearchStore),
    textareaStyle() {
      const lineBreaks = this.track.lyrics.match(/\n/g);
      if (!lineBreaks) {
        return {
          height: '5em',
        };
      }

      return {
        height: 1.5 * lineBreaks.length + 'em',
      };
    },
  },
  async mounted() {
    try {
      this.track = await api.tracksIdGet(this.$route.params.id);
    } catch (e) {
      this.notFound = true;
    }
  },
  methods: {
    edit() {
      this.editMode = true;
    },
    async save() {
      try {
        this.track = await api.tracksIdPatch(this.track.spotifyId, Lyrics.constructFromObject({
          lyrics: this.track.lyrics,
        }));
        this.$toast.success('You\'ve successfully updated the lyrics of this track!');
      } catch (e) {
        this.$toast.error(e.message ?? e.body);
      } finally {
        this.editMode = false;
      }
    },
    async importLyrics() {
      this.importingLyrics = true;

      try {
        this.track = await importApi.importLyricsTrackIdPost(this.track.spotifyId);
        this.$toast.success('Lyrics were found and successfully saved!');
      } catch (e) {
        this.$toast.info('Sadly, no lyrics were found.');
      } finally {
        this.importingLyrics = false;
      }
    },
  },
};
</script>

<style scoped>
.lyrics-text {
  white-space: pre-wrap;
}
</style>
