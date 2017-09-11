<style lang="" scoped>
  .chain-content {
    max-width: 730px;
    margin: 0 auto;
  }

  .box {
    border: none;
    background-color: rgb(237, 240, 244);
  }

  .panel-heading .sub {
    font-size: 12px;
    font-style: italic;
  }

  .option {
    font-weight: normal;
  }

  .option input {
    margin-right: 5px;
  }

  .option .help-text {
    font-size: 12px;
    margin: 4px 17px;
  }

  .panel-body {
    margin: 25px 0;
  }

  .panel-footer {
    border: none;
    background: none;
    text-align: right;
  }

  .start-node {
    display: inline-block;
    border: none;
    padding: 9px 24px;
    background-color: white;
    color: rgb(0, 96, 242);
  }

  .start-node:hover {
    color: black;
  }
</style>

<template lang="pug">
div.chain-content
  h4 start bitmark node
  div.panel.panel-default.box
    div.row
      div.col-md-3
        div.panel-heading
          h5.title select chain
          p.sub Bitmark provide two diffrent chains to let the bitmarkd join in. They are testing, bitmark.
      div.col-md-9
        div.panel-body
          label.option
            input(type="radio", value="testing", v-model="network")
            .
              TESTING
            p.help-text Link to public test bitmark network, to pay the transactions, please contact us.
          label.option
            input(type="radio", value="bitmark", v-model="network")
            .
              BITMARK
            p.help-text Link to public bitmark network, pay the transactions with real bitcoin.
        div.panel-footer
          button.start-node(@click="start") START NODE »

</template>

<script>
  import axios from "axios"
  import {
    setCookie,
    getCookie
  } from "../utils"

  export default {
    methods: {
      start() {
        if (!this.network) {
          this.$emit("error", "no network is selected")
          return
        }
        axios.post("/api/bitmarkd", {
            "option": "setup",
            "network": this.network
          })
          .then((resp) => {
            if (resp.data && resp.data.ok) {
              return axios.post("/api/prooferd", {
                "option": "setup",
                "network": this.network
              })
            } else {
              throw new Error('fail to setup bitmarkd service: ' + resp.data.result)
            }
          })
          .then((resp) => {
            if (resp.data && resp.data.ok) {
              setCookie("bitmark-webgui-network", this.network, 30)
              this.$router.push("/node")
            } else {
              throw new Error('fail to setup prooferd service: ' + resp.data.result)
            }
          })
          .catch((e) => {
            this.$emit("error", e.message)
          })
      }
    },
    data() {
      return {
        network: getCookie("bitmark-webgui-network")
      }
    }
  }
</script>
