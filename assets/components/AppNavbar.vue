<template>
  <b-navbar
    toggleable="lg"
    type="dark"
    variant="primary"
    class="py-2"
  >
    <b-navbar-brand href="#">
      <i class="fas fa-music" /> <strong>Spo</strong>lyr
    </b-navbar-brand>

    <b-navbar-toggle target="nav-collapse" />

    <b-collapse
      id="nav-collapse"
      is-nav
    >
      <b-navbar-nav>
        <b-nav-item href="#">
          Tracks
        </b-nav-item>
      </b-navbar-nav>

      <b-navbar-nav>
        <search-form />
      </b-navbar-nav>

      <b-navbar-nav class="ml-auto">
        <b-nav-item-dropdown right>
          <template #button-content>
            <b-avatar
              src="https://i.scdn.co/image/ab6775700000ee85c83eb00882a56e93684d7ccc"
              variant="dark"
            />
          </template>

          <b-dropdown-item href="#">
            Sign Out
          </b-dropdown-item>
        </b-nav-item-dropdown>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script>
import SearchForm from './SearchForm.vue';

export default {
  components: {
    SearchForm,
  },
};
</script>
