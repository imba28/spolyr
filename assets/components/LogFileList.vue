<template>
  <div
    v-if="text"
    class="lig-file-list"
  >
    <div v-html="text" />
  </div>
</template>

<script>
export default {
  props: {
    text: {
      type: String,
      default: '',
    },
  },
};
</script>

<style lang="scss">
.lig-file-list {
  position: relative;
  max-height: 200px;
  overflow-y: hidden;

  &::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    width: 100%;
    background: linear-gradient(0deg, rgba(42,58,72,1) 3%, rgba(42,58,72,0) 40%);
  }
}
</style>
