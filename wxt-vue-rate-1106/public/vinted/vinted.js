console.log('%cVinted价格转换插件启动','color:red;font-size:14px;');

//按url确定货币
let URL = document.URL;
//console.log(URL);
let currency = '';

if(URL.includes('vinted.pl/') ){
    currency = 'PLN';
}else if(   URL.includes('vinted.be/') || 
            URL.includes('vinted.at/') || 
            URL.includes('vinted.de/') ||
            URL.includes('vinted.es/') ||
            URL.includes('vinted.fi/') ||
            URL.includes('vinted.fr/') ||
            URL.includes('vinted.it/') ||
            URL.includes('vinted.lt/') ||
            URL.includes('vinted.lu/') ||
            URL.includes('vinted.nl/') ||
            URL.includes('vinted.sk/') ||
            URL.includes('vinted.pt/')
        ){
    currency = 'EUR';
}else if(URL.includes('vinted.cz/')){
    currency = 'CZK';
}else if( URL.includes('vinted.dk/') ){
    currency = 'DKK';
}else if(URL.includes('vinted.hu/')){
    currency = 'HUF';
}else if( URL.includes('vinted.ro/') ){
    currency = 'RON';
}else if( URL.includes('vinted.se/')  ){
    currency = 'SEK';
}else if( URL.includes('vinted.co.uk/') ){
    currency = 'GBP';
}else if(  URL.includes('vinted.com/')){
    currency = 'USD';
}


let rate = 0
let color = 'blueviolet'
let vinted_nodes = [
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

function  all() {
  
    if(currency === ''){
        return
    }
    //console.log('currency:',currency)
    chrome.storage.local.get(["my_rate",'vinted_nodes','color'], function (result) {
        rate = result.my_rate[`rate_${currency}`];
        vinted_nodes = result.vinted_nodes
        if(result.color){
            color = result.color
        }
        find_node(vinted_nodes);

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
            if(item[1] === 2){
                // special_price_do(item[0])
            }
        }
    }
}

function insetHtml(nodes){
    nodes.forEach(e => {
        let s = e.innerText
        if(!s.includes('¥')){
            let rmb = getRmb(s)
            if(rmb){
                let element_my = document.createElement('sub')
                element_my.style.color = color 
                // element_my.style.fontSize = '14px'
                element_my.style.fontWeight =  '500'
                // element_my.style["-webkit-text-fill-color"] = color
                element_my.setAttribute('translate', 'no');
                element_my.innerText = '¥' + rmb
                e.appendChild(element_my)
               
            }
        }
    })
}

function getRmb(s){
    let b = extractAndConvertCurrencyValue(s)
    let b1 = parseFloat(b)
    let rmb = b1/rate
    if(rmb){
        return rmb.toFixed(2)
    }else{
        return ''
    }
}

function extractAndConvertCurrencyValue(currencyString) {
    // 步骤1: 去除除数字、逗号、点号外的所有字符
    let numericValue = currencyString.replace(/[^\d,\.]/g, '');
  
    // 步骤2: 去除字符串前后的非数字字符
    numericValue = numericValue.replace(/^\D+|\D+$/g, '');
  
    // 步骤3: 将逗号换为点号
    numericValue = numericValue.replace(/,/g, '.');
  
    return numericValue;
  }
  