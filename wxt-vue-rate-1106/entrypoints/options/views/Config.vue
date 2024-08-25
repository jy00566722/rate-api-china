<template>
   <div class="config-set container">
    <el-card style="width: 400px;">
      <template #header>
        <div class="card-header">
          <span>刷新汇率</span>
          <!-- <el-button class="button" text>Operation button</el-button> -->
        </div>
      </template>
      <!-- <div v-for="o in 4" :key="o" class="text item">{{ 'List item ' + o }}</div> -->
       <span>点击刷新汇率,会重新获取汇率数据,并更新到本地</span>
      <el-button size="small"  type="primary"  @click="getRateAndNode">刷新汇率</el-button>
    </el-card>

    <el-card style="width: 400px;">
      <template #header>
        <div class="card-header">
          <span>设置页面(跨境平台)显示人民币的颜色</span>
        </div>
      </template>
      <!-- <div v-for="o in 4" :key="o" class="text item">{{ 'List item ' + o }}</div> -->
      <span>选择跨境平台页面显示人民币的颜色,尽量选鲜亮一点的便于识别:(选中后记得点OK)</span>
      <el-color-picker size="large" v-model="color" :predefine="predefineColors" @change="colorChange"/><br>
      <span>示例:  <span ref="selectedColor"> ¥99.9</span> </span>
    </el-card>

    <el-card style="width: 400px;">
      <template #header>
        <div class="card-header">
          <span>设置POP页面价格转换时:小数点位数 / 千分位</span>
        </div>
      </template>

      <span>最少2位小数,最大6位:</span>
      <el-input-number v-model="decimal" :step="1" step-strictly size="large" :min="2" :max="6" @change="saveDecimal()"/><br><br>
      <span>是否显示千分位:</span>
      <el-radio-group v-model="thousandsF"  size="small" @change="thousandsChange">
        <el-radio-button label="no" >不显示</el-radio-button>
        <el-radio-button label="yes" >显示</el-radio-button>
      </el-radio-group>
     
    </el-card>

    <el-card style="width: 400px;">
      <template #header>
        <div class="card-header">
          <span>设置POP页面主题</span>
        </div>
      </template>

      <span>默认随系统改变:</span>
      <el-radio-group v-model="theme" size="small" @change="themeChange">
        <el-radio-button label="light" >浅色系</el-radio-button>
        <el-radio-button label="dark" >深色系</el-radio-button>
        <el-radio-button label="excited" >清新蓝</el-radio-button>
        <el-radio-button label="ePurple" >优雅紫</el-radio-button>
        <el-radio-button label="classBlue" >【原】经典蓝</el-radio-button>
        <el-radio-button label="auto" >跟随系统</el-radio-button>

      </el-radio-group>
    </el-card>


    <el-card style="width: 400px;">
      <template #header>
        <div class="card-header">
          <span>设置POP布局 / 以及是否显示100货币量的转换</span>
        </div>
      </template>

      <span>默认为正常布局:</span>
      <el-radio-group v-model="layout" size="small" @change="layoutChange">
        <el-radio-button label="default" >正常布局</el-radio-button>
        <el-radio-button label="small" >紧凑布局</el-radio-button>
      </el-radio-group>
    
     
      <el-divider />
      <span>是否显示100货币量的转换:<br>(示例:<span style="color: red;font-size: 12px;">100USD = 728.56CNY</span>)</span>
      <el-radio-group v-model="show100" size="small" @change="show100Change">
        <el-radio-button label="yes" >显示</el-radio-button>
        <el-radio-button label="no" >不显示</el-radio-button>
      </el-radio-group><br>
      <span style="font-size: 12px;">紧凑布局时不显示</span>
    </el-card>
    <el-card style="width: 400px;">
      <template #header>
        <div class="card-header">
          <span>设置页面(跨境平台)转换为哪种货币(实验性质,默认CNY) 以及字体颜色 (目前只对美客多生效)</span>
        </div>
      </template>

      <span>转换为:</span>
      <el-select v-model="ConverterCurreny" placeholder="Select" @change="ConverterCurreny_change">
        <el-option
          v-for="item in optionsCurreny"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        />
      </el-select><br><br>
      <span>字体主题色(注意选深色 或 鲜亮的颜色):</span>   <el-color-picker size="large" v-model="currenyTheme" :predefine="predefineColors" @change="currenyThemeChange"/>
      <br><span>示例:</span>
      <span class="rmb" translate="no" style="display: inline-flex; border-radius: 8px; font-size: 14px; font-weight: 600; overflow: hidden;">
        <span :style="{ backgroundColor: currenyTheme, padding: '2px 2px 2px 5px', color: 'white', display: 'flex', alignItems: 'center', height: '16px', fontSize: '12px' }">
          {{ ConverterCurreny }}
        </span>
        <span :style="{ color: currenyTheme, padding: '2px 5px 2px 2px', backgroundColor: 'rgb(238, 238, 238)', display: 'flex', alignItems: 'center', height: '16px', fontSize: '14px' }">
          56.21
        </span>
      </span>
    </el-card>


   </div>
