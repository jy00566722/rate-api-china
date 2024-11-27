<template>    
      <section class="hero is-info">
        <div class="hero-body">
          <div class="container">
            <h1 class="title">
              汇率转换 <span class="is-size-6 tag is-link p-1">{{version}}</span>
            </h1>
            <h3 class="subtitle">
                支持平台:Amazon / Lazada / Shopee / Qoo10.jp / Gmarket / Temu / WowmaJp / Rakuten / Aliexpress / Ozon / Vinted / 美客多 / Kream / Wildberries / Coupang / Yandex
              </h3>
            <h3 class="subtitle">
                问题及建议点此
                <a :href="ext_url" target="_blank">
                    <span class="tag is-warning">反馈,</span>
            </a>
              <span>或者加QQ群反应: 586113791 </span><span> 还可以发邮件: ideey88@gmail.com </span>
            </h3>
          </div>
        </div>
      </section>

      <div class="container" style="margin-top: 1rem;">
        <div class="box">
            <h4 class="title is-4">重要说明</h4>
            <p>本扩展(插件)不收集关于您个人的任何信息,本扩展第一项功能为转换对应跨境电商平台的货币为人民币,直接操作页面相关DOM, 相关汇率数据直接来自于我们的服务器,没必要使用您的个人信息。第二项功能为自由货币转换，同样没必要使用个人信息。
                关于第一项功能:在跨境电商平台显示人民币价格,是直接操作页面DOM,不用收集用户的cookie,请完全放心。</p><br>
            <p>本扩展(插件)在
              <a href="https://chrome.google.com/webstore/detail/%E6%B1%87%E7%8E%87%E8%BD%AC%E6%8D%A2/bcpgdpedphodjcjlminjbdeejccjbimp/reviews" target="_blank"><span class="tag is-danger">Google Chrome应用商店</span></a>,
              <a href="https://microsoftedge.microsoft.com/addons/detail/%E6%B1%87%E7%8E%87%E8%BD%AC%E6%8D%A2/jaippeddpgbjcdnmepklnjmcakmgodah" target="_blank"><span class="tag is-info">微软Edge应用商店</span></a>,
              <span class="tag is-warning">360浏览器应用商店</span>
              上架。经过三方审核,请放心使用！</p>
            </div>
        </div>
      <div class="container mt-3 pt-4  has-text-centered box">
        <div class="tag  is-size-6 mb-2">本次更新重点</div><br>
        <div class="tag is-danger is-size-6 mb-2">增加yandex的页面价格转换</div><br>
        <div class="tag is-info is-size-6 mb-5">汇率优化:以下主流货币的汇率来自<span class="tag is-warning is-medium">中国银行实时数据</span>
        </div>
        <div class="mid1">
          <div class="mid1-item">USD - 美元</div>
          <div class="mid1-item">EUR - 欧元</div>
          <div class="mid1-item">GBP - 英镑</div>
          <div class="mid1-item">JPY - 日元</div>
          <div class="mid1-item">AUD - 澳大利亚元</div>
          <div class="mid1-item">KRW - 韩国元</div>
          <div class="mid1-item">CAD - 加拿大元</div>
          <div class="mid1-item">NZD - 新西兰元</div>
          <div class="mid1-item">MOP - 澳门币</div>
          <div class="mid1-item">HKD - 港币</div>
          <div class="mid1-item">TWD - 台币</div>
          <div class="mid1-item">AED - 阿联酋迪拉姆</div>
          <div class="mid1-item">SAR - 沙特阿拉伯里亚尔</div>
          <div class="mid1-item">BRL - 巴西雷亚尔</div>
          <div class="mid1-item">CHF - 瑞士法郎</div>
          <div class="mid1-item">DKK - 丹麦克朗</div>
          <div class="mid1-item">SEK - 瑞典克朗</div>
          <div class="mid1-item">NOK - 挪威克朗</div>
          <div class="mid1-item">IDR - 印度尼西亚盾</div>
          <div class="mid1-item">MYR - 马来西亚林吉特</div>
          <div class="mid1-item">PHP - 菲律宾比索</div>
          <div class="mid1-item">SGD - 新加坡元</div>
          <div class="mid1-item">THB - 泰铢</div>
          <div class="mid1-item">INR - 印度卢比</div>
          <div class="mid1-item">RUB - 俄罗斯卢布</div>
          <div class="mid1-item">TRY - 土耳其里拉</div>
          <div class="mid1-item">ZAR - 南非兰特</div>

        </div>
      </div>
</template>

<script setup>
import { ref,onMounted,onBeforeMount } from 'vue'

//chrome-extension://knbpfapfapfggdfnljbfpmlellpkhakn/popup.html#/feed_back
const ext_url = ref('')
const version =chrome.runtime.getManifest().version
onBeforeMount(async() => {
    ext_url.value = chrome.runtime.getURL('popup.html#/feed_back')
})
</script>

<style scoped>
  .mid1{
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(170px, 1fr)); /* 创建响应式网格列 */
    grid-auto-rows: minmax(30px, auto);
    grid-gap: 5px;
    align-items: center;
    justify-items: start;
    font-size: 14px;
  }
  .mid1-item{
    border: 1px solid rgb(170, 14, 168);
    border-radius: 18px;
    padding: 6px;
  }
</style>