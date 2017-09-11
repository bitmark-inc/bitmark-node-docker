<style>
  @font-face {
    font-family: 'Avenir Next LT Pro';
    font-style: normal;
    src: url(/static/fonts/AvenirNextLTPro-Regular.otf) format('opentype'), local('Avenir Next LT Pro');
    font-weight: normal;
    font-style: normal;
  }

  @font-face {
    font-family: 'Avenir Next LT Pro';
    font-style: normal;
    src: url(/static/fonts/AvenirNextLTPro-Bold.otf) format('opentype'), local('Avenir Next LT Pro');
    font-weight: bold;
    font-style: normal;
  }

  @font-face {
    font-family: 'Avenir Next LT Pro';
    font-style: normal;
    src: url(/static/fonts/AvenirNextLTPro-Bold.otf) format('opentype'), local('Avenir Next LT Pro');
    font-weight: bold;
    font-style: bold;
  }

  body {
    font-family: "Avenir Next LT Pro", Helvetica, Arial, sans-serif;
    padding: 0 20px;
  }

  .content {
    margin-top: 70px;
  }

  .content h3,
  .content h4,
  .content h5 {
    font-weight: bold;
    text-transform: uppercase;
  }
</style>

<template lang="pug">
  div
    nav-bar(v-if="this.$route.path !== '/login'")
    div.container-fluid.content
      alert(:text="errorMsg" v-on:close="this.closeAlert")
      router-view(v-on:error="this.handleError")
</template>

<script>
  import Navbar from "./navbar.vue"
  import Alert from "./alert.vue"
  export default {
    components: {
      "alert": Alert,
      "nav-bar": Navbar
    },
    methods: {
      handleError(message) {
        this.errorMsg = message
      },

      closeAlert() {
        this.errorMsg = ""
      }
    },
    data() {
      return {
        errorMsg: ""
      }
    }
  }
</script>
