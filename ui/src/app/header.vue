<template>
  <div class="header">
    <div class="header__top">
      <div class="logo">
        <img src="assets/img/bitmark-logo.svg" alt="Bitmark logo">
      </div>
      <div class="account">
        <div class="user">
          <span>Account: {{ abbrAccount }}</span>
          <i>
            <svg class="icon-account-circle">
              <use xlink:href="assets/img/icons.svg#icon-account-circle" xmlns:xlink="http://www.w3.org/1999/xlink"></use>
            </svg>
          </i>
          <ul class="menu account__menu">
            <li>
              <a :href="blockLink" target="_blank">View blocks won on Registry</a>
            </li>
            <li>
              <a href="#" @click="openRecoveryAlert">Write down recovery phrase</a>
            </li>
            <li>
              <a href="#" @click="copyAccount">Copy account address</a>
            </li>
          </ul>
        </div>
        <div class="hambuger">
          <svg class="icon-hamburger">
            <use xlink:href="assets/img/icons.svg#icon-hamburger" xmlns:xlink="http://www.w3.org/1999/xlink"></use>
          </svg>
          <ul class="menu hambuger__menu">
            <li>
              <svg class="icon-description">
                <use xlink:href="assets/img/icons.svg#icon-description" xmlns:xlink="http://www.w3.org/1999/xlink"></use>
              </svg>
              <a href="https://hub.docker.com/r/bitmark/bitmark-node/" target="_blank">Instructions</a>
            </li>
            <li>
              <svg class="icon-language">
                <use xlink:href="assets/img/icons.svg#icon-language" xmlns:xlink="http://www.w3.org/1999/xlink"></use>
              </svg>
              <span class="lan-title">Languages Settings</span>
              <span class="lan-opt">
                <a href="">EN</a> |
                <a class="active" href="">中文</a>
              </span>
            </li>
          </ul>
        </div>
      </div>
    </div>
    <div class="divider"></div>
    <div class="header__bottom">
      <div class="blockchain">
        <div class="chain-title">Blockchain:</div>
        <div class="chain-select">
          <ul class="chain-opt">
            <li class="active">
              <a href="#bitmark">{{ nodeInfo.network }}</a>
              <p>The official version of the Bitmark blockchain</p>
            </li>
            <li>
              <a href="#testing">Testing</a>
              <p>A testnet version of the blockchain used solely for development testing</p>
            </li>
          </ul>
        </div>
      </div>
      <div class="bitmark-version">
        <div class="bitmark-node-version">Bitmark node {{ nodeInfo.version }}</div>
        <div class="bitmark-node-update">
          <a href="https://hub.docker.com/r/bitmark/bitmark-node/tags/" target="_blank" v-if="(latestVersion > nodeInfo.version)">New update available</a>
        </div>
      </div>
    </div>
    <div class="divider"></div>
  </div>
</template>


<script>

  import axios from 'axios';

  export default {
    props: {
      nodeInfo: Object
    },

    computed: {
      blockLink() {
        let l = (this.nodeInfo.network === 'bitmark') ? "https://registry.bitmark.com" : "https://registry.test.bitmark.com"
        return l + '/' + this.nodeInfo.account
      },

      abbrAccount() {
        let a = this.nodeInfo.account || ""
        return a.slice(0, 6) + "......" + a.slice(-6)
      }
    },

    mounted() {
      setTimeout(this.getLatestVersion, 6000)
    },

    methods: {
      getLatestVersion() {
        axios.get('/api/latestVersion')
          .then((resp) => {
            this.latestVersion = resp.data.latest
          })
          .catch((error) => {
            this.$emit("error", error)
          })
      },
      copyAccount() {
        let i = document.createElement("input")
        document.body.appendChild(i)
        i.value = this.nodeInfo.account
        i.select()
        document.execCommand("copy");
        i.remove();
      },

      openRecoveryAlert () {
        this.$emit("openRecoveryAlert")
      }
    },

    data() {
      return {
        latestVersion: this.nodeInfo.version,
      }
    }
  }
</script>
