<style scoped>
  .action {
    float: right;
  }

  .action .btn {
    border: none;
    margin: 0;
    padding: 0 10px;
    height: 100%;
    opacity: 1;
    background: none;
    color: rgb(0, 96, 242);
    text-transform: uppercase;
    font-size: 16px;
    font-weight: bold;
    text-decoration: none;
  }

  .action .btn:hover {
    color: rgb(126, 211, 33);
  }

  .action .btn.stop:hover {
    color: red;
  }

  .action .btn[disabled],
  .action .btn[disabled]:hover {
    color: rgb(193, 193, 193);
    cursor: not-allowed;
  }

  .row {
    padding-bottom: 10px;
  }

  .fields {
    margin-top: 5px;
  }

  .info-box {
    background-color: rgb(249, 251, 255);
    text-align: center;
    padding: 15px;
  }

  .info-box>span {
    text-transform: uppercase;
    font-weight: bold;
  }

  .info-box>span p{
    text-transform: none;
    font-weight: normal;
  }
</style>

<template lang="pug">
  div
    h4 current chain
    p.info-box
      span {{this.network}}
    h4 bitmark node
      div.action
        button.btn(
          @click="this.startBitmarkd"
          :disabled="!this.bitmarkd.status || this.bitmarkd.status!=='stopped'") Start
        button.btn.stop(
          disabled, @click="this.stopBitmarkd"
          :disabled="!this.bitmarkd.status || this.bitmarkd.status==='stopped'") Stop
    p.info-box
      status-grid(
        v-if="this.bitmarkdInfo", title="bitmark info", :style='{backgroundColor: "rgb(249, 251, 255)"}'
        :data="this.bitmarkdInfo", sub-align="horizontal")
      span(v-else)
        span(v-if="this.bitmarkd.status === 'started'") Bitmarkd info is not available
        span(v-else-if="this.bitmarkd.status === 'stopped'") Bitmarkd is stopped
        span(v-else-if="this.bitmarkd.status === ''") Checking bitmarkd status...
        p(v-else) Bitmarkd is failed to start: {{ this.bitmarkd.status }}

    h4 prooferd node
      div.action
        button.btn(@click="this.startProoferd"
          :disabled="!this.prooferd.status || this.prooferd.status!=='stopped'") Start
        button.btn.stop(disabled, @click="this.stopProoferd"
          :disabled="!this.prooferd.status || this.prooferd.status==='stopped'") Stop
    p.info-box
      span(v-if="this.prooferd.status === ''") Checking prooferd status...
      span(v-else) Prooferd is {{this.prooferd.status || "loading status"}}
</template>

<script>
  import {
    getCookie
  } from "../utils"
  import axios from "axios"

  import statusGrid from "../components/statusGrid.vue"

  export default {
    components: {
      "status-grid": statusGrid
    },
    methods: {
      startBitmarkd() {
        this.bitmarkd.status = ""
        axios.post("/api/" + "bitmarkd", {
          option: "start"
        })
      },

      stopBitmarkd() {
        this.bitmarkd.status = ""
        this.bitmarkdInfo = null;
        axios.post("/api/" + "bitmarkd", {
          option: "stop"
        })
      },

      startProoferd() {
        this.prooferd.status = ""
        axios.post("/api/" + "prooferd", {
          option: "start"
        })
      },

      stopProoferd() {
        this.prooferd.status = ""
        axios.post("/api/" + "prooferd", {
          option: "stop"
        })
      },

      fetchBitmarkInfo() {
        if (this.bitmarkd.status === "started") {
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

    mounted() {
      let network = getCookie("bitmark-node-network")
      if (!network) {
        this.$router.push("/chain")
      } else {
        axios.post("/api/chain", {
            "network": network
          })
          .then((resp) => {
            this.network = network;
             this.periodicalTask = setInterval(() => {
              this.fetchStatus('bitmarkd')
              this.fetchStatus('prooferd')
              this.fetchBitmarkInfo()
            }, 2000)
          })
      }
    },

    destroyed() {
      clearInterval(this.periodicalTask)
    },

    data() {
      return {
        network: "",
        periodicalTask: null,
        bitmarkd: {
          querying: false,
          status: ""
        },
        prooferd: {
          querying: false,
          status: ""
        },
        bitmarkdInfo: null
      }
    }
  }
</script>
