console.log("Gmarket启动..");
const country = "KRW";

//====统一监听body的改变，触发总回调
let callback = function (records) {
    // console.log('callback..')
  all();
};
let throttle_callback = _.throttle(callback, 2600);
let mo = new MutationObserver(throttle_callback);
let option = {
  childList: true,
  subtree: true,
};
let fs_node = document.getElementsByTagName("body")[0];
try {
  mo.observe(fs_node, option);
} catch (e) {
  console.log('监听器启动失败body."');
}

let rate = 0;
let gmarket_nodes = [];
let color = "green";


//总回调
const all = function () {
  // console.log("总回调启动");
  chrome.storage.local.get(
    ["my_rate", "gmarket_nodes", "color"],
    function (result) {
      rate = result.my_rate[`rate_${country}`];
      gmarket_nodes = result.gmarket_nodes;

      if (result.color) {
        color = result.color;
      }
      foreach_nodes(gmarket_nodes);
    }
  );
};

//取出元素数组的元素处理
const foreach_nodes = function (node_all) {

  for (let node of node_all) {
    //这里判断node是三个元素还是一个元素,组成不周的query查询
    let query = ''
    if(node.length===3){
      query = `${node[0]}[${node[1]}="${node[2]}"]`
    }else if(node.length===1){
      query = node[0]
    }else{
      console.log('节点查询信息不正常')
      return
    }
    let a = document.querySelectorAll(query);

    if (a[0]) {
      changePriceOfTheOneNode(a);
    }
  }

};

const changePriceOfTheOneNode = function (nodes) {
  nodes.forEach((e) => {
    let c = e.innerHTML.trim();
    if (
      !e.nextElementSibling ||
      (e.nextElementSibling && e.nextElementSibling.tagName != "SUB" &&e.nextElementSibling.innerText!="%")
    ) {
      let rmb = getRmb(c);
      if (rmb) {
        let b = document.createElement("sub");
        b.setAttribute('translate', 'no');
        b.style.color = color;
        b.style.fontSize="14px"
        b.innerText = "¥" + rmb;
        e.after(b);
      }
    }
  });
};



//计算人民币
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
