<template>
  <div class="box">
    <div class="box__header">
      <span class="title">{{ title }}</span>
      <span class="btn">
        <slot name="header-button"></slot>
      </span>
    </div>
    <div class="box__body">
      <slot></slot>
    </div>
  </div>
</template>

<script>
  export default {
    props: {
      title: String
    }
  }
</script>
