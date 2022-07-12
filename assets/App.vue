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
          <app-sidebar />
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
        dev-build
        <a
          class="text-light"
          rel="noopener"
          target="_blank"
          href="https://github.com/imba28/spolyr"
        >github.com/imba28/spolyr</a>
      </div>
    </footer>
  </div>
</template>

<script>
import AppNavbar from './components/AppNavbar';
import AppSidebar from '@/components/AppSidebar';

export default {
  name: 'App',
  components: {
    AppSidebar,
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
