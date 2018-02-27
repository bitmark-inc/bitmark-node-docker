<template>
  <!-- POP UP - payment address -->
  <div class="pop-up-box" >
    <div class="pop-up-box__content">
      <div class="pop-up-box__header">
        <h1>WARNING!</h1>
        <span class="close" @click="close">
          <svg class="icon-hamburger">
            <use xlink:href="assets/img/icons.svg#icon-cancel" xmlns:xlink="http://www.w3.org/1999/xlink"></use>
          </svg>
        </span>
      </div>
      <p>
        Your recovery phrase is the only way to restore your Bitmark account. <br>
        We will show you a list of words to write down on a piece of paper and keep safe. <br>
        Make sure you are in a private location before writing down your recovery phrase.
      </p>
      <div class="button">
        <button type="button" class="btn-default-fill" @click="showPhrase">WRITE DOWN RECOVERY PHRASE</button>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    methods: {
      showPhrase() {
        this.$emit('showPhrase')
      },

      close() {
        this.$emit('close')
      }
    }
  }
</script>
