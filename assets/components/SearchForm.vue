<template>
  <b-nav-form
    class="search-form"
    @submit.prevent="$emit('search', query)"
  >
    <b-input-group>
      <b-form-input
        v-model="query"
        placeholder="Search for lyrics, titles, artists, album names ..."
        autocomplete="off"
      />

      <template #append>
        <b-button
          variant="secondary"
          type="submit"
          aria-label="Search"
        >
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

<style lang="scss">
.search-form {
  input[type='text'] {
    min-width: 25vw;
  }
}

@media (max-width: 992px) {
  .search-form {
    margin: 1em 0;

    .form-inline {
      display: block;
      width: 100%;
    }
  }
}
</style>
