var Vue = require('vue')
var VueRouter = require('vue-router')
Vue.use(VueRouter)

var Main = require('./app/main.vue')
var Login = require('./app/login.vue')
var Chain = require('./app/chain.vue')
var Node = require('./app/node.vue')
var Config = require('./app/config.vue')
var Console = require('./app/console.vue')

import axios from "axios";
import {getCookie, setCookie} from "./utils"

var routes = [
  {path: '/', component: Main, redirect: '/node'},
  {path: '/chain', component: Chain},
  {path: '/node', component: Node},
  {path: '/config', component: Config},
]

var router = new VueRouter({routes, linkActiveClass: "active"})

var app = new Vue({
  router,
  el: '#main',
  render: h => h(Main)
})
