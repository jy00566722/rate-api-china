console.log('%cozon价格转换插件启动','color:red;font-size:14px;');

let currency = ''
let rate = 0
let color = 'blueviolet'
let ozon_nodes = [

]

//====统一监听body的改变，触发总回调
let callback = function (records) {
    all()
};
let throttle_callback = _.throttle(callback, 2000);
let mo = new MutationObserver(throttle_callback);
let option = {
    'childList': true,
    'subtree': true
};
let fs_node = document.getElementsByTagName("body")[0];
try {
    mo.observe(fs_node, option);
} catch (e) {
    console.log('监听器启动失败body."');
}

async function  all() {
    chrome.storage.local.get(["my_rate",'ozon_nodes','ozon_nodeStr','color'], function (result) {
        currency = getCurrency(result.ozon_nodeStr)
        if(!currency){
            console.log('%cozon获取货币代码失败','color:red;font-size:14px;');
            console.log('%c 采用默认的RUB','color:green;font-size:14px;'); 
            //如果获取失败,默认使用RUB
            currency = 'RUB'
            //return
        }else{
            // console.log('ozon获取货币代码:',currency);
        }
        if(currency==='CNY'){
            console.log('%cozon获取货币代码为CNY,不转换','color:red;font-size:14px;');
            return
        }
        let ozon_nodes = result.ozon_nodes
        rate = result.my_rate[`rate_${currency}`];
        ozon_nodes = result.ozon_nodes
        if(result.color){
            color = result.color
        }
        find_node(ozon_nodes);

    })
}

//找出元素
function find_node(node_all) {
    for (let item of node_all) {
        let a = document.querySelectorAll(item[0])
        if (a[0]) {
            if(item[1] === 1){
                insetHtml(a); 
            }
        }
    }
}

function insetHtml(nodes){
    nodes.forEach(e => {
        let s = e.innerText
        if(currency==='JPY'){
            if(/^[^¥]*¥[^¥]*$/.test(s)){
                let rmb = getRmb(s)
                if(rmb){
                    let element_my = document.createElement('sub')
                    element_my.style.color = color 
                    // element_my.style.fontSize = '14px'
                    element_my.style.fontWeight =  '500'
                    //position: inherit;
                    // element_my.style.position = "inherit"
                    element_my.style["-webkit-text-fill-color"] = color
                    element_my.setAttribute('translate', 'no');
                    //如果是俄罗斯站点，字体设置为细体
                    // if(currency === 'RUB'){
                    //     //字体颜色设置为紫色
                    //     element_my.style.color = color
                    //     element_my.style.fontSize = '14px'
                    //     element_my.style.fontWeight =  '500'
                    // }
                    element_my.innerText = '¥' + rmb
                    e.appendChild(element_my)
                   
                }
            }
            return
        }
        if(!s.includes('¥')){
            let rmb = getRmb(s)
            if(rmb){
                let element_my = document.createElement('sub')
                element_my.style.color = color 
                // element_my.style.fontSize = '14px'
                element_my.style.fontWeight =  '500'
                //position: inherit;
                // element_my.style.position = "inherit"
                element_my.style["-webkit-text-fill-color"] = color
                element_my.setAttribute('translate', 'no');
                //如果是俄罗斯站点，字体设置为细体
                // if(currency === 'RUB'){
                //     //字体颜色设置为紫色
                //     element_my.style.color = color
                //     element_my.style.fontSize = '14px'
                //     element_my.style.fontWeight =  '500'
                // }
                element_my.innerText = '¥' + rmb
                e.appendChild(element_my)
               
            }
        }
    })
}

//计算人民币价格
function getRmb(s){
    let r1 = /[^0-9,\.]/g
    let s1 = s.replace(r1,'')
    let s2 = s1.replace(',','.')
    let b = parseFloat(s2)
    let rmb = b/rate
    if(rmb){
        return rmb.toFixed(1)
    }else{
        return ''
    }
}

//获取当前页面的货币种类
const getCurrency = function (nodeStr){
    let element = document.querySelector(nodeStr);
    let Currency = element ? element.textContent : null;
    if(Currency==null){
        return null
    }else if(/^[A-Z]{3}$/.test(Currency)){
        return Currency
    }else{
        return null
    }
}