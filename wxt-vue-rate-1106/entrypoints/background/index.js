export default defineBackground(() => {
    console.log('%c%s', 'color: #0000ff; background: #ffffff', ' backgroud-server worker start...')
    //用红色字黄色底打印当前时间,按中国的时间规格
    console.log('%c%s', 'color: #ff0000; background: #ffff00', (new Date()).toLocaleString('zh-CN', { hour12: false }))

      const get_rate_v3=async (my_id,version,cache)=>{
        let url = 'https://rate.lizudi.top/v3/rate?my_id='+my_id+'&version='+version
        if(cache === 'noCache'){
          url = url + '&noCache=yes'
        }
        let a = await fetch(url)
        let b = await a.json()
        if(b&&b.code === 20000){
          let my_rate={}
          Object.keys(b.data).forEach(el=>{
            my_rate['rate_'+el] = b.data[el].rate
          })
          chrome.storage.local.set({my_rate,my_rate_zz:b.data},function(){
            console.log('保存汇率成功','my_rate,my_rate_zz')
          })
        }else{
          console.log('保存汇率失败')
        }
      }
      

      
      //向激活的窗口，发送消息，弹出主界面
      function sendMessageToContentScript(message, callback) {
        chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
          chrome.tabs.sendMessage(tabs[0].id, message, function (response) {
            if (callback) callback(response);
          });
        });
      }
      
      function generateUUID() {
        let d = new Date().getTime();
      /*   if (window.performance && typeof window.performance.now === "function") {
            d += performance.now()
        } */
        let uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
            let r = (d + Math.random() * 16) % 16 | 0;
            d = Math.floor(d / 16);
            return (c == 'x' ? r : (r & 0x3 | 0x8)).toString(16);
        });
        return uuid
      }
      //设置uuid,如果有就不设置了，没有才设置
      function setUUid(){
        chrome.storage.local.get(['my_id'],function(s){
          if(s.my_id){
            // console.log('存在my_id')
          }else{
            // console.log('不存在my_id')
            chrome.storage.local.set({my_id:generateUUID()},function(){console.log('设置my_id成功.')})
          }
        })
      }

      //监听安装，更新事件
      chrome.runtime.onInstalled.addListener(function(e){
        if(e.reason === 'install'){
                console.log('插件安装')
                //注册站点-之前已经请求好权限的
                register_site()
                setUUid() 
                //初始化货币国家
                let countries = ['CNY','USD','EUR','JPY','GBP','KRW','CAD','AUD','MXN','PHP','TWD','MYR','VND']
                chrome.storage.local.set({ "my_rate_zz_countries": countries }, function (s) {  
                  // console.log("初始化::货币列表在my_rate_zz_countries中...");
                });
                checkAlarmState("install!");
                start_all_function()
      
                if (chrome.runtime.openOptionsPage) {
                  chrome.runtime.openOptionsPage();
                } else {
                  window.open(chrome.runtime.getURL('options.html'));
                }
                
        }
        if(e.reason === 'update'){
                console.log('插件更新')
                //注册站点-之前已经请求好权限的
                register_site()
                setUUid()
                start_all_function()

                // 在这里，可以处理与更新相关的任务，比如清除旧的Alarms
                checkAlarmState("update!");
                if (chrome.runtime.openOptionsPage) {
                  // chrome.runtime.openOptionsPage();
                } else {
                  // window.open(chrome.runtime.getURL('options.html'));
                }
              
  
        }
          if(e.reason === 'chrome_update'){
                  console.log('chrome更新') 
                  setUUid()
                  //注册站点-之前已经请求好权限的
                  register_site()
                  checkAlarmState("chrome_update!");
                  start_all_function()
                  
          }
      
      })
      
      
      const get_all_nodes_and_data = async (cache)=>{
        // let url = 'https://rate.lizudi.top/v2/all_nodes_and_date'
        let url = 'http://localhost:8001/v2/all_nodes_and_date'
        if (cache === 'noCache'){
          url = url + '?noCache=yes'
        }
        const a = await fetch(url)
        const b = await a.json()
        if(b && b.code === 20000){
          chrome.storage.local.set(b.all_data,()=>{
            // console.log('是否成功的保存了所有数据?')
          })
          //处理请求频率的问题
          let periodInMinutes = parseInt( b.all_data.periodInMinutes)
          if(periodInMinutes>2&&periodInMinutes<121){
            checkAlarmState("fetchChange!",periodInMinutes)
          }
        }else{
          // console.log('======>>不正确了吗')
        }
      }
      //请求时加上my_id与version
      function go_get_rate(cache){
        chrome.storage.local.get(['my_id'],function(s){
          let version =chrome.runtime.getManifest().version
          if(!version){
            version = 'v0.0.0.0'
          }else{
            version = "v"+ version
          }
          if(s&&s.my_id){
            get_rate_v3(s.my_id,version,cache)
          }else{
            get_rate_v3('0000000000',version,cache)
          }
        })
      }
      
      const start_all_function = (cache)=>{
        console.log('%c%s startAllFunction启动..', 'color: #ff0000; background: #ffff00', (new Date()).toLocaleString('zh-CN', { hour12: false }))
        go_get_rate(cache);//初始化汇率

        get_all_nodes_and_data(cache)

      
        chrome.storage.local.set({'get_Data_Time':(new Date()).toString()},function(){
          // console.log('获取所有数据..')
        })
      }
      
    async function checkAlarmState(s,periodInMinutes) {     
          console.log('检查alarm状态:',s)
          if(periodInMinutes){//如果有传入时间
            try{
              let a = await chrome.alarms.get("my-alarm")
              if (a&&a.periodInMinutes === periodInMinutes){//如果alarm已经存在,不需要修改
                // console.log('alarm已经存在,不需要修改')
              }else{//alarm不存在,或是时间不同,需要修改
                // console.log('alarm不存在-或时间不同,需要修改')
                await chrome.alarms.clearAll();
                await chrome.alarms.create("my-alarm",{ periodInMinutes: periodInMinutes });
              }
            }catch(e){
              console.log(e)
              await chrome.alarms.create("my-alarm",{ periodInMinutes: 10 });
            }

          }else{ //如果没有传入时间
            try{
              await chrome.alarms.clearAll();
              let b = await chrome.storage.local.get("periodInMinutes")
              if(b&&b.periodInMinutes){
                await chrome.alarms.create("my-alarm",{ periodInMinutes: b.periodInMinutes });
              }else{ 
                await chrome.alarms.create("my-alarm",{ periodInMinutes: 10 });
            }
            }catch(e){
              console.log(e)
              await chrome.alarms.create("my-alarm",{ periodInMinutes: 10 });
            }
    }
  }

      
      checkAlarmState("top!");
      
      chrome.alarms.onAlarm.addListener((r) => {
        if(r.name === 'my-alarm'){
          start_all_function()
        }
      })
          
    //设置卸载反馈
    chrome.runtime.setUninstallURL(
        'https://rate.lizudi.top/v2/uninstall',()=>{
            console.log('设置反馈URL成功')
        }
    ) 

    //设置接收消息
    chrome.runtime.onMessage.addListener(
        function(request, sender, sendResponse) {
          if (request.do === "getStart" && request.cache === "no")
            sendResponse({status: "good"});
            start_all_function('noCache')
        }
      );



    //在启动/安装/更新时，获取已经注册的站点代号，重新注册
    const register_site = ()=>{
      const sites = [
        {
              site:"gmarket",
              permissions: ["scripting"],
              origins: ["http://*.gmarket.co.kr/*", "https://*.gmarket.co.kr/*"],
              js: ["js/Underscore.js", "gmarket/gmarket.js"],
              siteName:"韩国Gmarket",
              url:"https://www.gmarket.co.kr/",
              status:'',
              goUrl:[{
                  text:'https://www.gmarket.co.kr/',
                  url:'https://www.gmarket.co.kr/'
              }]
          },
          {
              site:"temu",
              permissions: ["scripting"],
              origins: ["https://*.temu.com/*"],
              js: ["js/Underscore.js", "temu/temu.js"],
              siteName:"Temu",
              url:"https://www.temu.com/",
              status:'',
              goUrl:[{
                  text:'https://www.temu.com/',
                  url:'https://www.temu.com/'
              }]

          },
       {
              site:"amazon",
              permissions: ["scripting"],
              origins: ["https://*.amazon.com/*",//美国
                      "https://*.amazon.co.jp/*",//日本
                      "https://*.amazon.co.uk/*",//英国
                      "https://*.amazon.de/*",//德国
                      "https://*.amazon.com.br/*",//巴西
                      "https://*.amazon.com.mx/*",//墨西哥
                      "https://*.amazon.com.au/*",//澳大利亚
                      "https://*.amazon.fr/*",//法国
                      "https://*.amazon.ca/*",//加拿大
                      "https://*.amazon.es/*",//西班牙
                      "https://*.amazon.it/*",//意大利
                      "https://*.amazon.in/*",  //印度
                      "https://*.amazon.com.be/*",//比利时
                      "https://*.amazon.sg/*",//新加坡
                      "https://*.amazon.nl/*",//荷兰
                      "https://*.amazon.pl/*",//波兰
                      "https://*.amazon.se/*",//瑞典
                      "https://*.amazon.com.tr/*",//土耳其
                      "https://*.amazon.eg/*",//埃及
                      "https://*.amazon.sa/*",//沙特
                      "https://*.amazon.ae/*",//阿联酋

              ],
              js: ["js/Underscore.js", "amazon/amazon-all.js"],
              siteName:"亚马逊",
              url:"https://www.amazon.com/",
              status:'',
              goUrl:[
                  {
                      text:'美国站',
                      url:'https://www.amazon.com/'
                  },
                  {
                      text:'日本站',
                      url:'https://www.amazon.co.jp/'
                  },
                  {
                      text:'英国站',
                      url:'https://www.amazon.co.uk/'
                  },
                  {
                      text:'德国站',
                      url:'https://www.amazon.de/'
                  },
                  {
                      text:'巴西站',
                      url:'https://www.amazon.com.br/'
                  },
                  {
                      text:'墨西哥站',
                      url:'https://www.amazon.com.mx/'
                  },
                  {
                      text:'澳大利亚站',
                      url:'https://www.amazon.com.au/'
                  },
                  {
                      text:'法国站',
                      url:'https://www.amazon.fr/'
                  },
                  {
                      text:'加拿大站',
                      url:'https://www.amazon.ca/'
                  },
                  {
                      text:'西班牙站',
                      url:'https://www.amazon.es/'
                  },
                  {
                      text:'意大利站',
                      url:'https://www.amazon.it/'
                  },
                  {
                      text:'印度站',
                      url:'https://www.amazon.in/'
                  },
                  {
                      text:'比利时站',
                      url:'https://www.amazon.com.be/'
                  },
                  {
                      text:'新加坡站',
                      url:'https://www.amazon.sg/'
                  },
                  {
                      text:'荷兰站',
                      url:'https://www.amazon.nl/'
                  },
                  {
                      text:'波兰站',
                      url:'https://www.amazon.pl/'
                  },
                  {
                      text:'瑞典站',
                      url:'https://www.amazon.se/'
                  },
                  {
                      text:'土耳其站',
                      url:'https://www.amazon.com.tr/'
                  },
                  {
                      text:'埃及站',
                      url:'https://www.amazon.eg/'
                  },
                  {
                      text:'沙特站',
                      url:'https://www.amazon.sa/'
                  },
                 {
                      text:'阿联酋站',
                      url:'https://www.amazon.ae/'
                 }
              ]
          },
         {
              site:"shopee",
              permissions: ["scripting"],
              origins: [
                  "https://*.shopee.com.my/*",//马来西亚
                  "https://*.shopee.ph/*",//菲律宾
                  "https://*.shopee.sg/*",//新加坡
                  "https://*.shopee.co.id/*",//印度尼西亚
                  "https://*.shopee.tw/*",//台湾
                  "https://*.shopee.co.th/*",//泰国
                  "https://*.shopee.vn/*",//  越南
                  "https://*.shopee.com.br/*",//巴西
                  "https://*.shopee.com.mx/*",//墨西哥
                  "https://*.shopee.com.co/*",//哥伦比亚
                  "https://*.shopee.cl/*",//智利
                  "https://*.co.xiapibuy.com/*",
                  "https://*.cl.xiapibuy.com/*",
                  "https://*.th.xiapibuy.com/*",
                  "https://*.vn.xiapibuy.com/*",
                  "https://*.ph.xiapibuy.com/*",
                  "https://*.my.xiapibuy.com/*",
                  "https://*.id.xiapibuy.com/*",
                  "https://*.sg.xiapibuy.com/*",
                  "https://*.xiapi.xiapibuy.cc/*",
                  "https://*.xiapi.xiapibuy.com/*",
                  "https://*.br.xiapibuy.com/*",
                  "https://*.mx.xiapibuy.com/*",
              ],
              js: ["js/Underscore.js", "shopee/my.js"],
              siteName:"Shopee",
              url:"https://shopee.vn/",
              status:'',
              goUrl:[
                  {
                      text:'马来西亚站',
                      url:'https://shopee.com.my/'
                  },
                  {
                      text:'菲律宾站',
                      url:'https://shopee.ph/'
                  },
                  {
                      text:'新加坡站',
                      url:'https://shopee.sg/'
                  },
                  {
                      text:'印度尼西亚站',
                      url:'https://shopee.co.id/'
                  },
                  {
                      text:'台湾站',
                      url:'https://shopee.tw/'
                  },
                  {
                      text:'泰国站',
                      url:'https://shopee.co.th/'
                  },
                  {
                      text:'越南站',
                      url:'https://shopee.vn/'
                  },
                  {
                      text:'巴西站',
                      url:'https://shopee.com.br/'
                  },
                  {
                      text:'墨西哥站',
                      url:'https://shopee.com.mx/'
                  },
                  {
                      text:'哥伦比亚站',
                      url:'https://shopee.com.co/'
                  },
                  {
                      text:'智利站',
                      url:'https://shopee.cl/'
                  },
              ]
              },
              {
                  site:"lazada",
                  permissions: ["scripting"],
                  origins: [
                      "https://*.lazada.com.my/*",//马来西亚
                      "https://*.lazada.com.ph/*",//菲律宾
                      "https://*.lazada.sg/*",//新加坡
                      "https://*.lazada.co.id/*",//印度尼西亚
                      "https://*.lazada.vn/*",//越南
                      "https://*.lazada.co.th/*"//泰国
                  ],
                  js: ["js/Underscore.js", "lazada/lazada_all.js"],
                  siteName:"Lazada",
                  url:"https://www.lazada.com.my/",
                  status:'',
                  goUrl:[{
                      text:'菲律宾站',
                      url:'https://www.lazada.com.ph/'
                  },
                  {
                      text:'马来西亚站',
                      url:'https://www.lazada.com.my/'
                  },
                  {
                      text:'新加坡站',
                      url:'https://www.lazada.sg/'
                  },
                  {
                      text:'印度尼西亚站',
                      url:'https://www.lazada.co.id/'
                  },
                  {
                      text:'越南站',
                      url:'https://www.lazada.vn/'
                  },
                  {
                      text:'泰国站',
                      url:'https://www.lazada.co.th/'
                  }
              ]
              },
               {
                  site:"qoo10",
                  permissions: ["scripting"],
                  origins: [
                      "https://www.qoo10.jp/*"
                  ],
                  js: ["js/Underscore.js", "qoo10/qoo10jp.js"],
                  siteName:"趣天Qoo10",
                  url:"https://www.qoo10.jp/",
                  status:'',
                  goUrl:[{
                      text:'https://www.qoo10.jp/',
                      url:'https://www.qoo10.jp/'
                  }]
              },
              {
                  site:"wowmajp",
                  permissions: ["scripting"],
                  origins: [
                      "https://*.wowma.jp/*"
                  ],
                  js: ["js/Underscore.js", "wowmajp/wowmajp.js"],
                  siteName:"wowmajp",
                  url:"https://www.wowma.jp/",
                  status:'',
                  goUrl:[{
                      text:'https://www.wowma.jp/',
                      url:'https://www.wowma.jp/'
                  }]
              },
               {
                  site:"aliexpress",
                  permissions: ["scripting"],
                  origins: [
                      "https://*.aliexpress.com/*",
                      "https://*.aliexpress.ru/*",
                      "https://*.aliexpress.us/*",
                  ],
                  js: ["js/Underscore.js", "aliexpress/aliexpress.js"],
                  siteName:"aliExpress",
                  url:"https://www.aliexpress.com/",
                  status:'',
                  goUrl:[{
                      text:'https://www.aliexpress.com/',
                      url:'https://www.aliexpress.com/'
                  }]
              },
              {
                  site:"rakuten",
                  permissions: ["scripting"],
                  origins: [
                      "https://*.rakuten.co.jp/*",
                  ],
                  js: ["js/Underscore.js", "rakuten/rakuten.js"],
                  siteName:"乐天rakuten",
                  url:"https://www.rakuten.co.jp/",
                  status:'',
                  goUrl:[{
                      text:'https://www.rakuten.co.jp/',
                      url:'https://www.rakuten.co.jp/'
                  }]
              },
              {
                site:"ozon",
                permissions: ["scripting"],
                origins: [
                    "https://www.ozon.ru/*",
                ],
                js: ["js/Underscore.js", "ozon/ozon.js"],
                siteName:"ozon",
                url:"https://www.ozon.ru/",
                status:'',
                goUrl:[{
                    text:'https://www.ozon.ru/',
                    url:'https://www.ozon.ru/'
                }]
            },
            {
              site:"Kream",
              permissions: ["scripting"],
              origins: [
                  "https://*.kream.co.kr/*",
              ],
              js: ["js/Underscore.js", "kream/kream.js"],
              siteName:"Kream",
              url:"https://kream.co.kr/",
              status:'',
              goUrl:[{
                  text:'https://kream.co.kr/',
                  url:'https://kream.co.kr/'
              }]
          },
            {
              site:"mercadolibre",
              permissions: ["scripting"],
              origins: [
                  "https://*.mercadolibre.com.mx/*",
                  "https://*.mercadolibre.cl/*",
                  "https://*.mercadolivre.com.br/*",
                  "https://*.mercadolibre.com.co/*"
              ],
              js: ["js/Underscore.js", "mercadolibre/mercadolibre.js"],
              siteName:"美客多",
              url:"https://www.mercadolibre.com.mx/",
              status:'',
              goUrl:[{
                  text:'墨西哥站',
                  url:'https://www.mercadolibre.com.mx/'
              },
              {
                  text:'智利站',
                  url:'https://www.mercadolibre.cl/'
              },
              {
                  text:'巴西站',
                  url:'https://www.mercadolivre.com.br/'
              },
              {
                  text:'哥伦比亚',
                  url:'https://www.mercadolibre.com.co/'
              }
              ]
          },
            {
              site:"Vinted",
              permissions: ["scripting"],
              origins: [
                  "https://www.vinted.pl/*",
                  "https://www.vinted.at/*",
                  "https://www.vinted.be/*",
                  "https://www.vinted.cz/*",
                  "https://www.vinted.de/*",
                  "https://www.vinted.dk/*",
                  "https://www.vinted.es/*",
                  "https://www.vinted.fi/*",
                  "https://www.vinted.fr/*",
                  "https://www.vinted.hu/*",
                  "https://www.vinted.it/*",
                  "https://www.vinted.lt/*",
                  "https://www.vinted.lu/*",
                  "https://www.vinted.nl/*",
                  "https://www.vinted.ro/*",
                  "https://www.vinted.se/*",
                  "https://www.vinted.sk/*",
                  "https://www.vinted.co.uk/*",
                  "https://www.vinted.com/*",
                  "https://www.vinted.pt/*"
              ],
              js: ["js/Underscore.js", "vinted/vinted.js"],
              siteName:"Vinted",
              url:"https://www.vinted.pl/",
              status:'',
              goUrl:[{
                      text:'波兰站',
                      url:'https://www.vinted.pl/'
                      },
                      {
                      text:'比利时站',
                      url:'https://www.vinted.be/'
                      },
                      {
                      text:'奥地利',
                      url:'https://www.vinted.at/'
                      },
                      {
                      text:'捷克',
                      url:'https://www.vinted.cz/'
                      },
                      {
                      text:'德国',
                      url:'https://www.vinted.de/'
                      },
                      {
                      text:'丹表',
                      url:'https://www.vinted.dk/'
                      },
                      {
                      text:'西班牙',
                      url:'https://www.vinted.es/'
                      },
                      {
                      text:'芬兰',
                      url:'https://www.vinted.fi/'
                      },
                      {
                      text:'法国',
                      url:'https://www.vinted.fr/'
                      },
                      {
                      text:'匈牙利',
                      url:'https://www.vinted.hu/'
                      },
                      {
                      text:'意大利',
                      url:'https://www.vinted.it/'
                      },
                      {
                      text:'立陶宛',
                      url:'https://www.vinted.lt/'
                      },
                      {
                      text:'卢森堡',
                      url:'https://www.vinted.lu/'
                      },
                      {
                      text:'荷兰',
                      url:'https://www.vinted.nl/'
                      },
                      {
                      text:'罗马尼亚',
                      url:'https://www.vinted.ro/'
                      },
                      {
                      text:'葡萄牙',
                      url:'https://www.vinted.pt/'
                      },
                      {
                      text:'瑞典',
                      url:'https://www.vinted.se/'
                      },
                      {
                      text:'斯洛伐克',
                      url:'https://www.vinted.sk/'
                      },
                      {
                      text:'英国',
                      url:'https://www.vinted.co.uk/'
                      },
                      {
                      text:'美国',
                      url:'https://www.vinted.com/'
                      }
          ]
          },
          {
            site:"wildberries",
            permissions: ["scripting"],
            origins: [
                "https://www.wildberries.ru/*",
            ],
            js: ["js/Underscore.js", "wildberries/wildberries.js"],
            siteName:"wildberries",
            url:"https://www.wildberries.ru/",
            status:'',
            goUrl:[{
                text:'https://www.wildberries.ru/',
                url:'https://www.wildberries.ru/'
            }]
          },
          {
            site:"coupang",
            permissions: ["scripting"],
            origins: [
                "https://*.coupang.com/*",
            ],
            js: ["js/Underscore.js", "coupang/coupang.js"],
            siteName:"coupang",
            url:"https://www.coupang.com/",
            status:'',
            goUrl:[{
                text:'https://www.coupang.com/',
                url:'https://www.coupang.com/'
            }]
        },
        {
          site:"walmart",
          permissions: ["scripting"],
          origins: [
              "https://*.walmart.com/*",
          ],
          js: ["js/Underscore.js", "walmart/walmart.js"],
          siteName:"walmart",
          url:"https://www.walmart.com/",
          status:'',
          goUrl:[{
              text:'https://www.walmart.com/',
              url:'https://www.walmart.com/'
              }]
          }

          ]
      let sitesObj = {}
      sites.forEach(el=>{
        sitesObj[el.site] = el
      })
      chrome.storage.local.get(['registedSites'],function(s){
        let register_sites = s['registedSites']
        console.log('register_sites:',register_sites)
        if(register_sites){
          register_sites.forEach(el=>{
            if(sitesObj[el]){
              chrome.permissions.contains({
                origins: sitesObj[el].origins,
                permissions: ['scripting'],
              }, async function(result) {
                if (result) {
                  // console.log('有权限') 重新注册脚本
                  await handRegister(sitesObj[el])

                } else {
                  // console.log('没有权限,也不要请求权限')

                }
              });
            }
          })
        }
      })
    }
  
      //注册脚本
      async function handRegister(site) {
        let permissions = site.permissions
        let origins = site.origins
        let js = site.js
        let id = site.site
        //判断是否已经注册
        let isReg = await isRegister(id)
          if (isReg) {
              console.log('清理之前的脚本')
              await handleDelete(id)
          }
          //之前已经判断过，有权限，这里直接注册，不再请求权限
          console.log("权限请求成功"); 
          try{
              chrome.scripting.registerContentScripts(
              [
                  {
                      id,
                      matches: origins,
                      js,
                      runAt: "document_end",
                  },
              ],
                  () => {
                  console.log("后台注册成功:",id);
              }
          );
          }catch(e){
              console.log('发生错误01:',e)
          }

}
    //删除单个脚本
    async function handleDelete(id) {
      let isReg = await isRegister(id)
      if (!isReg) {
          return
      }
      try {
          await chrome.scripting.unregisterContentScripts({ids:[id]});
      } catch (error) {
          console.log(error);

      }
  }

      //判断脚本是否已经注册
      async function isRegister(id) {
        try {
            const scripts = await chrome.scripting.getRegisteredContentScripts();
            const scriptIds = scripts.map((script) => script.id); //已经注册的脚本id 列表
            if (scriptIds.includes(id)) {
                return true
            } else {
                return false
            }
        } catch (error) {
            console.log(error);
        }
    }
  });