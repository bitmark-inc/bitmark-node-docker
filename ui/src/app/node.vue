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

  .info-box span.error {
    color: red;
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
        span(v-if="this.bitmarkd.status === 'started'") Loading bitmarkd info…
        span(v-else-if="this.bitmarkd.status === 'stopped'")
          span(v-if="this.bitmarkd.errorMsg", class="error") {{ this.bitmarkd.errorMsg }}
          span(v-else) Bitmarkd is stopped
        span(v-else-if="this.bitmarkd.status === ''") Checking bitmarkd status…
        p(v-else) Bitmarkd is failed to start: {{ this.bitmarkd.status }}

    h4 recorderd node
      div.action
        button.btn(@click="this.startRecorderd"
          :disabled="!this.recorderd.status || this.recorderd.status!=='stopped'") Start
        button.btn.stop(disabled, @click="this.stopRecorderd"
          :disabled="!this.recorderd.status || this.recorderd.status==='stopped'") Stop
    p.info-box
      span(v-if="this.recorderd.status === ''") Checking recorderd status…
      span(v-else)
        span(v-if="this.recorderd.errorMsg", class="error") {{ this.recorderd.errorMsg }}
        span(v-else) Recorderd is {{this.recorderd.status || "loading status"}}
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
      }
      this.network = network;
      this.periodicalTask = setInterval(() => {
        this.fetchStatus('bitmarkd')
        this.fetchStatus('recorderd')
        this.fetchBitmarkInfo()
      }, 2000)
    },

    destroyed() {
      clearInterval(this.periodicalTask)
    },

    data() {
      return {
        network: "",
        periodicalTask: null,
        bitmarkd: {
          errorMsg: "",
          querying: false,
          status: ""
        },
        recorderd: {
          errorMsg: "",
          querying: false,
          status: ""
        },
        bitmarkdInfo: null
      }
    }
  }
</script>
