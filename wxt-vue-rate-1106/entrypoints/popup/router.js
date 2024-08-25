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
    {
      path: '/add_list',
      name: 'AddList',
      component: () => import('./views/AddList.vue'),
    },
    {
      path: '/feed_back',
      name: 'FeedBack',
      component: () => import('./views/FeedBack.vue'),
    },
    {
      path: '/note_page',
      name: 'NotePage',
      component: () => import('./views/NotePage.vue'),
    },
  ],
})