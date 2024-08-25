console.log('%c COUPANG 价格转换启动','color:green;font-size:16px;');

let cu_code = 'KRW' //货币代码
let rate = 0   //汇率
let coupang_nodes = [

]
let color = 'green'
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
//总回调
const all=async function(){

    chrome.storage.local.get(['my_rate','coupang','color'],function(r){
    
            coupang_nodes = r.coupang  //测试时直接读取配置好的节点
            rate = r.my_rate[`rate_${cu_code}`]
            if(r.color){
                color = r.color
            }
            find_node_list(coupang_nodes)
      
    })
}

//找出元素
const find_node_list = function(nodes){
    for(let node of nodes){
        let list = document.querySelectorAll(node[0])
        if(list[0]){
           
            if(node[1] === 1){
                changePriceOfTheOneNode(list,node[2]) //node[2] 指配置好的样式
            }
            if(node[1] === 2){
                changePriceOfTheOneNode_main(list,node[2]) //node[2] 指配置好的样式
            }
        }
    }
}

//处理每个元素
const changePriceOfTheOneNode=function(list,styles){
    list.forEach(e=>{
        let c = e.textContent.trim()
        // 获取最后一个子元素
        var lastChild = e.lastElementChild;
        if (!lastChild || (lastChild && lastChild.tagName.toLowerCase() !== 'sub') ) {
            if(c.includes('%')){
                return
            }
            let rmb = getRmb(c);
            if (rmb) {
              let b = document.createElement("sub");
              b.setAttribute('translate', 'no');
              b.style.color = color;
              b.style.fontSize="14px"
              if(Array.isArray(styles) && styles.length >= 1){
                styles.forEach(e=>{
                    b.style[e.styleName] = e.styleValue
                })
              }
              b.innerText = "¥" + rmb;
              e.appendChild(b);
            }
          }
    })
}

const changePriceOfTheOneNode_main=function(list,styles){
    list.forEach(e=>{
        let c = e.textContent.trim()
        var lastChild = e.lastElementChild;
        if (!lastChild || (lastChild && lastChild.tagName.toLowerCase() !== 'sub')  ) {
            if(c.includes('%')){
                return
            }
            let rmb = getRmb(c);
            if (rmb) {
              let b = document.createElement("sub");
              b.setAttribute('translate', 'no');
              b.style.color = color;
              b.style.fontSize="14px"
              if(Array.isArray(styles) && styles.length >= 1){
                styles.forEach(e=>{
                    b.style[e.styleName] = e.styleValue
                })
              }
              b.innerText = "¥" + rmb;
              e.after(b);
            }
          }else if(lastChild && lastChild.tagName.toLowerCase() === 'sub' && lastChild.innerText.includes('¥')){
            if(c.includes('%')){
                return
            }
            let rmb = getRmb(c);
            if (rmb) {
                lastChild.innerText = "¥" + rmb;
            }
          }
    })
}
//计算人民币价格
const  getRmb = function(s){
    let r1 = /^[^0-9]*/i
    let r2 = /[^0-9]*$/i
    let r3 = /,/g
    let s1 = s.replace(r1,'').replace(r2,'').replace(r3,'')
    let b = parseFloat(s1)
    let rmb = b/rate
    if(rmb){
        return rmb.toFixed(1)
    }else{
        return 0
    }
}