</template>

<script setup>
import { ref,onMounted,onBeforeMount } from 'vue'
import { ElMessage } from 'element-plus'

const theme = ref('auto') //POP页面主题
const layout = ref('default') //POP页面布局
const show100 = ref('yes') //是否显示100货币量的转换
const color = ref('#008000') //页面价格转换显示的颜色
const decimal = ref(2) //POP页面价格转换时的小数点位数
const thousandsF = ref('no') //是否显示千分位,
//预计的颜色
const predefineColors = ref([
  '#ff4500',
  '#ff8c00',
  '#ffd700',
  '#90ee90',
  '#00ced1',
  '#1e90ff',
  '#c71585',
  '#008000',
  '#008080',
  '#8a2be2'
])
const selectedColor = ref(null)

//记录在页面中要转换为的货币
const ConverterCurreny = ref('CNY')
//一个options数组，里面记录select的值，用于select组件的v-model
const optionsCurreny = ref([
  {value: 'CNY', label: '人民币'},
  {value: 'USD', label: '美元'},
  {value: 'EUR', label: '欧元'},
  {value: 'JPY', label: '日元'},
  {value: 'GBP', label: '英镑'},

]);



//监听ConverterCurreny的变化，并设置storage中的值
const ConverterCurreny_change=(v)=>{
  chrome.storage.local.set({ConverterCurreny: v}, function() {
  });
}

//记录页面价格的主题
const currenyTheme = ref('#8a2be2')
const currenyThemeChange=()=>{
  chrome.storage.local.set({currenyTheme: currenyTheme.value}, function() {
    console.log('currenyTheme is ' + currenyTheme.value);
  });
}


//给后台发信息,请求最新汇率及节点信息
async function getRateAndNode() {
  const response = await chrome.runtime.sendMessage({do: "getStart",cache:"no"});
  if (response.status==='good') {
    ElMessage.success('刷新成功');
  }
  //获取storage中的所有数据,打印出来
  chrome.storage.local.get(null, function(items) {
    console.log('Storage中的数据:',items)
  });
}
const colorChange=()=>{
  selectedColor.value.style.color = color.value
  chrome.storage.local.set({color: color.value}, function() {
    console.log('color is ' + color.value);
  });
}
const saveDecimal=()=>{
  chrome.storage.local.set({decimal: decimal.value}, function() {
    console.log('decimal位数是 : ' + decimal.value);
  });
}
const themeChange=()=>{
  chrome.storage.local.set({theme: theme.value}, function() {
    console.log('theme is ' + theme.value);
  });
  //根据theme.value设置主题
  if(theme.value==='auto'){
    // console.log('auto')
    let theme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    // console.log(theme)
    // document.documentElement.setAttribute('data-theme',theme);
  }else{
    // console.log('light or dark')
    // document.documentElement.setAttribute('data-theme',theme.value)
  }
}
const layoutChange=()=>{
  chrome.storage.local.set({layout: layout.value}, function() {
    console.log('layout is ' + layout.value);
    // document.documentElement.setAttribute('data-layout',layout.value);
  });
}
const show100Change=()=>{
  chrome.storage.local.set({show100: show100.value}, function() {
    console.log('show100 is ' + show100.value);
  });
}
const thousandsChange=()=>{
  chrome.storage.local.set({thousandsF: thousandsF.value}, function() {
    console.log('thousandsF is ' + thousandsF.value);
  });
}
onMounted(async() => {
  // selectedColor.value = document.querySelector('#selectedColor')
  const response = await chrome.storage.local.get(['color','decimal','theme','layout','thousandsF','ConverterCurreny','currenyTheme','show100']);
  if (response.color) {
    color.value = response.color
    selectedColor.value.style.color = color.value
  }else{
    color.value = '#008000'
    selectedColor.value.style.color = color.value
  }
  if (typeof response.decimal == 'number' && response.decimal >= 2 && response.decimal <= 6) {
    decimal.value = response.decimal
  }else{
    decimal.value = 2
    //删除原来的decimal
    await chrome.storage.local.remove('decimal')
  }
  if(response.show100&&response.show100==='no'){
    show100.value = 'no'
  }else{
    show100.value = 'yes'
  }
  if(response.thousandsF&&response.thousandsF==='yes'){
    thousandsF.value = 'yes'
  }else{
    thousandsF.value = 'no'
  }
  if(response.theme=='light'||response.theme=='dark'||response.theme=='auto'){
    theme.value = response.theme
  }
  if(response.layout=='default'||response.layout=='small'){
    layout.value = response.layout
  }
  if(response.ConverterCurreny){
    ConverterCurreny.value = response.ConverterCurreny
  }else{
    ConverterCurreny.value = 'CNY'
  }
  if(response.currenyTheme){
    currenyTheme.value = response.currenyTheme
  }

  
})
</script>

<style scoped>
  .config-set{
    display: flex;
    flex-wrap: wrap;
    gap: 20px; 
  }
  .card-header {
    background-color: rgb(73, 71, 71);
    color: white;
    padding: 5px 10px;
  }
</style>