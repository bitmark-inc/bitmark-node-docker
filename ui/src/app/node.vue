<template>
  <div class="bitmark-node-wrapper">
    <div class="content-body">
      <div class="row">
        <h3 class="paragraph-title">Bitmark node status</h3>
        <div class="row__box">
          <Box :class="{'running': this.bitmarkd.status.started, 'stop-running': !this.bitmarkd.status.started}" title="Bitmark Node (bitmarkd)">
            <button slot="header-button" class="btn-default" @click="toggleBitmarkd">{{ this.bitmarkd.status.started ? 'Stop' : 'Start' }}</button>
            <ul>
              <li v-if="!this.bitmarkd.status.error">
                <span class="label">Status:</span>
                <span class="status">{{ this.bitmarkd.status.started ? "Running" : "Stopped" }}</span>
              </li>
              <li v-else>
                <span class="label">Error: </span>
                <span class="status">{{ this.bitmarkd.status.error }}</span>
              </li>
              <li v-if="this.bitmarkd.status.running === 'started'">
                <span class="label">Connection:</span>
                <span class="status" v-if="bitmarkdConnStat !== null">You’re connected to {{ this.bitmarkdConnStat.connections }} nodes.
                  <span v-if="!bitmarkdConnStat.port_state.broadcast">
                    <br> Broadcast port (2135) is not accessible.</span>
                  <span v-if="!bitmarkdConnStat.port_state.listening">
                    <br> Listening port (2136) is not accessible.</span>
                </span>
                <span class="status" v-else>Checking networking…</span>
              </li>
            </ul>
          </Box>
          <!-- End: box -->
          <Box :class="{'running': this.recorderd.status.started, 'stop-running': !this.recorderd.status.started}" title="Recorder Node (recorderd)">
            <button slot="header-button" class="btn-default" @click="toggleRecorderd">{{ this.recorderd.status.started ? 'Stop' : 'Start' }}</button>
            <ul>
              <li>
                <span class="label ">Status:</span>
                <span class="status ">{{ this.recorderd.status.started ? "Running" : "Stopped" }}</span>
              </li>
            </ul>
          </Box>
          <!-- End: box -->
        </div>
      </div>
      <div class="divider "></div>
      <template v-if="bitmarkdInfo">
        <div class="row ">
          <h3 class="paragraph-title ">Bitmark node info</h3>
          <div class="row__box ">
            <Box title="Network ID ">
              <p>{{ this.bitmarkdInfo.public_key }}</p>
            </Box>
            <!-- End: box -->
            <Box title="Current Block ">
              <div class="blocks ">
                <span class="blocks__num ">{{ this.bitmarkdInfo.blocks }}/{{ this.bitmarkdInfo.block_height || this.bitmarkdInfo.blocks }}</span>
                <span class="blocks__label ">
                  <template v-if="this.bitmarkdInfo.mode === 'Resynchronise'">Updating blockchain</template>
                  <template v-else-if="this.bitmarkdInfo.mode === 'Normal'">Latest block</template>
                  <template v-else>Unknown block status</template>
                </span>
              </div>
            </Box>
            <!-- End: box -->
            <Box title="Transaction Counters ">
              <ul>
                <li>
                  <span class="label ">Pending:</span>
                  <span class="status ">{{ this.bitmarkdInfo.transactionCounters.pending }}</span>
                </li>
                <li>
                  <span class="label ">Verified: </span>
                  <span class="status ">{{ this.bitmarkdInfo.transactionCounters.verified }}</span>
                </li>
              </ul>
            </Box>
            <!-- End: box -->
            <Box title="Uptime ">
              <div class="times ">
                {{ this.bitmarkdInfo.uptime || '0s' }}
              </div>
            </Box>
            <!-- End: box -->
          </div>
        </div>
        <div class="divider "></div>
      </template>
      <div class="row ">
        <h3 class="paragraph-title ">Bitmark wallet</h3>
        <div class="row__box ">
          <Box class="full-width" title="Payment Addresses">
            <button slot="header-button" class="btn-default " @click="openConfig">Edit</button>
            <div class="btc-address complete ">
              <i>
                <img src="assets/img/icons/ic_bitcoin.svg " alt="icon bitcoin ">
              </i>
              <span class="coin-title ">BTC Address:</span>
              <span class="field ">
                <span>{{ paymentAddrs.btc || 'NOT SET' }}</span>
              </span>
            </div>
            <div class="ltc-address ">
              <i>
                <img src="assets/img/icons/ic_litecoin.svg " alt="icon litecoin ">
              </i>
              <span class="coin-title ">LTC Address:</span>
              <span class="field ">
                <span>{{ paymentAddrs.ltc || 'NOT SET' }}</span>
              </span>
            </div>
          </Box>
          <!-- End: box -->
        </div>
      </div>
      <div class="divider "></div>
      <PaymentPopUp v-if="showPaymentConfig" v-on:saved="saveConfig" v-on:close="closeConfig"
        :initBtcAddr="paymentAddrs.btc" :initLtcAddr="paymentAddrs.ltc"></PaymentPopUp>
    </div>
    <!-- End: content-body -->
  </div>
</template>

