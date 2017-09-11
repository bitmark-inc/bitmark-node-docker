<style scoped>
  h3 {
    font-size: 18px;
    font-weight: bold;
  }

  .container {
    margin: 0 auto;
    text-align: center;
    text-transform: uppercase;
  }

  .login-form {
    margin-top: 100px;
    display: inline-block;
  }

  .login-form .password,
  .login-form .login-btn {
    font-size: 16px;
    text-transform: uppercase;
    text-align: center;
    width: 400px;
    height: 50px;
    border: none;
    margin-top: 20px;
  }

  .login-form .password {
    background-color: #edf0f4;
    font-style: italic;
  }

  .login-form .login-btn {
    color: white;
    padding: 1px;
    background-color: #0060f2;
  }

  .login-form .login-btn:hover {
    background-color: black;
  }
</style>

<template lang="pug">
    div.container
      form.login-form(@submit="this.login")
        h3 login bitmark node web
        div
          input.password(type="password", v-model="password" placeholder="PASSWORD")
        div
          button.login-btn(type="submit") login
</template>

<script>
  const axios = require("axios");
  import {
    getCookie
  } from "../utils"

  export default {
    methods: {
      login(e) {
        e.preventDefault();

        let redirect = this.$route.query.redirect || "/"

        if (getCookie("bitmark-webgui")) {
          this.$router.push({
            path: redirect
          })
        }
        axios.post('/api/login', {
            password: this.password
          })
          .then((response) => {
            if (response.data.ok) {
              this.$router.push({
                path: redirect
              })
            } else {
              throw new Error(response.data.result);
            }
          })
          .catch((e) => {
            this.$emit("error", e.message)
          });
      }
    },

    data() {
      return {
        password: ""
      }
    }
  }
</script>
