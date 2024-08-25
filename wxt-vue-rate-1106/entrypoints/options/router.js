import { createRouter, createWebHashHistory } from 'vue-router'
import Home from './views/Home.vue'

export default createRouter({
  // history: createWebHistory(),
  "history": createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    // {
    //   path: '/add_list',
    //   name: 'AddList',
    //   component: () => import('./views/AddList.vue'),
    // },

  ],
})