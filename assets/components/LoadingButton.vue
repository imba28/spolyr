<template>
  <b-button
    :variant="variant"
    :disabled="disabled"
    @click="$emit('click')"
  >
    <div v-if="loading">
      <b-spinner
        small
        label="Busy"
      />
      <slot name="loading">
        Loading...
      </slot>
    </div>
    <div
      v-else
      class="ready"
    >
      <slot />
    </div>
  </b-button>
</template>

<script>
export default {
  props: {
    variant: {
      type: String,
      default: 'primary',
    },
    loading: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
