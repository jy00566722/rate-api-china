<template>
    <div class="nav">
        <div class="nav_up">
          <span class="head_start" @click="$router.push('/feed_back')">反馈</span>
          <span class="head_start" style="color: var(--nav-text-note-color);" @click="$router.push('/note_page')">消息</span>
          <span class="head_start" style="color: var(--nav-text-sort-color);
          font-family:Verdana, Geneva, Tahoma, sans-serif;
          font-style: oblique;
          ">{{props.darpE==='no'?' ':'拖动旗帜排序'}}</span>
          <span @click="get_5starts" class="head_start">好评</span>
          <span @click="showOptionPage" class="head_start">设置</span>
          <span @click="popWindows" class="head_start">弹出</span>
        </div>

      </div>
</template>
<script setup>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
const props = defineProps({
  darpE:{
    type:String,
    default:'yes'
    // required: true
  }
}) //接受是否显示拖动提示的flag

//打开好评页面
const get_5starts = () => {
  let url =
          "https://chromewebstore.google.com/detail/%E6%B1%87%E7%8E%87%E8%BD%AC%E6%8D%A2/bcpgdpedphodjcjlminjbdeejccjbimp";

        if (navigator.userAgent.includes("Edg")) {
          url =
            "https://microsoftedge.microsoft.com/addons/detail/%E6%B1%87%E7%8E%87%E8%BD%AC%E6%8D%A2/jaippeddpgbjcdnmepklnjmcakmgodah";
        }
        window.open(url);
}

//打开option页面
const showOptionPage = () => {
  if (chrome.runtime.openOptionsPage) {
            chrome.runtime.openOptionsPage();
          } else {
            window.open(chrome.runtime.getURL('options.html'));
          }
}
//弹出桌面
const popWindows = async () => {
  let layout = await chrome.storage.local.get("layout");
  let width = 390;
  let height = 620;
  if (layout&&layout.layout == "small") {
        width = 295;
        height = 380;
      }
        
  chrome.windows.create({
            url: chrome.runtime.getURL("popup.html"),
            type: "popup",
            width,
            height,
          });
}
</script>
<style>

.nav{
  /* position:fixed; */
  top:0px;
  /* z-index: 1900; */
  display: grid;
  /* grid-template-rows: 40% 60%; */
  width: 100%;
  margin-right: 8px;
}
.nav_up{
  padding: 5px 0 0 0;
  height: 25px;
  display: grid;
  grid-template-columns: var(--nav-grid-template-columns);
  justify-items: end;

  /* background-color: rgb(39,46,56); */
}
.head_start {
      font-size: 14px;
      color: var(--nav-text-color);
    }
  
.head_start:hover {
  font-weight: 900;
  color: var(--nav-text-color-hover);
  cursor: pointer;
}
</style>