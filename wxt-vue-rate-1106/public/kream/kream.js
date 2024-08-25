console.log('%cKream价格转换插件启动','color:red;font-size:14px;');



let currency = 'KRW';

let rate = 0
let color = 'blueviolet'
let kream_nodes = [
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
    chrome.storage.local.get(["my_rate",'kream_nodes','color'], function (result) {
        rate = result.my_rate[`rate_${currency}`];
        kream_nodes = result.kream_nodes
        if(result.color){
            color = result.color
        }
        find_node(kream_nodes);

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
        let s = e.textContent
        if(!s.includes('¥')){
            let rmb = getRmb(s)
            if(rmb){
                // let element_my = document.createElement('sub')
                // element_my.style.color = color 
                // // element_my.style.fontSize = '14px'
                // element_my.style.fontWeight =  '500'
                // // element_my.style["-webkit-text-fill-color"] = color
                // element_my.setAttribute('translate', 'no');
                // element_my.innerText = '¥' + rmb
                // e.appendChild(element_my)

                e.innerHTML = e.innerHTML  + `<sub translate="no" title="¥${rmb}" style="color:${color}"> ¥${rmb}</sub>`
               
            }
        }
    })
}

function getRmb(s){
    let b = s.replace(/\D/g,'')
    let b1 = parseFloat(b)
    let rmb = b1/rate
    if(rmb){
        return rmb.toFixed(1)
    }else{
        return ''
    }
}


  