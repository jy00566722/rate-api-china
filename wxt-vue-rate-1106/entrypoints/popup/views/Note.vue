<template>
    <div class="note" v-if="noteInfo.main">
          <el-alert
          :title="noteInfo.context"
          :type="noteInfo.type"
          center
          @close="closeHander"
          >
        </el-alert>
        </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, computed ,onBeforeMount} from 'vue'
import { useMessageStore } from '../pinia'


const noteInfo=ref({
    "context": "各站点页面价格转换,需要手动在设置中开启",
    "index": 1,
    "main": false,
    "type": "info"
})
const messageStore = useMessageStore()

onBeforeMount(async()=>{
  let r = await chrome.storage.local.get(['closeMessage'])
  
  if(getChromeVersion()){
    noteInfo.value.main = true
    noteInfo.value.context = '浏览器版本低于102,请升级!'
  }else{
        let url = 'https://rate.lizudi.top/v2/message'
        let a = await fetch(url)
        let b = await a.json()
        if(b.code===20000){
          b.data.forEach(item => {
            if(item.main){
              if(r.closeMessage&& item.index===r.closeMessage.messageId && r.closeMessage.closeTime === getToday()){
                return
              }
              noteInfo.value.main = true
              noteInfo.value.context = item.context
              noteInfo.value.index = item.index
              noteInfo.value.type = item.type
            }
          });
        //保存进入pinia
        messageStore.message = b.data
        }else{
          console.log('出错了啊',b)
        }
  }
})

//返回浏览器Chrome版本是否大于等于102
const getChromeVersion=()=>{
let  ua = navigator.userAgent;
let chrome = /Chrome\/([\d]+)./i.exec(ua);
let version = parseInt(chrome[1])
// console.log(version)
return version < 102
}
const getToday=()=>{
    const today = new Date();
    const formattedDate = today.toLocaleDateString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit'
          });
        return formattedDate
}

const closeHander=()=>{
  let closeMessage = {
    messageId: noteInfo.value.index,
    message: noteInfo.value.context,
    closeTime: getToday()
  }
  chrome.storage.local.set({closeMessage},()=>{} )

}
</script>

<style>
    .note {
      margin: 2px 10px 2px 10px;
      /* background-color: rgb(39, 46, 56); */
      
    }
    .note .el-alert{
       padding: 4px 12px 4px 12px; 
    }
</style>