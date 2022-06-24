<template>
  <div :class="`bg-${variant} text-white status-card`">
    <b-row>
      <b-col cols="4 d-flex justify-content-center align-items-center">
        <slot name="icon" />
      </b-col>
      <b-col cols="8 d-flex justify-content-start align-items-center">
        <div>
          <div class="status-card__headline">
            <slot name="text" />
          </div>
          <div class="status-card__number">
            <slot name="number" />
          </div>
        </div>
      </b-col>
    </b-row>
  </div>
</template>

<script>
export default {
  props: {
    variant: {
      type: String,
      default: 'primary',
    },
  },
};
</script>

<style lang="scss">
.status-card {
  padding: 1rem;

  &__number {
    font-size: 2rem;
    font-weight: 700;
  }
}
</style>
