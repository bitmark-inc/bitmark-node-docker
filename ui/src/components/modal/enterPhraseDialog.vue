<style lang="scss" scoped>
  .items {
    display: flex;
    flex-wrap: wrap;
    flex-direction: column;
    justify-content: center;
    height: 475px;
  }

  .items .item {
    flex: 1;
    box-sizing: border-box;
    margin: 8px 20px 8px 0;
    color: #171e42;
  }

  .item input {
    border: none;
    outline: none;
    color: blue;
    border-bottom: solid 1px black;
  }

  p.error {
    color: red;
    font-weight: bold;
  }

  .error .item input {
    color: red;
  }
</style>

<template>
  <!-- POP UP - welcome dialog -->
  <div class="pop-up-box fullscreen">
    <div class="pop-up-box__content">
      <div class="pop-up-box__header">
        <h1>Enter recovery phrase</h1>
      </div>
      <p>
        Enter your 24-word recovery phrase to access your existing account:
      </p>
      <form class="items" :class="{error: errMsg}">
        <template v-for="(word, index) in words">
          <div class="item">
            {{ 1 + parseInt(index) + '.' }}
            <input type="text" v-model='words[index]'>
          </div>
        </template>
      </form>
      </p>
      <p v-if="errMsg" class="error">Can not recover the account. Error: {{ errMsg }}</p>
      <div class="button">
        <button type="button" class="btn-default-fill" @click="backToWelcome">BACK</button>
        <button type="button" class="btn-default-fill" @click="recoverAccount">NEXT</button>
      </div>
    </div>
  </div>
</template>

<script>
  import axios from "axios"

  export default {
    methods: {
      recoverAccount() {
        let phrase = this.words.join(" ")
        axios.post("/api/account/phrase", {
            "phrase": phrase
          })
          .then(resp => {
            this.$emit('accountRecovered', )
          })
          .catch(err => {
            this.errMsg = err.response.data.message
          })
      },

      backToWelcome() {
        this.$emit('backToWelcome')
      }
    },

    data() {
      let words = [];
      for (let i = 0; i < 24; i++) {
        words[i] = ""
      }
      return {
        words: words,
        errMsg: "",
      }
    }
  }
</script>
