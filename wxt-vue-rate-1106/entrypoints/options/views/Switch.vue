<template>
    <div class="switch">
        <el-table :data="sites" style="width: 100%">
            <el-table-column type="index" label="序号" width="80"></el-table-column>
            <el-table-column prop="siteName" label="站点名称" width="120"></el-table-column>
            <!-- <el-table-column prop="site" label="站点" width="80"></el-table-column> -->
            <!-- <el-table-column prop="url" label="连接" width="180"></el-table-column> -->
            <el-table-column prop="status" label="开启状态" width="120">
                <template #default="scope">
                    <span class="tag is-success" v-if="scope.row.status == '已开启'">{{scope.row.status}}</span>
                    <span class="tag" v-else>{{scope.row.status}}</span>
                </template>
            </el-table-column>
            <el-table-column label="页面价转转换" width="180">
                <template #header>
                    <!-- <el-button size="small" type="info" disabled>手动开启</el-button> -->
                    <el-button size="small" type="danger" @click="delAll">全部关闭</el-button>
                  </template>
                <template #default="scope">
                    <el-button size="small"  type="primary"  @click="handRegister(scope.row,scope.$index)">开启</el-button>
                    <el-button size="small"  type="danger"  @click="handleDelete(scope.row.site)">关闭</el-button>
                </template>
            </el-table-column>

            <el-table-column label="跳转到" width="360">
                <template #default="scope">
                    <el-button size="small"  type="info" link v-for="item in scope.row.goUrl" :key="item.text" @click="openUrl(item.url)">{{item.text}}</el-button>
                </template>
            </el-table-column>
        </el-table>
    </div>
</template>

<script setup>
    import { ref,onMounted,onBeforeMount,computed } from 'vue'
    import { ElMessage,ElMessageBox } from 'element-plus'

    const sites = ref([
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

            ])  

    //检查已经注册的脚本,返回给siteList，显示状态
    async function chechStatus() {
            console.log('检查状态...')
            //先判断浏览器版本
            if (getChromeVersion()){
                    console.log("浏览器版本符合要求")
                }else{
                    ElMessage.error("请将浏览器版本升级到102及以上,再使用!")
                return
            }
        try {
            const scripts = await chrome.scripting.getRegisteredContentScripts();
            const scriptIds = scripts.map((script) => script.id); //已经注册的脚本id 列表
            sites.value.forEach(item => {
                if (scriptIds.includes(item.site)) {
                    item.status = "已开启"
                } else {
                    item.status = "未开启"
                }
            })
        } catch (error) {
            console.log(error);
            ElMessage.error("出现错误,错误代码005,可以反馈给开发者。或者:请将浏览器版本升级到102及以上!");
        }
    }

    //返回浏览器Chrome版本是否大于等于102
    function getChromeVersion(){
        let  ua = navigator.userAgent;
        let chrome = /Chrome\/([\d]+)./i.exec(ua);
        let version = parseInt(chrome[1])
        return version >= 102
    }

    //注册脚本
    async function handRegister(row,$index) {
        let site = sites.value[$index]
        let id = site.site
        let permissions = Array.from(new Set(site.permissions))
        let origins = Array.from(new Set(site.origins))
        let js = Array.from(new Set(site.js))
        let isReg = await isRegister(id)
        //弹出提示，提示用户要在接下来的弹出窗口中选择允许
        ElMessageBox({
            title: '请给予权限',
            message: '要开启页面价格转换，请在接下来的弹出窗口中选择<span class="tag is-danger">允许</span>.插件不会收集用户信息,请放心使用!',
            type: 'warning',
            showCancelButton: false,
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            dangerouslyUseHTMLString:true,
        }).then(async() => {
            //判断是否已经注册
                        if (isReg) {
                            console.log('清理之前的脚本')
                            await handleDelete(id,'noShow')
                        }
                        //点击确定后，请求权限
                        chrome.permissions.request(
                            {
                                permissions,
                                origins,
                            },
                            async (granted) => {
                                if (granted) {
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
                                            (e) => {
                                            console.log("注册成功");
                                            chechStatus()
                                            saveStorage()
                                            ElMessage.success("开启成功");
                                        }
                                    );
                                    }catch(e){
                                        console.log('发生错误01:',e)
                                    }
                                
                                } else {
                                    console.log("权限请求失败");
                                    ElMessage.error("请求权限失败,请插件给予权限,或者联系开发者!");
                                }
                            }
                );
        }).catch(() => {

            //点击取消后，提示用户
            // ElMessage.error("开启失败,错误代码006,可以反馈给开发者.或者:请将浏览器版本升级到102及以上!");
        });
}
    
    //删除所有脚本
    async function delAll() {
        try {
            const scripts = await chrome.scripting.getRegisteredContentScripts(); //unregisterContentScripts
            const scriptIds = scripts.map((script) => script.id); //已经注册的脚本id 列表
            await chrome.scripting.unregisterContentScripts({ids:scriptIds});
            chechStatus();
            saveStorage()
            ElMessage.success("关闭成功");
        } catch (error) {
            console.log(error);
            chechStatus()
            saveStorage()
            ElMessage.error("出现错误,错误代码010,可以反馈给开发者。");
        }
    }
    //删除单个脚本
    async function handleDelete(id,noShow) {
        // let site = sites.value[$index]
        // let id = site.site
        let isReg = await isRegister(id)
        if (!isReg) {
            ElMessage.error("脚本未注册,无需关闭");
            return
        }
        try {
            await chrome.scripting.unregisterContentScripts({ids:[id]});
            chechStatus();
            saveStorage()
            if (noShow == 'noShow') return
            ElMessage.success("关闭成功");
        } catch (error) {
            console.log(error);
            chechStatus()
            saveStorage()
            ElMessage.error("出现错误,错误代码011,可以反馈给开发者。");
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
            ElMessage.error("出现错误,错误代码0012,可以反馈给开发者。");
        }
    }
    //打开url
    function openUrl(url) {
        chrome.tabs.create({ url: url });
    }
    //把注册过的站点代号，保存到storage中
    async function saveStorage() {
        try {
            const scripts = await chrome.scripting.getRegisteredContentScripts();
            const scriptIds = scripts.map((script) => script.id); //已经注册的脚本id 列表
            if (scriptIds.length > 0) {
                chrome.storage.local.set({ registedSites: scriptIds }, function () {
                    console.log('保存成功');
                });
            } else {
                await chrome.storage.local.remove('registedSites')
            }
        } catch (error) {
            console.log(error);
            ElMessage.error("出现错误,错误代码0012,可以反馈给开发者。");
        }
    }
    onMounted(async() => {
                chechStatus()
            })
</script>

<style>

</style>