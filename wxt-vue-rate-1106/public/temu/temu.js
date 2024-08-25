console.log("TEMU价格转换启动..")

let cu_code = '' //货币代码
let rate = 0   //汇率
let temu_nodes = []
let color = 'green'
//====统一监听body的改变，触发总回调
let callback = function (records){
    //console.log('总回调')
    all();
}
//从cookie中获取货币代码
const get_cu_code = async()=>{
    let a = await cookieStore.get('currency')
    if(a&&a.value){
        return a.value
    }else{
        return ''
    }
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
    let cu = await get_cu_code()
    if(cu){
        cu_code = cu
    }else{
        return
    }
    chrome.storage.local.get(['my_rate','temu_nodes_new','color'],function(r){
    
            temu_nodes = r.temu_nodes_new
            rate = r.my_rate[`rate_${cu_code}`]
            if(r.color){
                color = r.color
            }
            find_node_list(temu_nodes)
      
    })
}

//找出元素
const find_node_list = function(nodes){
    for(let node of nodes){
        let list = document.querySelectorAll(node[0])
        if(list[0]){
           
            if(node[1] === 1){
                changePriceOfTheOneNode(list,node[2])
            }
            if(node[1] === 2){
                changePriceOfTheOneNode_main(list,node[2])
            }
        }
    }
}

//处理每个元素
const changePriceOfTheOneNode=function(list,styles){
    list.forEach(e=>{
        let c = e.innerText.trim()
        if (!e.nextElementSibling || (e.nextElementSibling && e.nextElementSibling.tagName != "SUB" && e.nextElementSibling.innerText!="%")  ) {
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
            //   b.style.lineHeight='normal'
              b.innerText = "¥" + rmb;
              e.after(b);
            }
          }
    })
}

const changePriceOfTheOneNode_main=function(list,styles){
    list.forEach(e=>{
        let c = e.innerText.trim()
        if (!e.nextElementSibling || (e.nextElementSibling && e.nextElementSibling.tagName != "SUB" && e.nextElementSibling.innerText!="%")  ) {
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
            //   b.style.lineHeight='normal'
              b.innerText = "¥" + rmb;
              e.after(b);
            }
          }else if(e.nextElementSibling && e.nextElementSibling.tagName === "SUB" && e.nextElementSibling.innerText.includes('¥')){
            if(c.includes('%')){
                return
            }
            let rmb = getRmb(c);
            if (rmb) {
              e.nextElementSibling.innerText = "¥" + rmb;
            }
          }
    })
}
//计算人民币价格
const getRmb = function(s){
    if(cu_code==='EUR'){
        let r1 = /^[^0-9]*/i
        let r2 = /[^0-9]*$/i
        let r3 = /,/i
        let s1 = s.replace(r1,'').replace(r2,'').replace(r3,'.')
        let b = parseFloat(s1)
        let rmb = b/rate
        if(rmb){
            return rmb.toFixed(1)
        }else{
            return 0
        }
    }
    if(cu_code === 'JPY'){
        let r1 = /^[^0-9]*/i
        let r2 = /[^0-9]*$/i
        let r3 = /,/i
        let s1 = s.replace(r1,'').replace(r2,'').replace(r3,'')
        let b = parseFloat(s1)
        let rmb = b/rate
        if(rmb){
            return rmb.toFixed(1)
        }else{
            return 0
        }
    }

    let r1 = /^[^0-9]*/i
    let r2 = /[^0-9]*$/i
    let s1 = s.replace(r1,'').replace(r2,'')
    let b = parseFloat(s1)
    let rmb = b/rate
    if(rmb){
        return rmb.toFixed(1)
    }else{
        return 0
    }
}