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
                <span v-if="track.artists.length > 0">
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
                Lyrics

                <div v-if="track.hasLyrics">
                  <b-button
                    v-if="editMode"
                    variant="primary"
                    @click="save"
                  >
                    <i class="fa fa-save" /> Save
                  </b-button>
                  <b-button
                    v-else
                    variant="primary"
                    @click="edit"
                  >
                    <i class="fa fa-edit" /> Edit
                  </b-button>
                </div>
              </div>
            </template>
            <b-card-body>
              <b-card-text>
                <div
                  v-if="track.lyrics"
                >
                  <textarea
                    v-if="editMode"
                    v-model="track.lyrics"
                    class="form-control"
                    :style="textareaStyle"
                  />
                  <span
                    v-else
                    class="lyrics-text"
                  >{{ track.lyrics }}</span>
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
import {Lyrics, TracksApi} from '@/openapi';
import ErrorBox from '@/components/ErrorBox';
const api = new TracksApi();

export default {
  components: {ErrorBox},
  data: () => ({
    track: null,
    notFound: false,
    maxImportErrorCount: 3,
    editMode: false,
  }),
  computed: {
    textareaStyle() {
      return {
        height: 1.5 * this.track.lyrics.match(/\n/g).length + 'em',
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
        this.$toast.error(e.message);
      } finally {
        this.editMode = false;
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
