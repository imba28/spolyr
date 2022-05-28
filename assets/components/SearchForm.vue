<template>
  <b-nav-form
    style="min-width: 25vw"
    @submit.prevent="$emit('search', query)"
  >
    <b-input-group>
      <b-form-input
        v-model="query"
        placeholder="Search for lyrics, titles, artists, album names ..."
        autocomplete="off"
      />

      <template #append>
        <b-button variant="secondary">
          <i class="fas fa-search" />
        </b-button>
      </template>
    </b-input-group>
  </b-nav-form>
</template>

<script>
export default {
  data() {
    return {
      query: '',
    };
  },
};
</script>
