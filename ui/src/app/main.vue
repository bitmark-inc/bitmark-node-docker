<style lang="scss">
  @import "../scss/app.scss";
  .fullscreen {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
  }
</style>

<template>
  <div>
    <HeaderBar :node-info="nodeInfo" v-on:openRecoveryAlert="openRecoveryAlert"></HeaderBar>
    <router-view :node-info="nodeInfo" :payment-addrs="paymentAddrs" v-on:openPaymentConfig="openPaymentConfig" v-on:error="this.handleError"></router-view>
    <WelcomeDialog v-if="showWelcomePage" v-on:createAccount="createAccount" v-on:existingAccount="existingAccount"></WelcomeDialog>
    <EnterPhraseDialog v-if="showEnterPhrase" v-on:backToWelcome="backToWelcome" v-on:accountRecovered="accountRecovered"></EnterPhraseDialog>

    <RecoveryAlert v-if="showRecoveryAlert" v-on:showPhrase="openRecoveryPhrase" v-on:close="closeRecoveryAlert"></RecoveryAlert>
    <RecoveryPhrase v-if="showRecoveryPhrase" v-on:close="closeRecoveryPhrase"></RecoveryPhrase>

    <PaymentPopUp v-if="showPaymentConfig" v-on:saved="saveConfig" v-on:close="closePaymentConfig" :initBtcAddr="paymentAddrs.btc"
      :initLtcAddr="paymentAddrs.ltc"></PaymentPopUp>
  </div>
</template>

<script>
  import axios from "axios"
  import EnterPhraseDialog from '../components/modal/enterPhraseDialog.vue'
  import RecoveryAlert from '../components/modal/recoveryAlert.vue'
  import RecoveryPhrase from '../components/modal/recoveryPhrase.vue'
  import WelcomeDialog from '../components/modal/welcomeDialog.vue'
  import PaymentPopUp from '../components/paymentPopUp.vue'

  const HeaderBar = require('./header.vue')

  export default {
    components: {
      HeaderBar: HeaderBar,
      RecoveryAlert: RecoveryAlert,
      RecoveryPhrase: RecoveryPhrase,
      WelcomeDialog: WelcomeDialog,
      EnterPhraseDialog: EnterPhraseDialog,
      PaymentPopUp: PaymentPopUp
    },

    async created() {
      await this.checkAccount()
    },

    methods: {
      async checkAccount() {
        try {
          await axios.get("/api/" + "account");
          this.getNodeInfo();
        } catch (err) {
          this.showWelcomePage = true
        }
      },

      getConfig() {
        return axios
          .get("/api/config")
          .then(response => {
            let {
              btcAddr,
              ltcAddr
            } = response.data.result

            return {
              btcAddr,
              ltcAddr
            }
          })
      },

      saveConfig(paymentAddrs) {
        axios.post("/api/config", {
            btcAddr: paymentAddrs.btcAddr,
            ltcAddr: paymentAddrs.ltcAddr
          })
          .then(() => {
            this.paymentAddrs = {
              btc: paymentAddrs.btcAddr,
              ltc: paymentAddrs.ltcAddr
            }
            this.showPaymentConfig = false
          })
          .catch((e) => {
            console.log(e)
          })
      },

      async getNodeInfo() {
        let resp = await axios.get("/api/" + "info")
        let data = resp.data
        if (data.ok) {
          this.nodeInfo = data.result
          this.getPaymentAddr()
        }
      },

      async getPaymentAddr() {
        let {
          btcAddr,
          ltcAddr
        } = await this.getConfig();
        if (!btcAddr || !ltcAddr) {
          this.showPaymentConfig = true;
        } else {
          this.paymentAddrs = {
            btc: btcAddr,
            ltc: ltcAddr
          }
        }
      },

      openPaymentConfig() {
        this.showPaymentConfig = true;
      },

      closePaymentConfig() {
        this.showPaymentConfig = false;
      },

      handleError(message) {
        this.errorMsg = message
      },

      createAccount() {
        axios.post("/api/" + "account")
          .then(resp => {
            let data = resp.data
            if (data.ok) {
              this.getNodeInfo();
              this.showWelcomePage = false;
            }
          })
          .catch(err => {
            console.log(err)
          })
      },

      existingAccount() {
        this.showEnterPhrase = true;
        this.showWelcomePage = false;
      },

      backToWelcome() {
        this.showEnterPhrase = false;
        this.showWelcomePage = true;
      },

      accountRecovered() {
        this.showEnterPhrase = false;
        this.getNodeInfo();
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
        showWelcomePage: false,
        showEnterPhrase: false,
        showRecoveryAlert: false,
        showRecoveryPhrase: false,
        showPaymentConfig: false,
        paymentAddrs: {},
        nodeInfo: {}
      }
    }
  }
</script>
