<template>
  <div class="text-center">
    <img
      src="@/static/undraw_happy_music_g6wc.svg"
      alt="error image"
      style="width: 35%;"
    >
    <h1 class="h3 mt-4">
      {{ message }}
    </h1>
    <div
      v-if="text"
      class="text-muted"
    >
      {{ text }}
    </div>
  </div>
</template>

<script>
export default {
  props: {
    message: {
      type: String,
      required: true,
    },
    text: {
      type: String,
      default: null,
    },
  },
};
</script>