<script>
  import axios from "axios"

  import {
    getCookie,
    setCookie
  } from "../utils"
  import Box from './box.vue'
  import PaymentPopUp from '../components/paymentPopUp.vue'

  export default {
    components: {
      Box: Box,
      PaymentPopUp: PaymentPopUp
    },

    methods: {

      toggleBitmarkd() {
        if (this.bitmarkd.status.started) {
          this.stopBitmarkd()
        } else {
          this.startBitmarkd()
        }
      },

      startBitmarkd() {
        if (!this.paymentAddrs.btc || !this.paymentAddrs.ltc) {
          this.openConfig()
          return;
        }
        this.bitmarkd.status = ""
        this.bitmarkd.errorMsg = ""
        axios.post("/api/" + "bitmarkd", {
            option: "start"
          })
          .catch((err, resp) => {
            this.bitmarkd.errorMsg = err.response.data.msg
          })
      },

      stopBitmarkd() {
        this.bitmarkd.status = ""
        this.bitmarkd.errorMsg = ""
        this.bitmarkdInfo = null;
        axios.post("/api/" + "bitmarkd", {
            option: "stop"
          })
          .catch((err, resp) => {
            this.bitmarkd.errorMsg = err.response.data.msg
          })
      },

      toggleRecorderd() {
        if (this.recorderd.status.started) {
          this.stopRecorderd()
        } else {
          this.startRecorderd()
        }
      },

      startRecorderd() {
        this.recorderd.status = ""
        this.recorderd.errorMsg = ""
        axios.post("/api/" + "recorderd", {
            option: "start"
          })
          .catch((err, resp) => {
            this.recorderd.errorMsg = err.response.data.msg
          })
      },

      stopRecorderd() {
        this.recorderd.status = ""
        this.recorderd.errorMsg = ""
        axios.post("/api/" + "recorderd", {
            option: "stop"
          })
          .catch((err, resp) => {
            this.recorderd.errorMsg = err.response.data.msg
          })
      },

      getConfig() {
        return axios
          .get("/api/config")
          .then((response) => {
            let {
              btcAddr,
              ltcAddr
            } = response.data.result

            this.paymentAddrs.btc = btcAddr
            this.paymentAddrs.ltc = ltcAddr
          })
          .catch((e) => {
            console.log(e)
          })
      },
      openConfig() {
        this.showPaymentConfig = true
      },
      closeConfig() {
        this.showPaymentConfig = false
      },

      saveConfig(paymentAddrs) {
        axios.post("/api/config", {
            btcAddr: paymentAddrs.btcAddr,
            ltcAddr: paymentAddrs.ltcAddr
          })
          .then(() => {
            this.paymentAddrs.btc = paymentAddrs.btcAddr
            this.paymentAddrs.ltc = paymentAddrs.ltcAddr
            this.showPaymentConfig = false
          })
          .catch((e) => {
            this.$emit('error', e)
          })
      },

      getBitmarkdConnectionStatus() {
        if (!this.bitmarkd.status.started) {
          return
        }
        axios.get("/api/" + "bitmarkd/conn_stat")
          .then((resp) => {
            let data = resp.data
            if (resp.status === 200) {
              this.bitmarkdConnStat = data
            }
          })
      },

      fetchBitmarkInfo() {
        if (this.bitmarkd.status.started) {
          axios.post("/api/" + "bitmarkd", {
            option: "info"
          }).then((resp) => {
            let data = resp.data
            if (data.ok) {
              this.bitmarkdInfo = data.result
            }
          })
        }
      },

      fetchStatus(serviceName) {
        let service = this[serviceName]
        if (service.querying) {
          return
        }
        service.querying = true
        axios.post("/api/" + serviceName, {
            option: "status"
          })
          .then((resp) => {
            if (resp.data.ok) {
              service.status = resp.data.result
            } else {
              throw new Error(resp.data.result)
            }
          }).catch((e) => {
            this.$emit("error", e.message)
          })
          .then(() => {
            service.querying = false
          })
      }
    },

    created() {
      this.getConfig()
        .then(() => {
          let lastRun = this.lastRun
          if (!lastRun && (!this.paymentAddrs.btc || !this.paymentAddrs.ltc)) {
            this.showPaymentConfig = true;
            setCookie("lastRun", Date(), 30)
          }
        })
    },

    mounted() {
      let network = getCookie("bitmark-node-network") || 'bitmark'
      if (!network) {
        this.$router.push("/chain")
      }

      this.network = network;
      this.periodicalTask = setInterval(() => {
        this.fetchStatus('bitmarkd')
        this.fetchStatus('recorderd')
        this.fetchBitmarkInfo()
        this.getBitmarkdConnectionStatus()
      }, 2000)
    },

    destroyed() {
      clearInterval(this.periodicalTask)
    },

    data() {
      return {
        lastRun: getCookie("lastRun"),

        showPaymentConfig: false,
        paymentAddrs: {
          btc: "",
          ltc: ""
        },

        network: "",
        periodicalTask: null,
        bitmarkd: {
          errorMsg: "",
          querying: false,
          status: "",
          error: ""
        },
        recorderd: {
          errorMsg: "",
          querying: false,
          status: "",
          error: ""
        },

        bitmarkdInfo: null,
        bitmarkdConnStat: null
      }
    }
  }
</script>
