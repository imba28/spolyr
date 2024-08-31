<template>
  <WordHighlighter
    :query="keywords"
    :caseSensitive="false"
    :splitBySpace="true"
    :diacriticsSensitive="true"
    highlightClass="mark"
  >
    <slot></slot>
  </WordHighlighter>
</template>

<script>
import WordHighlighter from 'vue-word-highlighter';

export default {
  components: {WordHighlighter},
  props: {
    keywords: {
      type: String,
    },
  },
};
</script>

<style>
.mark {
    background-color: var(--primary);
    padding: 0;
    box-shadow: 0 0 3px rgba(0,0,0,0.15);
    color: var(--white);
}
</style>
