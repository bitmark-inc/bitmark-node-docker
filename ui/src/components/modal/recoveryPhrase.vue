<style scoped>
.phrases {
  list-style-type: none;
  columns: 2;
  -webkit-columns: 2;
  -moz-columns: 2;
}
</style>

<template>
  <!-- POP UP - payment address -->
  <div class="pop-up-box">
    <div class="pop-up-box__content">
      <div class="pop-up-box__header">
        <h1>Write Down Recovery Phrase</h1>
        <span class="close" @click="close">
          <svg class="icon-hamburger">
            <use xlink:href="assets/img/icons.svg#icon-cancel" xmlns:xlink="http://www.w3.org/1999/xlink"></use>
          </svg>
        </span>
      </div>
      <p>
        Please write down your recovery phrase in the exact sequence below:
      </p>
      <p>
        <ol class="phrases">
          <li v-for="(phrase, index) in phrases">{{ index + 1 }}. {{ phrase }}</li>
        </ol>
      </p>
      <div class="button">
        <button type="button" class="btn-default-fill" @click="close">DONE</button>
      </div>
    </div>
  </div>
</template>

<script>
  import axios from "axios"

  export default {

    created() {
      axios.get("/api/" + "account/phrase")
        .then((resp) => {
          let data = resp.data
          if (data.ok) {
            this.phrases = data.result.split(" ")
          }
        })
    },

    methods: {
      close() {
        this.$emit('close')
      }
    },

    data() {
      return {
        phrases: []
      }
    }
  }
</script>
