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
    h5 bitmark rpc
    div.row
      div.col-xs-6.col-sm-4 CHAIN
      div.upper.col-xs-6.col-sm-8 {{bitmarkdConfig.chain}}
    div.row
      div.col-xs-6.col-sm-4 NODES
      div.col-xs-6.col-sm-8
        input.input-form(v-model="bitmarkdConfig.nodes")
    div.row
      div.col-xs-6.col-sm-4 Announce
      div.col-xs-6.col-sm-8
        input.input-form(v-model="bitmarkdConfig.client_rpc.announce[0]")
  div
    h5 bitmark peer
    div.row
      div.col-xs-6.col-sm-4 PUBLICKEY
      div.upper.col-xs-6.col-sm-8 {{bitmarkdConfig.peering.public_key}}
    div.row
      div.col-xs-6.col-sm-4 BROADCAST
      div.col-xs-6.col-sm-8
        input.input-form(v-model="bitmarkdConfig.peering.announce.broadcast[0]")
    div.row
      div.col-xs-6.col-sm-4 LISTEN
      div.col-xs-6.col-sm-8
        input.input-form(v-model="bitmarkdConfig.peering.announce.listen[0]")
  div
    h5 bitmark proofing
    div.row
      div.col-xs-6.col-sm-4 PUBLICKEY
      div.upper.col-xs-6.col-sm-8 {{bitmarkdConfig.proofing.public_key}}
    div.row
      div.col-xs-6.col-sm-4 CURRENCY
      div.col-xs-6.col-sm-8
        input.input-form(v-model="bitmarkdConfig.proofing.currency")
    div.row
      div.col-xs-6.col-sm-4 ADDRESS
      div.col-xs-6.col-sm-8
        input.input-form(v-model="bitmarkdConfig.proofing.address")
    div.row
      div.col-xs-6.col-sm-4 PUBLISH
      div.col-xs-6.col-sm-8
        input.input-form(v-model="bitmarkdConfig.proofing.publish[0]")
    div.row
      div.col-xs-6.col-sm-4 SUBMIT
      div.col-xs-6.col-sm-8
        input.input-form(v-model="bitmarkdConfig.proofing.submit[0]")
  div
    h5 prooferd peer
    div.row
      div.col-xs-6.col-sm-4 CONNNECT
        a.add(@click="this.addConnect") +
      div.upper.col-xs-6.col-sm-8
        table
          tbody
            template(
              v-for="(conn, index) in prooferdConfig.peering.connect"
            )
              tr
                td \#{{ index+1 }}
                td
                td
              tr
                td PUBLICKEY
                td
                  input.input-form(v-model="conn.public_key")
                td(rowspan="3", @click="() => {removeConnect(index)}") x
              tr
                td BLOCKS
                td
                  input.input-form(v-model="conn.blocks")
                td
              tr
                td SUBMIT
                td
                  input.input-form(v-model="conn.submit")
                td

</template>

<script>
  import axios from "axios"
  import {
    getCookie,
    setCookie
  } from "../utils"

  function checkAndMergeConfig(oldConf, newConf) {
    let _conf = new oldConf.constructor()
    for (var k in oldConf) {
      if (newConf[k]) {
        let v = newConf[k]
        // The order of type checking here is important
        if (oldConf[k] instanceof Array) {
          _conf[k] = (newConf[k] instanceof Array && newConf[k].length > 0) ? newConf[k] : oldConf[k]
        } else if (oldConf[k] instanceof Object) {
          _conf[k] = checkAndMergeConfig(oldConf[k], newConf[k])
        } else {
          _conf[k] = v
        }
      } else {
        _conf[k] = oldConf[k]
      }
    }
    return _conf
  }

  export default {
    methods: {
      addConnect() {
        this.prooferdConfig.peering.connect.push({})
      },

      removeConnect(index) {
        this.prooferdConfig.peering.connect.splice(index, 1)
      },

      save() {
        axios.post("/api/config", {
            bitmarkConfig: this.bitmarkdConfig,
            prooferdConfig: this.prooferdConfig
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
      axios
        .get("/api/config")
        .then((response) => {
          let {
            bitmarkd,
            prooferd
          } = response.data.result

          if (bitmarkd.err !== "") {
            this.$emit("error", bitmarkd.err)
            return
          }

          if (prooferd.err !== "") {
            this.$emit("error", prooferd.err)
            return
          }
          this.bitmarkdConfig = checkAndMergeConfig(this.bitmarkdConfig, bitmarkd.data)
          this.prooferdConfig = checkAndMergeConfig(this.prooferdConfig, prooferd.data)
        })
        .catch((e) => {
          console.log(e)
        })
    },

    data() {
      return {
        bitmarkdConfig: {
          chain: getCookie("bitmark-webgui-network"),
          nodes: "chain",
          client_rpc: {
            announce: ["127.0.0.1:2130"]
          },
          peering: {
            public_key: "",
            announce: {
              broadcast: ["127.0.0.1:2135"],
              listen: ["127.0.0.1:2136"]
            }
          },
          proofing: {
            public_key: "",
            currency: "bitcoin",
            address: "",
            publish: ["127.0.0.1:2140"],
            submit: ["127.0.0.1:2141"]
          }
        },
        prooferdConfig: {
          peering: {
            connect: []
          }
        }
      }
    }
  }
</script>
