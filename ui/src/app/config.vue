<style scoped>
  .input-form {
    width: 280px;
    border: 0.8px solid #EDF0F4;
    max-width: 280px;
    min-width: 280px;
    height: 30px;
    min-height: 30px;
    max-height: 280px;
    padding-left: 10px;
    margin-top: 4px;
    background-color: white;
  }

  .row {
    min-height: 34px;
  }

  .add {
    text-decoration: none;
  }

  .upper {
    text-transform: uppercase;
  }

  .action {
    margin-top: -5px;
    float: right;
  }

  h3 .action .btn {
    box-shadow: none;
    border: none;
    margin-right: 10px;
    background: none;
    color: rgb(0, 96, 242);
    text-transform: uppercase;
    font-size: 16px;
    font-weight: bold;
    text-decoration: none;
  }

  h3 .action .btn:hover {
    color: rgb(126, 211, 33);
  }

  h3 .action .btn.cancel:hover {
    color: red;
  }

  h3 .action .btn[disabled],
  h3 .action .btn[disabled]:hover {
    color: rgb(193, 193, 193);
    cursor: not-allowed;
  }
</style>

<template lang="pug">
div
  h3 configuration
    div.action
      button.btn(@click="this.save") Save
      router-link(tag="button", class="btn cancel",to="/node") Cancel
  div
    h5 bitmark proofing
    div.row
      div.col-xs-6.col-sm-4 BTC Address
      div.col-xs-6.col-sm-8
        input.input-form(v-model="btcAddr")
    div.row
      div.col-xs-6.col-sm-4 LTC Address
      div.col-xs-6.col-sm-8
        input.input-form(v-model="ltcAddr")
</template>

<script>
  import axios from "axios"
  import {
    getCookie,
    setCookie
  } from "../utils"


  export default {
    methods: {
      save() {
        axios.post("/api/config", {
            btcAddr: this.btcAddr,
            ltcAddr: this.ltcAddr
          })
          .then(() => {
            this.$router.push("/node")
            console.log("saved")
          })
          .catch((e) => {
            this.$emit('error', e)
          })
      }
    },
    mounted() {
      let network = getCookie("bitmark-node-network")
      if (!network) {
        this.$router.push("/chain")
        return
      }

      axios
        .get("/api/config")
        .then((response) => {
          let {
            btcAddr,
            ltcAddr
          } = response.data.result

          this.btcAddr = btcAddr
          this.ltcAddr = ltcAddr
        })
        .catch((e) => {
          console.log(e)
        })
    },

    data() {
      return {
        ltcAddr: "",
        btcAddr: ""
      }
    }
  }
</script>
