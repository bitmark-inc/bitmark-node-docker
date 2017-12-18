var Vue = require('vue')
var VueRouter = require('vue-router')
Vue.use(VueRouter)

var Main = require('./app/main.vue')
// var Login = require('./app/login.vue')
var Node = require('./app/node.vue')
var NotFoundView = require('./app/notFound.vue')
// var Console = require('./app/console.vue')

import axios from "axios";
import {getCookie, setCookie} from "./utils"

var routes = [
  {path: '/', component: Main, component: Node},
  {path: '/node', component: Node},
  {path: '*', component: NotFoundView}
]

var router = new VueRouter({routes, linkActiveClass: "active"})

// Set chain cookie on start up
axios.get("/api/chain")
  .then((resp) => {
    var app = new Vue({
      router,
      el: '#main',
      render: h => h(Main)
    })
  })
