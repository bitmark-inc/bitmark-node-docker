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
              <a href="https://registry.bitmark.com">View blocks won on Registry</a>
            </li>
            <li>
              <a href="#">Write down recovery phrases</a>
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
              <a href="https://hub.docker.com/r/bitmark/bitmark-node/">Instructions</a>
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
        <div class="bitmark-node-version">Bitmark node v{{ nodeInfo.version }}</div>
        <div class="bitmark-node-update">
          <a href="#">New update available</a>
        </div>
      </div>
    </div>
    <div class="divider"></div>
  </div>
</template>


<script>
  export default {
    props: {
      nodeInfo: Object
    },

    computed: {
      abbrAccount() {
        let a = this.nodeInfo.account || ""
        return a.slice(0, 6) + "......" + a.slice(-6)
      }
    },

    methods: {
      copyAccount() {
        let i = document.createElement("input")
        document.body.appendChild(i)
        i.value = this.nodeInfo.account
        i.select()
        document.execCommand("copy");
        i.remove();
      }
    },

    data() {
      return {}
    }
  }
</script>
