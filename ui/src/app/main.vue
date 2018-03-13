<style lang="scss">
  @import "../scss/app.scss"
</style>

<template>
  <div>
    <HeaderBar :node-info="nodeInfo" v-on:openRecoveryAlert="openRecoveryAlert"></HeaderBar>
    <router-view :node-info="nodeInfo" v-on:error="this.handleError"></router-view>
    <RecoveryAlert v-if="showRecoveryAlert" v-on:showPhrase="openRecoveryPhrase" v-on:close="closeRecoveryAlert"></RecoveryAlert>
    <RecoveryPhrase v-if="showRecoveryPhrase" v-on:close="closeRecoveryPhrase"></RecoveryPhrase>
  </div>
</template>

<script>
  import axios from "axios"
  import RecoveryAlert from '../components/modal/recoveryAlert.vue'
  import RecoveryPhrase from '../components/modal/recoveryPhrase.vue'

  const HeaderBar = require('./header.vue')

  export default {
    components: {
      HeaderBar: HeaderBar,
      RecoveryAlert: RecoveryAlert,
      RecoveryPhrase: RecoveryPhrase,
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

      openRecoveryAlert() {
        if (this.showRecoveryPhrase) {
          return
        }
        this.showRecoveryAlert = true
      },

      openRecoveryPhrase() {
        this.showRecoveryAlert = false
        this.showRecoveryPhrase = true
      },

      closeRecoveryAlert() {
        this.showRecoveryAlert = false
      },

      closeRecoveryPhrase() {
        this.showRecoveryPhrase = false
      }
    },

    data() {
      return {
        showRecoveryAlert: false,
        showRecoveryPhrase: false,
        nodeInfo: {}
      }
    }
  }
</script>
