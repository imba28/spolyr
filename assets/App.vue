<template>
  <div
    id="app"
    class="d-flex flex-column min-vh-100"
  >
    <app-navbar />

    <b-container
      fluid
      class="d-flex flex-grow-1"
    >
      <b-row class="d-flex flex-grow-1">
        <b-col
          md="3"
          lg="2"
          class="bg-dark-lighter pt-4"
        >
          <aside
            class="position-sticky"
            style="stop: 1em"
          >
            <ul class="nav flex-column">
              <li class="nav-item">
                <router-link
                  :to="{name:'home'}"
                  class="nav-link"
                >
                  <i class="fas fa-home" /> Home
                </router-link>
              </li>
              <li class="nav-item">
                <router-link
                  :to="{name:'dashboard'}"
                  class="nav-link"
                >
                  <i class="fas fa-music" /> Dashboard
                </router-link>
              </li>
              <li class="nav-item">
                Songs
              </li>
              <li class="nav-item">
                <router-link
                  :to="{name:'home'}"
                  class="nav-link"
                >
                  <i class="fas fa-question" /> Without lyrics
                </router-link>
              </li>
              <li class="nav-item">
                <router-link
                  :to="{name:'home'}"
                  class="nav-link"
                >
                  <i class="fas fa-exclamation-triangle" /> Manual action required
                </router-link>
              </li>
            </ul>
          </aside>
        </b-col>
        <b-col
          md="9"
          lg="10"
          class="p-0 py-4"
        >
          <main role="main">
            <router-view />
          </main>
        </b-col>
      </b-row>
    </b-container>

    <footer class="bg-dark text-lg-start">
      <div
        class="text-center p-3"
        style="background-color: rgba(0, 0, 0, 0.2)"
      >
        0.0.4
        <a
          class="text-light"
          href="/"
        >spolyr.app</a>
      </div>
    </footer>
  </div>
</template>

<script>
import AppNavbar from './components/AppNavbar';

export default {
  name: 'App',
  components: {
    AppNavbar,
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
