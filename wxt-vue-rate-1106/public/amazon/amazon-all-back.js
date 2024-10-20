console.log('Amazon价格转换插件启动...');

//
let URL = document.URL
let origin_=document.location.origin
let country = '' //rate  
let color = 'green'
if(URL.includes('amazon.com/')){ //美国
    country = 'USD'
}else if(URL.includes('amazon.co.jp/')){ //
    country = 'JPY'
}else if(URL.includes('amazon.co.uk/')){ //英国
    country = 'GBP'
}else if(URL.includes('amazon.de/')|| URL.includes('amazon.fr/') || URL.includes('amazon.es/') || URL.includes('amazon.it/')){
    country = 'EUR'
}else if(URL.includes('amazon.com.br/')){//巴西
    country = 'BRL'
}else if(URL.includes('amazon.com.mx/')){//墨西哥
    country = 'MXN'
}else if(URL.includes('amazon.com.au/')){//澳大利亚
    country = 'AUD'
}else if(URL.includes('amazon.ca/')){// 加拿大
    country = 'CAD'
}else if(URL.includes('amazon.in/')){ //印度
    country = 'INR'
}else if(URL.includes('amazon.com.be/')){//比利时
    country = 'EUR'
}else if(URL.includes('amazon.sg/')){ //新加坡
    country = 'SGD'
}else if(URL.includes('amazon.nl/')){ //荷兰
    country = 'EUR'
}else if(URL.includes('amazon.pl/')){ //波兰
    country = 'PLN'
}else if(URL.includes('amazon.se/')){ //瑞典
    country = 'SEK'
}else if(URL.includes('amazon.com.tr/')){//土耳其
    country = 'TRY'
}else if(URL.includes('amazon.eg/')){//埃及
    country = 'EGP'
}else if(URL.includes('amazon.sa/')){//沙特
    country = 'SAR'
}else if(URL.includes('amazon.ae/')){//阿联酋
    country = 'AED'
}
console.log('确定的国家为:'+country)


//====统一监听body的改变，触发总回调
let callback = async function (records){
    //判断是否为默认货币种类
    let a = await cookieStore.get('i18n-prefs')
    if(!a){ //没有默认货币种类-直接按域名设置的货币种类
        all()
    }else{
        let b = a.value
        if(b ===country){ //默认货币种类为当前域名的货币种类
            all()
        }else{//默认货币种类不是当前域名的货币种类
            console.log('%c 默认货币种类不是当前域名的货币种类!!!-不进行价格转换', 'color:red')
        }
    }
};
let throttle_callback = _.throttle(callback,2200);
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
let amazon_nodes = []

//总回调
const all=function(){
    //console.log('总回调启动');
    chrome.storage.local.get(["my_rate",'amazon_nodes','color'],function(result){
        rate = result.my_rate[`rate_${country}`];
        amazon_nodes = result.amazon_nodes
        if(result.color){
            color = result.color
        }
        foreach_nodes(amazon_nodes)
    //特别价格处理
    let e1 = document.querySelectorAll('span[class="a-offscreen"]')
        if(e1[0]){
            t1(e1)
        }
    })
}

let rg = /^\$(\d{1,3}\,){0,1}\d{1,3}\.\d{2}( \- \$(\d{1,3}\,){0,1}\d{1,3}\.\d{2}){0,1}$/

//计算人民币
const  getRmb = function(s){
        if(!s.includes('-')){
            return getRmb_of(s)
        }else{
            let [sa,sb] = s.split('-')
            let s1 = getRmb_of(sa)
            let s2 = getRmb_of(sb)
            if(s1 && s2){
                return s1 + ' - ¥'+ s2
            }else{
                return null
            }
            
        }
}

const getRmb_of = function(s){
    let r1 = /^[^0-9]*/i
    let r2 = /[^0-9]*$/i
    let r3 = /,/g
    let r4 = /\./g
    let r5 = / /g
    let r6 = /[^0-9\|]/g
    if(country === 'JPY' || country==='INR'){ //japan
        let s1 = s.replace(r1,'').replace(r2,'')
        let s2 = s1.replace(r3,'')
        let b = parseFloat(s2)
        let rmb = b/rate
        if(rmb){
            return rmb.toFixed(1)
        }else{
            return 0
        }

    }else{ //orther
        let s1 = s.replace(r1,'').replace(r2,'')
        let s_ = s1.slice(-3,-2)
        if(s_==='.' || s_===','){
            let s2 = s1.slice(0,-3)+ '|' +s1.slice(-2)
            //s3 = s2.replace(r3,'').replace(r4,'').replace(r5,'').replace(/\|/g,'.')
            let s3 = s2.replace(r6,'').replace(/\|/g,'.')
            let b = parseFloat(s3)
            let rmb = b/rate
            if(rmb){
                return rmb.toFixed(1)
            }else{
                return 0
            }
        }else{
            return 0
        }
    }

}
//取出元素数组的元素处理
const foreach_nodes = function(node_all){

    for(let node of node_all){
        let a = document.querySelectorAll(`${node[0]}[${node[1]}="${node[2]}"]`)

        if(a[0]){
            changePriceOfTheOneNode(a)
        }
    }

}

const changePriceOfTheOneNode=function(nodes){
    nodes.forEach(e=>{
        let c=e.innerHTML.trim()
        if((!e.nextElementSibling) || (e.nextElementSibling && e.nextElementSibling.tagName != 'SUB')){
            let rmb = getRmb(c)
            if(rmb){
                let b = document.createElement('sub')
                    b.setAttribute('translate', 'no');
                    b.style.color = color
                    b.innerText = '¥'+rmb
                    e.after(b)
            }
        }
    })
}

const t1=function(nodes){
    nodes.forEach(x=>{
        let s = x.innerHTML.trim();
        if(s === 'price'){//非常特别的价格
            let a = x.nextElementSibling
            let b = a.innerText
            let c = b.replace(/[^0-9\,\.]/g,'')
            let d = c.replace(',','.')
            //console.log(d)
            if(x.nextElementSibling&&x.nextElementSibling.lastElementChild){
                if(!x.nextElementSibling.lastElementChild.innerHTML.includes('¥')){
                 let rmb = getRmb(d);
                 if(rmb){
                     let b = document.createElement('sub');
                     b.setAttribute('translate', 'no');
                     b.style.color = color;
                     b.innerText = '¥'+rmb;
                     x.nextElementSibling.appendChild(b);
                 }
                
                }
             }



        }else{
            if(x.nextElementSibling&&x.nextElementSibling.lastElementChild){
                if(!x.nextElementSibling.lastElementChild.innerHTML.includes('¥')){
                 let rmb = getRmb(s);
                 if(rmb){
                     let b = document.createElement('sub');
                     b.setAttribute('translate', 'no');
                     b.style.color = color;
                     b.innerText = '¥'+rmb;
                     x.nextElementSibling.appendChild(b);
                 }
                
                }
             }
        }

    })
}

