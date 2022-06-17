<template>
  <div
    class="card"
    style="height: 100%"
  >
    <img
      class="card-img-top"
      :src="cover"
      :alt="`album cover of ${title} from ${artists.join(', and ')}`"
    >

    <div class="card-body d-flex justify-content-between flex-column">
      <div>
        <h5 class="card-title">
          {{ title }}
        </h5>
        <p class="card-text">
          <router-link
            v-for="artist in artists"
            :key="title + artist"
            :to="{name:'search', params: {q: artist}}"
          >
            {{ artist }}
          </router-link>
        </p>
      </div>

      <div class="card-text mt-1">
        <router-link
          :to="{name: 'track-detail', params: {id: id}}"
          class="d-inline-block btn btn-primary"
        >
          Details
        </router-link> <a
          v-if="spotifyId"
          :href="`https://open.spotify.com/track/${spotifyId}`"
          class="d-inline-block btn btn-primary"
          target="_blank"
          rel="noopener"
        >
          to Spotify
        </a>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    id: {
      type: Number,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
    cover: {
      type: String,
      default: null,
    },
    spotifyId: {
      type: String,
      default: null,
    },
    artists: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
