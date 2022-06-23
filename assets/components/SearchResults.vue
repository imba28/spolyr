<template>
  <div>
    <b-table
      hover
      class="search-results"
      head-variant="light"
      :busy="loading"
      :fields="fields"
      :items="items"
      @row-hovered="(row) => playing && row.previewURL !== playing ? $emit('play', row.previewURL) : null"
    >
      <template #cell(icons)="data">
        <div class="d-flex justify-content-end">
          <b-icon-music-note
            v-if="data.item.previewURL"
            v-b-popover.hover.top="'Lyrics available'"
            class="mr-1"
          />
          <b-icon-chat-quote-fill
            v-if="data.item.hasLyrics"
            v-b-popover.hover.top="'Audio snippet available'"
          />
        </div>
      </template>

      <template #cell(actions)="data">
        <div class="d-flex justify-content-end">
          <b-button
            v-if="data.item.previewURL"
            variant="primary"
            size="sm"
            class="player-button mr-1"
            @click="$emit(playing ? 'stop' : 'play', data.item.previewURL)"
          >
            <b-icon-pause-fill v-if="playing === data.item.previewURL" />
            <b-icon-play-fill v-else />
          </b-button>
          <router-link
            class="btn btn-primary"
            :to="{name:'track-detail', params: {id: data.item.spotifyId}}"
          >
            Detail
          </router-link>
        </div>
      </template>
    </b-table>
  </div>
</template>

<script>
import {BIconMusicNote, BIconChatQuoteFill, BIconPlayFill, BIconPauseFill} from 'bootstrap-vue';

export default {
  components: {
    BIconMusicNote,
    BIconChatQuoteFill,
    BIconPlayFill,
    BIconPauseFill,
  },
  props: {
    loading: {
      type: Boolean,
      default: false,
    },
    items: {
      type: Array,
      required: true,
    },
    playing: {
      type: String,
      default: null,
    },
  },
  data: () => ({
    fields: [
      {key: 'title', label: 'Title'},
      {key: 'artists', label: 'Artists', formatter: (v) => v.join(', ')},
      {key: 'album', label: 'Album'},
      {key: 'icons', label: ''},
      {key: 'actions', label: ''},
    ],
  }),
};
</script>


<style lang="scss">
.search-results {
  tbody {
    td {
      padding: .25rem;
      vertical-align: middle;
    }
  }
}
</style>
