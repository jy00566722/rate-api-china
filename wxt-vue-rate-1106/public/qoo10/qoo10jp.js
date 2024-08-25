console.log('这里是趣天日本站价格转换插件,启动,GO!!!')

let country = 'JPY'

//====统一监听body的改变，触发总回调
let callback = function (records){
    //console.log('总回调')
    all();
}

let throttle_callback = _.throttle(callback,3500);
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

let rate = 0
let qoo10jp_nodes = []
let color = 'green'

//总回调
const all=function(){
    //console.log('总回调启动');
    chrome.storage.local.get(["my_rate",'qoo10jp_nodes','color'],function(result){
        rate = result.my_rate[`rate_${country}`];
        qoo10jp_nodes = result.qoo10jp_nodes
        if(result.color){
            color = result.color
        }
        foreach_nodes(qoo10jp_nodes)
    })
}

//取出元素数组的元素处理
const foreach_nodes = function(node_all){
    for(let node of node_all){
        let a = document.querySelectorAll(node)
        if(a[0]){
            changePriceOfTheOneNode(a)
        }
    }
}

const changePriceOfTheOneNode=function(nodes){
    nodes.forEach(e=>{
        let c=e.innerHTML.trim()
        if(!c.includes('円')){
            return
        }
/*         if((!e.nextElementSibling) || (e.nextElementSibling && e.nextElementSibling.tagName != 'SUB')){
            let rmb = getRmb(c)
            if(rmb){
                let b = document.createElement('sub')
                    b.style.color = 'green'
                    b.innerText = '¥'+rmb
                    e.after(b)
            }
        } */
        
        if(!c.includes('¥')){
            let rmb = getRmb(c)
            if(rmb){
                e.innerHTML = e.innerHTML  + `<sub translate="no" title="¥${rmb}" style="color:${color}"> ¥${rmb}</sub>`
            }
        }

    })
}

const getRmb = function(s){
    let a = s.replace(/[\,\.円]/g,'')
    let b = parseFloat(a)
    let rmb = b/rate
    if(rmb){
        return rmb.toFixed(1)
    }else{
        return 0
    }
}

