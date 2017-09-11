<style lang="" scoped>
  .fields {
    margin-top: 5px;
  }

  .key {
    text-transform: uppercase;
  }

  .row {
    width: 100%;
    text-align: left;
  }

  .horizontal .key {
    float:  left;
    width: 30%;
  }
</style>

<template lang="pug">
  div.row
    div.col-md-2
      h5 {{this.title}}
    div.col-md-10
      div.fields.row(v-for="(val, key) in this.data")
        div.col-sm-6.col-md-3.key {{ key }}
        div.col-sm-6.col-md-9.val(v-if="typeof val === 'object'")
          dl(:class="dlClass")
            template(v-for="(subVal, subKey) in val")
              template(v-if="typeof subVal !== 'object'")
                dt.key {{subKey}}
                dd {{subVal}}
              template(v-else, v-for="(v, k) in subVal")
                dt.key {{subKey}}.{{k}}
                dd {{v}}
        div.col-sm-6.col-md-9.val(v-else) {{ val }}
</template>

<script>
  export default {
    props: {
      subAlign: {
        type: String,
        default: function() {
          return "vertical"
        }
      },
      title: String,
      data: Object
    },

    computed: {
      dlClass() {
        let ret = {
          horizontal: this.subAlign === 'horizontal'
        }
        return ret
      }
    },

    data() {
      return {}
    }
  }
</script>
