<style scoped>
  .console {
    width: 100%;
    height: 300px;
    border: none;
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

  iframe {
    background-color: black;
  }
</style>

<template lang="pug">
div
  iframe.console(src="/console")
  p.info-box
    status-grid(
      v-if="this.bitmarkdInfo", title="bitmark info", :style='{backgroundColor: "rgb(249, 251, 255)"}'
      :data="this.bitmarkdInfo", sub-align="horizontal")
    span(v-else) Bitmarkd info is not available

</template>

<script>
  import axios from "axios"
  import statusGrid from "../components/statusGrid.vue"

  export default {
    components: {
      "status-grid": statusGrid
    },

    methods: {
      fetchBitmarkInfo() {
        axios.post("/api/" + "bitmarkd", {
          option: "info"
        }).then((resp) => {
          let data = resp.data
          if (data.ok) {
            this.bitmarkdInfo = data.result
          }
        })
      },
    },

    mounted() {
      this.periodicalTask = setInterval(() => {
        this.fetchBitmarkInfo()
      }, 2000)
    },

    destroyed() {
      clearInterval(this.periodicalTask)
    },

    data() {
      return {
        periodicalTask: null,
        bitmarkdInfo: null
      }
    }
  }
</script>
