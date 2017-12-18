<style lang="scss">
  @import "../scss/app.scss"
</style>

<template>
  <div>
    <HeaderBar :node-info="nodeInfo"></HeaderBar>
    <router-view v-on:error="this.handleError"></router-view>
  </div>
</template>

<script>
  import axios from "axios"

  const HeaderBar = require('./header.vue')

  export default {
    components: {
      HeaderBar: HeaderBar,
    },

    created() {
      axios.get("/api/" + "info")
        .then((resp) => {
          let data = resp.data
          if (data.ok) {
            this.nodeInfo = data.result
          }
        })
    },

    methods: {
      handleError(message) {
        this.errorMsg = message
      },
    },

    data() {
      return {
        nodeInfo: {}
      }
    }
  }
</script>
