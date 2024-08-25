console.log('%c walmart价格转换插件启动','color:red;font-size:14px;');

let cu_code = 'USD' //货币代码
let rate = 0   //汇率
let walmart_nodes = [

]
//紫色
let color = 'blueviolet' 
//====统一监听body的改变，触发总回调
let callback = function (records){
    all();
}
let throttle_callback = _.throttle(callback,2800);
let mo = new MutationObserver(throttle_callback);
let option = {
  'childList': true,
  'subtree': true
};
let fs_node = document.getElementsByTagName("body")[0];
try{   
    mo.observe(fs_node, option);
}catch(e){
    console.log('监听器启动失败body."');
}

const all=async function(){

    chrome.storage.local.get(['my_rate','walmart','color'],function(r){

        if(r.walmart){
            walmart_nodes = r.walmart
        }
        if(r.my_rate){
            rate = r.my_rate[`rate_${cu_code}`]
        }
        if(r.color){
            color = r.color
        }
        if(walmart_nodes.length>0){

            find_node_list(walmart_nodes)
        }
    })
}

const find_node_list=function(nodes){
    // console.log('%c 查找节点','color:red;font-size:14px;');
    for(let i=0;i<nodes.length;i++){
        let node = nodes[i]
        let node_list = document.querySelectorAll(node[0])
        if(node_list.length>0){
            node_list.forEach(function(n){
                // console.log(n)
                if(node[1]==1){
                    changePriceOfTheOneNode([n], node[2]) // 修改这里，传递单个节点的数组
                }
            })
        }
    }
}

//处理每个元素
const changePriceOfTheOneNode = function(list, styles) {
    list.forEach(e => {
        // 检查是否已经转换过或者是百分比
        if (!e.nextElementSibling || (e.nextElementSibling && e.nextElementSibling.tagName !== "SUB" && e.nextElementSibling.innerText !== "%")) {
            let usdPrice = extractComplexPrice(e);
            if (usdPrice) {
                let rmbPrice = convertToRMB(usdPrice);
                addRMBElement(e, rmbPrice, styles);
            }
        }
    })
}

const extractComplexPrice = function(element) {
    // 如果只有一个文本节点
    if (element.childNodes.length === 1 && element.childNodes[0].nodeType === Node.TEXT_NODE) {
        let text = element.textContent.trim();
        // 匹配价格，包括可能的区间价格和千分位
        let match = text.match(/\$([\d,]+(?:\.\d+)?)/);
        return match ? parseFloat(match[1].replace(/,/g, '')) : null;
    }

    // 如果有多个子元素（假设是4个span）
    let spans = element.querySelectorAll('span');
    if (spans.length >= 3) {
        // 忽略第一个span（可能是广告语或空）
        let integerPart = spans[2].textContent.trim().replace(/[^\d]/g, '');
        let fractionalPart = spans[3] ? spans[3].textContent.trim().replace(/[^\d]/g, '') : '00';
        
        return parseFloat(integerPart + '.' + fractionalPart);
    }

    // 如果都不匹配，返回null
    return null;
}

const convertToRMB = function(usdPrice) {
    return (usdPrice / rate).toFixed(1);
}

const addRMBElement = function(element, rmbPrice, styles) {
    let b = document.createElement("sub");
    b.setAttribute('translate', 'no');
    b.style.color = color;
    b.style.fontSize = "14px"
    if (Array.isArray(styles) && styles.length >= 1) {
        styles.forEach(e => {
            b.style[e.styleName] = e.styleValue
        })
    }
    b.innerText = "¥" + rmbPrice;
    element.after(b);
}