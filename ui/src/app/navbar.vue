<style scoped>
  ul,
  li {
    margin: 0;
    padding: 0;
    list-style: none;
  }
  .navbar {
    top: 0;
    left: 0;
    width: calc(100% - 70px);
    position: fixed;
    margin: 0 20px;
    padding: 0 15px;
    min-height: 50px;
    border-radius: 0;
    border-bottom: 0.5px solid #b2c1d8;
    background-color: rgba(255, 255, 255, 0.8);
  }
  .bitmark {
    width: 120px;
    height: 46px;
    margin: 0 10px;
  }
  .bitmark img {
    margin-top: 9px;
    width: 120px;
  }
  .menu {
    float: left;
    height: 46px;
  }
  .menu .trademark {
    margin: 0;
    padding: 0 16px;
    line-height: 50px;
  }
  .menu.right {
    float: right;
  }
  .menu>a.active,
  .menu>a:hover {
    border-bottom: 3px solid #b2c1d8
  }
  .menu>a {
    display: inline-block;
    line-height: 0;
    height: 0;
    box-sizing: border-box;
    padding: 23px;
    color: black;
    border-bottom: 3px solid transparent;
    text-decoration: none;
  }
  .dropdown-menu {
    margin: -4px;
    padding: 10px 0;
  }
</style>

<template lang="pug">
div.navbar
  div.menu.bitmark
    img(src="/static/images/bitmark-logo.png")
  div.menu
    router-link(to="/node") Node
</template>

<script>
  import {
    setCookie
  } from "../utils"
  export default {
    methods: {
      openAdmin() {
        this.showAdmin = true
        setTimeout(() => window.addEventListener('click', this.closeAdmin), 0)
      },
      closeAdmin() {
        this.showAdmin = false
        window.removeEventListener('click', this.closeAdmin)
      },
      logout(e) {
        e.preventDefault();
        setCookie("bitmark-webgui", "", 0)
        setCookie("bitmark-webgui-network", "", 0)
        this.$router.push({
          path: "/login"
        })
      }
    },
    data() {
      return {
        showAdmin: false
      }
    }
  }
</script>
