<template>
  <div
    class="card"
    style="height: 100%"
  >
    <img
      class="card-img-top"
      :src="cover"
      :alt="`cover of playlist ${name} created by ${owner}`"
    >

    <div class="card-body d-flex justify-content-between flex-column">
      <div>
        <h5 class="card-title">
          {{ name }}
        </h5>
        <p class="card-text">
          created by <span class="text-muted">{{ owner }}</span>
        </p>
      </div>

      <div class="card-text mt-1">
        <b-button
          variant="primary"
          :disabled="isImporting"
          :block="true"
          @click="$emit('button-click')"
        >
          <div v-if="isImporting">
            <b-spinner
              small
              label="Busy"
            /> Importing...
          </div>
          <span
            v-else
          >
            <i class="fa fa-music" /> Import {{ trackCount }} tracks
          </span>
        </b-button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    spotifyId: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    cover: {
      type: String,
      default: null,
    },
    owner: {
      type: String,
      default: null,
    },
    trackCount: {
      type: Number,
      default: null,
    },
    isPublic: {
      type: Boolean,
      default: false,
    },
    isCollaborative: {
      type: Boolean,
      default: false,
    },
    isImporting: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
