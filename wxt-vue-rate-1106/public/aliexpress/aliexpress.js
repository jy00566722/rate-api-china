console.log('%c速卖通价格转换插件启动','color:red;font-size:14px;');

const getRate = async()=>{
    let a = await cookieStore.get('aep_usuc_f')
    let b = a.value.split('&')
    let c = '' 
    let d = ''
    b.map(e=>{
        const [a1 ,a2] = e.split('=')
        if(a1 === 'c_tp'){
            c = a2
        }
        if(a1 === 'site'){
            d = a2
        }

    })
    return [c,d]
}

let currency = ''
let site = ''
let rate = 0
let color = 'green'
let aliexpress_nodes = [

]

//====统一监听body的改变，触发总回调
let callback = function (records) {
    //console.log('监听器启动.');
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
    //console.log('总回调启动.'); 
    if(currency === ''){
        [currency,site] = await getRate()
        console.log('info:',currency,site)
    }
    //console.log('currency:',currency)
    chrome.storage.local.get(["my_rate",'aliexpress_nodes','color'], function (result) {
        rate = result.my_rate[`rate_${currency}`];
        aliexpress_nodes = result.aliexpress_nodes
        if(result.color){
            color = result.color
        }
        find_node(aliexpress_nodes);

    })
}

//找出元素
function find_node(node_all) {
    for (let item of node_all) {
        let a = document.querySelectorAll(item[0])
        if (a[0]) {
            if(item[1] === 1){
                insetHtml1(a,item[2]); //item[2]是特殊的class,是主价格-会随规则的变化而变化
            }
            if(item[1] === 2){
                special_price_do(item[0])
            }
        }
    }
}

//按传过来的css选择器找出元素,并替换html。第一种模式是取出元素中的innerText,按innerText的规则计算出人民币价格。然后插入到元素的后面。
const insetHtml1 = function (nodes,main_price_class) {
    nodes.forEach(e => {
            let s = e.innerText
            if(s.includes('%')){return}
            if(s.includes('Order')){return}
            if(s.includes('-')){return}

            if(!s.includes('¥')){
                let rmb = getRmb(s)
                if(rmb ){
                    // e.innerHTML = e.innerHTML+ `<sub title="¥${rmb}" style="color:green"> ¥${rmb}</sub>`
                    let element_my = document.createElement('sub')
                    element_my.style.color = color
                    element_my.setAttribute('translate', 'no');
                    //如果是俄罗斯站点，字体设置为细体
                    if(currency === 'RUB'){
                        //字体颜色设置为紫色
                        element_my.style.color = color
                        element_my.style.fontSize = '14px'
                        element_my.style.fontWeight =  '500'
                    }

                    
                    element_my.innerText = '¥' + rmb
                    if(main_price_class){
                        //检查元素e的下一个兄弟元素是不是sub,如果不是就当做兄弟元素直接把element_my插入后面,如果是就改变sub的innerText
                        let next_element = e.nextElementSibling
                        if(next_element && next_element.tagName === 'SUB'){
                            next_element.innerText = '¥' + rmb
                        }else{
                            e.insertAdjacentElement('afterend', element_my);
                        }
                    }else{
                    e.appendChild(element_my)
                    }
                }
            }
    })
}



const getRmb = function (s) {
    //console.log('====')
    //console.log(s)
    let r1 = /^[^0-9]*/i
    let r2 = /[^0-9]*$/i
    let r3 = /,/g
    let r4 = /\./g
    let r5 = /\s/g
    let r6 = /^[0-9\,\s]{1,}$/
    let r7 = /руб\./
    let r8 = /&nbsp;/g
    let s1 = s.replace(r1,'').replace(r2,'')
    //这里对只取了一位小数点的价格进行处理,比如 ￥25.2，在后面加一个0,方便后面的正则处理
    if(s1.includes('.') && s1.slice(-2,-1) ==='.'){
        s1 = s1 + '0'
    }
    if(s1.includes(',') && s1.slice(-2,-1) ===','){
        s1 = s1 + '0'
    }
 

    if(s1.includes('-')){
        return 0
    }

    let s1_ = s1.slice(-3,-2)
    if(s1_ === '.' || s1_ === ','){
        let s2 = s1.slice(0,-3)+ '|' +s1.slice(-3)
        let s3 = s2.replace(r3,'').replace(r4,'').replace(r5,'').replace(r7,'').replace(r8,'').replace(/\|/g,'.')
        let b = parseFloat(s3)
        let rmb = b/rate
        if(rmb){
            return rmb.toFixed(2)
        }else{
            return 0
        }
    }

    //all Number  256 256,215,125
    if(r6.test(s1)){
        let s1_ = s1.replace(r3,'').replace(r5,'')
        let b = parseFloat(s1_)
        let rmb = b/rate
        if(rmb){
            return rmb.toFixed(2)
        }else{
            return 0
        }
    }
}

//document.querySelectorAll('.U9mS2 ._2FkhA')[0].setAttribute('style','overflow:visible;')
//这里是为了不让人民币价格被隐藏
const setElementStyle = function(list){
    list.map(el=>{
        let a = document.querySelectorAll(el[0])
        if(a[0]){
            a.forEach(e=>{
                e.setAttribute(el[1],el[2])
            })
        }
    })
}

const special_price_do = function(node_list){
    node_list.map(el=>{
        let a = document.querySelectorAll(el)
        if(a[0]){
            a.forEach(e=>{
                let b = e.innerText
                if(!b.includes('¥')){
                    //let c = b.replace(/[^0-9\,\.]/g,'')
                    //let d = c.replace(',','.')
                    let rmb = getRmb(b)
                    if(rmb){
                        let element_my = document.createElement('sub')
                        element_my.style.color = 'green'
                        element_my.setAttribute('translate', 'no');
                        element_my.innerText = '¥' + rmb
                        e.appendChild(element_my)
                    }
                }                
            })
        }
    })
}


//获取html元素的所有直接子元素的标签，如果标签中有一个不是span，直接回返false.使用for循环，而不是forEach,是因为forEach不能使用break提前结束循环
const get_element_child = function(e){
    let a = e.children
    if(a.length === 0){
        return false
    }
    for(let i = 0;i<a.length;i++){
        if(a[i].tagName != 'SPAN'){
            return false
        }
    }
    return true
}

