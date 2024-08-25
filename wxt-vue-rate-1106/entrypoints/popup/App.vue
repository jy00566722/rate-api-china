<template>


    <router-view/>
</template>
<script setup>
import { ref,onMounted,onBeforeMount } from 'vue'
import { useMessageStore } from './pinia'

const messageStore = useMessageStore()

const theme = ref('')

onBeforeMount(()=>{
//从storage中获取主题,theme.theme可能的值为light,dark,auto,如果为auto,则根据系统主题来设置,如果为light或者dark,则直接设置.如果没有获取到，则默认为auto.为auto时用window.matchMedia来获取系统主题。
//按获取的主题设置主题,设置方法为获取html上的data-theme属于，然后设置其值为light或dark
chrome.storage.local.get(["theme","layout"],(result)=>{
  if(result.theme){
    theme.value = result.theme
  }else{
    theme.value = 'auto'
  }
  // console.log(theme.value)
  if(theme.value==='auto'){
    // console.log('auto')
    let theme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    // console.log(theme)
    document.documentElement.setAttribute('data-theme',theme);
  }else{
    // console.log('light or dark')
    document.documentElement.setAttribute('data-theme',theme.value)
  }
  if(result.layout){
    document.documentElement.setAttribute('data-layout',result.layout);
    messageStore.layout = result.layout
  }else{
    document.documentElement.setAttribute('data-layout','default');
    messageStore.layout = 'default'
  }
})
})




onMounted(async() => {
      //打开一次记录一次
      let url = 'https://rate.lizudi.top/v2/feed_pop_page'
      let a = await fetch(url)
      let b = await a.json()
      // console.log(b)

    })
</script>
<style>
body{
  margin: 0;
  padding: 0;
  /* display: flex; */
  align-items:center;
  justify-content:center;
  /* background-color:rgb(39,46,56); */
  background: var(--main-backgroud-color);
  width: var(--body-width);
}
#app {
  width: var(--app-width);
  /* background-color:var(--main-backgroud-color); */
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  /* color: #2c3e50; */
}

body::-webkit-scrollbar {
    width: 5px;
    height: 10px;
  }

  body::-webkit-scrollbar-thumb {
    border-radius: 10px;
    -webkit-box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.2);
    box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.2);
    background: var(--scrollbar-thumb-color);
  }

  body::-webkit-scrollbar-track {
    -webkit-box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.2);
    box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.2);
    border-radius: 2px;
    background: var(--scrollbar-track-color);
  }

/* 

.nav a.router-link-exact-active {
  color: #42b983;
}
.add_country_button{
   position:fixed;
   bottom: 10px;
   right: 10px;
   z-index:999;
} */
</style>
