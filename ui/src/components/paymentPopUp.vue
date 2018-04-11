<template>
  <!-- POP UP - payment address -->
  <div class="pop-up-box fullscreen" >
    <div class="pop-up-box__content">
      <div class="pop-up-box__header">
        <h1>Please enter your payment addresses.</h1>
        <span class="close" @click="close">
          <svg class="icon-hamburger">
            <use xlink:href="assets/img/icons.svg#icon-cancel" xmlns:xlink="http://www.w3.org/1999/xlink"></use>
          </svg>
        </span>
      </div>
      <p>Bitcoin and litecoin payments for blocks won will be transferred to the following addresses:</p>
      <div class="wallet-address">
        <div class="bitcoin">
          <i>
            <img src="assets/img/icons/ic_bitcoin.svg " alt="icon bitcoin ">
          </i>
          <span class="coin-title">BTC Address:</span>
          <span class="field ">
            <input type="text" v-model="btcAddr">
          </span>
        </div>
        <div class="litecoin">
          <i>
            <img src="assets/img/icons/ic_litecoin.svg " alt="icon litecoin ">
          </i>
          <span class="coin-title ">LTC Address:</span>
          <span class="field ">
            <input type="text" v-model="ltcAddr">
          </span>
        </div>
      </div>
      <div class="button">
        <button type="button" class="btn-default-fill" :disabled="!btcAddr || !ltcAddr" @click="save">Save</button>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    props: {
      initBtcAddr: String,
      initLtcAddr: String,
      show: true
    },

    methods: {
      save() {
        this.$emit('saved', {
          btcAddr: this.btcAddr,
          ltcAddr: this.ltcAddr,
        })
      },

      close() {
        this.$emit('close')
      }
    },

    data() {
      return {
        btcAddr: this.initBtcAddr,
        ltcAddr: this.initLtcAddr
      }
    }
  }
</script>
