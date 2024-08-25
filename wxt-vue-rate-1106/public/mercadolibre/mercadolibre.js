console.log('%c美客多价格转换插件启动','color:red;font-size:14px;');

//按url确定货币
let URL = document.URL;
//console.log(URL);
let currency = ''; //当前面页的货币种类

//需要转换的货币种类,目标货币种类,默认转换为人民币
let ConverterCurreny = 'CNY';


if(URL.includes('mercadolibre.com.mx/') ){
    currency = 'MXN';
}else if( URL.includes('mercadolibre.cl/') ){
    currency = 'CLP';
}else if( URL.includes('mercadolivre.com.br/') ){
    currency = 'BRL';
}else if( URL.includes('mercadolibre.com.co/') ){
    currency = 'COP';
}

let rate = 0
let color = 'blueviolet'
let currenyTheme = 'blueviolet'
let mkd_nodes = [
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
    chrome.storage.local.get(["my_rate",'mkd_nodes','color','ConverterCurreny','currenyTheme'], function (result) {
        if(result.ConverterCurreny){
            ConverterCurreny = result.ConverterCurreny;
        }
        if(result.currenyTheme){
            currenyTheme = result.currenyTheme;
        }
        rate_b = result.my_rate[`rate_${currency}`];
        rate_c = result.my_rate[`rate_${ConverterCurreny}`];
        rate = rate_b / rate_c;
        mkd_nodes = result.mkd_nodes
        if(result.color){
            color = result.color
        }
        find_node(mkd_nodes);
    })
}

//找出元素
function find_node(node_all) {
    for (let item of node_all) {
        let a = document.querySelectorAll(item[0])
        if (a[0]) {
            if(item[1] === 1){
                insetHtml2(a);  //这种元素，在aria-label属性中取得金额的数字
            }
            if(item[1] === 2){
                
                insetHtml2(a,item[2],item[3]);  //这种元素中没有小数点
            }
        }
    }
}

function insetHtml2(nodes,integerChildIndex, decimalChildIndex){
    nodes.forEach(e => {
        let rmb = ''
        if(integerChildIndex !== undefined && decimalChildIndex !== undefined){
            rmb = extractNumber(e, integerChildIndex, decimalChildIndex)
        }else{
            let s = e.getAttribute('aria-label');
            rmb = getRmb(s)
        }
        
        if(!rmb){
            return
        }
        let nextSibling = e.nextElementSibling;
        if (nextSibling && (nextSibling.tagName.toLowerCase() === 'span') && nextSibling.classList.contains('rmb')) {
            if(nextSibling.querySelector('.rmb_num')){
                let s = nextSibling.querySelector('.rmb_num').innerText
                s = s.replace(/^\D+|\D+$/g, '');
                let rmb_old = parseFloat(s)
                if(rmb_old !=rmb){
                    nextSibling.querySelector('.rmb_num').innerText =  rmb
                }
                return;
            }

            
          } else {
                    let element_my = document.createElement('span')
                    element_my.classList.add('rmb')
                    element_my.style.display = 'inline-flex'
                    element_my.style.borderRadius = '8px'
                    element_my.style.fontSize = '14px'
                    element_my.style.fontWeight =  '500'
                    element_my.style.overflow = 'hidden'
                    // element_my.style["-webkit-text-fill-color"] = color
                    element_my.setAttribute('translate', 'no');
                    let element_a = document.createElement('span')
                    element_a.style.padding = '2px 2px 2px 5px'
                    element_a.style.color = '#FFFFFF'
                    element_a.style.backgroundColor = currenyTheme
                    element_a.style.display = 'flex'
                    element_a.style.alignItems = 'center'
                    element_a.style.height = '16px'
                    element_a.style.fontSize = '12px'
                    element_a.innerText = ConverterCurreny
                    let element_b = document.createElement('span')
                    element_b.classList.add('rmb_num')
                    element_b.style.padding = '2px 5px 2px 2px'
                    element_b.style.color = currenyTheme
                    element_b.style.backgroundColor = '#EEEEEE'
                    element_b.style.display = 'flex'
                    element_b.style.alignItems = 'center'
                    element_b.style.height = '16px'
                    element_b.style.fontSize = '14px'
                    element_my.style.fontWeight =  '600'
                    element_b.innerText = rmb

                    element_my.appendChild(element_a)
                    element_my.appendChild(element_b)
                    //把span插入到元素中
                    e.insertAdjacentElement('afterend', element_my);           
          }
    })
}
function getRmb(str){
    // 使用正则表达式匹配所有数字
    const matches = str.match(/\d+/g);
    if (!matches) return null; // 如果没有匹配到数字，返回null

    // 提取整数部分
    const integerPart = matches[0];

    // 尝试提取小数部分（如果存在）
    const decimalPart = matches[1] || "00"; // 如果没有小数部分，默认为"00"

    // 将整数部分和小数部分组合成一个浮点数
    const value = parseFloat(`${integerPart}.${decimalPart}`);
    const rmb = value / rate
    if(value && value>0){
        return rmb.toFixed(2)
    }else{
        return 0
    }
    
}

 //根据给出的父元素，以及整数位置入小数僧来计算人民币价格
  function extractNumber(parentElement, integerChildIndex, decimalChildIndex) {
  
    // 获取所有子元素
    const children = Array.from(parentElement.children);
  
    // 通过索引选中整数部分和小数部分的子元素
    // 整数部分
    const integerPartElement = children[integerChildIndex];
    let decimalPartElement = null
    //如果小数部分为no字符，表示没有小数部分
    if( decimalChildIndex != 'no'){
        decimalPartElement = children[decimalChildIndex];
    }

  
    if (!integerPartElement ) {
      console.error('Integer part element not found');
      return null;
    }
  
    // 提取数字文本内容，并去除非数字字符
    const integerPart = integerPartElement.textContent.replace(/[^\d]/g, '');
  
    let decimalPart = '0'

    if(decimalPartElement){
        decimalPart =  decimalPartElement.textContent.replace(/[^\d]/g, '') 
    }

  
    // 组合成一个完整的数值
    const fullNumber = parseFloat(`${integerPart}.${decimalPart}`);
  
    let rmb = fullNumber/rate
    if(rmb){
        return rmb.toFixed(2)
    }else{
        return ''
    }
  }
