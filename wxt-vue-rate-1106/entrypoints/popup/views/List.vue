<template>
    <div class="list">
      <div  ref="el">
        <!-- <TransitionGroup name="rateList"> -->

        <div class="rate_item" v-for="item in country_list_rate" :key="item.tcur">
            <img :src="
                'icons/flags/4x3/' +
                item.tcur.slice(0, 2).toLowerCase() +
                '.svg'
              " alt="" width="50" class="handle1 flage_img"/>
            <div class="rate_item_in">
              <div class="rate_item_in_symbol"> {{ item.currency_symbol }}</div>
              <el-input size="default" 
              v-model="item.num" 
              placeholder="0" 
              clearable
              :formatter="formatter"
              @input="numChange(item)"
              @clear="clear()"
              >
                <!-- <template #prepend><span style="
                      color: var(--list-item-symbol-color);
                      font-size: var(--list-item-symbol-color-font-size);
                      padding-left: 0px;
                      font-weight: bold;
                    ">
                    {{ item.currency_symbol }}
                  </span></template> -->
              </el-input>
              <div class="item-hot-f">
              <span style="white-space: nowrap;">{{item.currency_name}}</span>
              <span style="white-space: nowrap;" class="item-hot-f-c" v-if="layout == 'default'&&show100 == 'yes'">
                100 {{item.tcur}} = {{ item.num_100}}
              </span>
            </div>
            </div>
            <span class="del_span" @click="delItem(item)" >
              <el-icon style="font-size: var(--list-item-text-del-span-font-size);"><CloseBold /></el-icon>
            </span>
          </div>
   
        <!-- </TransitionGroup> -->
        </div>
    </div>
</template>

<script setup lang="ts">
    import { ref,onMounted,onBeforeMount,onBeforeUnmount } from 'vue'
    import { useDraggable, type UseDraggableReturn } from 'vue-draggable-plus'

    const country_list = ref([])
    const country_list_rate = ref([])
    const layout = ref('default')
    const el = ref() //拖动列表的ref
    //定义记住的当前货币入数字
    const currentCurrency = ref({
      currency:'',
      num:0
    })
    //定义小数位置
    const decimal = ref(2)
    //定义是否显示千分位
    const thousandsF = ref('no')
    const show100 = ref('yes') //是否显示100货币量的转换
    onBeforeMount(() => {
      // console.log('onBeforeMount')
        chrome.storage.local.get(
          ["my_rate_zz", "my_rate_zz_countries", "decimal","current_currency",'thousandsF','layout','show100'],
          async function (result) {
            let { my_rate_zz, my_rate_zz_countries } = result;
            if (result.decimal){
              if (typeof result.decimal == 'number' && result.decimal >= 2 && result.decimal <= 6) {
                  decimal.value = result.decimal
                }else{
                  //删除原来的decimal
                  await chrome.storage.local.remove('decimal')
                }
            }
            if(result.thousandsF==='yes'){
              thousandsF.value = 'yes'
            }
            if(result.layout){
              layout.value = result.layout
            }
            if(result.show100==='no'){
              show100.value = 'no'
            }else{
              show100.value = 'yes'
            }
            // console.log('my_rate_zz_countries',my_rate_zz_countries)
            if (!my_rate_zz_countries) {
              my_rate_zz_countries = ["CNY", "USD", "TWD","THB","JPY"];
              chrome.storage.local.set(
                { my_rate_zz_countries},
                function (s) {
                  console.log("初始化国家列表");
                }
              );
            }
            country_list.value = my_rate_zz_countries;
            let c = [];
            country_list.value.forEach((el) => {
              let rate_country = {};
              rate_country.country_img = my_rate_zz[el].flagURL;
              rate_country.currency_symbol = my_rate_zz[el].symbol;
              rate_country.num = "";
              rate_country.currency_name = my_rate_zz[el].name;
              rate_country.currency_status = "";
              rate_country.tcur = my_rate_zz[el].tcur;
              rate_country.rate = my_rate_zz[el].rate;
              rate_country.hot = my_rate_zz[el].hot;
              rate_country.num_100 = "" // 100对应的当前货币数值 可以设置关闭
              c.push(rate_country);
            });
            country_list_rate.value = c;
            //从之前保存的当前货币中恢复数字
            if(result.current_currency&&result.current_currency.currency&&result.current_currency.num){
              if(country_list.value.includes(result.current_currency.currency)){
                let num = parseFloat(result.current_currency.num)
                if(num>0){
                  // console.log('num',num)
                  country_list_rate.value.forEach((el)=>{
                    if(el.tcur===result.current_currency.currency){
                      el.num = result.current_currency.num
                      numChange(el)
                      currentCurrency.value = {
                        currency:result.current_currency.currency,
                        num:result.current_currency.num
                      }
                    }
                  })
                  
                }else{
                  // console.log('NO,3')
                }
              }else{
                // console.log('NO,2')
              }
            }else{
              // console.log('NO,1')
            }

          }
        );

    })
    onMounted(() => {
      // console.log('1:',document.documentElement.dataset.layout)
      // isDefaultLayout.value = document.documentElement.dataset.layout === 'default';
    })

    onBeforeUnmount(()=>{
  })
    const numChange=(val)=>{
      // console.log(val.num)
      //处理粘贴进来的带千分号的情况
       val.num = val.num.trim().replace(/[\s,]+/g, '');
        country_list_rate.value.forEach((el) => {
          if (val.tcur === el.tcur) {
            el.num_100 =    "100" + val.tcur;
          } else {
            el.num = ((Number(el.rate) / Number(val.rate)) * Number(val.num)).toFixed(decimal.value);
            el.num_100 =   ((Number(val.rate) / Number(el.rate)) *100).toFixed(decimal.value) +" " + val.tcur;
          }
        });
        //记录输入的货币与数字
        currentCurrency.value.currency = val.tcur,
        currentCurrency.value.num = val.num
        chrome.storage.local.set({current_currency:currentCurrency.value},function(){
        // console.log('设置current_currency成功')
      })
    }
      //格式化input输入框的值
  const formatter = (value) => {
    if (!value) return "";
    if ( thousandsF.value != 'yes'){ return value   }
    let [decimalSeparator, thousandsSeparator] = './,'.split("/");
    if (!thousandsSeparator) {
      thousandsSeparator = "";
    }
    if (thousandsSeparator === " ") {
      thousandsSeparator = "\u00A0";
    }
    const parts = value.toString().split(".");
    parts[0] = parts[0]
      .replaceAll(",", "")
      .replaceAll(/(\d)(?=(?:\d{3})+$)/g, "$1"+thousandsSeparator);
    return parts.join(".");
  }
    //清空输入框时,去掉保存的当前货币与数字
    const clear = ()=>{
      chrome.storage.local.set({current_currency:{
        currency:'',
        num:0
      }},function(){
        // console.log('清空current_currency成功')
      })
    }
    //删除列表中的某一项
    const delItem=(val)=>{
      // console.log(val)
      country_list_rate.value.forEach((el,index) => {
        if (val.tcur === el.tcur) {
          country_list_rate.value.splice(index,1)
        }
      });
      //把当前的列表存入storage
      let c = country_list_rate.value.map((el)=>{
        return el.tcur
      })
      chrome.storage.local.set({my_rate_zz_countries:c},function(){
        console.log('修改成功')
      })
    }
    //拖动列表
    const draggable = useDraggable<UseDraggableReturn>(el, country_list_rate, {
        animation: 150,
        handle:'.handle1',
        onStart() {

        },
        onUpdate() {
          //把当前拖动后的列表保存到storage
          let c = country_list_rate.value.map((el)=>{
            return el.tcur
          })
          chrome.storage.local.set({my_rate_zz_countries:c},function(){
            // console.log('修改成功')
          })
        }
      })

</script>

<style>

.list {
    margin: var(--list-margin); 
    /* background-color: rgb(39, 46, 56); */
    /* height: 120px; */
  }
  .rate_item {
    margin: var(--list-rate-item-margin);
    width: var(--list-rate-item-width);
    height: var(--list-rate-item-height);

    background: var(--list-item-backgroud-color);

    border-radius: var(--list-rate-item-border-radius);
    display: grid;
    /* grid-template-columns: 80px auto 35px; */
    grid-template-columns: var(--list-rate-item-grid1-template-columns);
    grid-template-rows: 1fr;
    place-items: center center;
  }

  .rate_item:hover {
    background: var(--list-item-backgroud-color-hover);
  }
  .flage_img{
    border-radius: 5px;
    transition: transform 0.5s ease;
    /* width: 50px; */
    /* height: 50px; */
    overflow: hidden;
    box-shadow: 3px 3px 10px 0 rgba(0,0,0,0.4);
  }

  .flage_img:hover {
    transform: scale(1.1);
  }
  .rate_item:hover .flage_img {
    transform: scale(1.1);
  }
  .rate_item_in {
    display: grid;
    grid-template-rows:  5fr 4fr;
    grid-template-columns: 1fr 5fr;
    place-items: center start;
  }
  .rate_item_in_symbol{
    color: var(--list-item-symbol-color);
    font-size: var(--list-item-symbol-color-font-size);
    font-weight: bold;
    text-shadow: var(--list-item-symbol-text-shadow);/* 设置文字阴影 */
    grid-row-start: 1;
    grid-row-end: 3;

  }

  .rate_item_in .el-input-group__prepend {
    background-color: rgba(10, 31, 61, 0) ;
    box-shadow: none ;
}
  .rate_item_in .el-input__wrapper {
    background-color: rgba(53, 75, 105, 0);
    box-shadow: none;
  }
  .rate_item_in .el-input-group__prepend+.el-input__wrapper{
    border-radius: 5px;
  }
  .rate_item_in .el-input__inner {
    color: var(--list-item-text-input);
    font-family:Verdana, Geneva, Tahoma, sans-serif;
    font-size: var(--list-rate-item-in-input-font-size);
  }
  .rate_item_in i svg{
    color: var(--list-item-clear-icon-color)
  }
  .item-hot-f {
  display: grid;
  width: 100%;
  grid-template-rows: 3fr 3fr;
  align-items: center;
}

.item-hot-f > :first-child {
  justify-self: start; /* 第一个子元素靠左对齐 */


}

.item-hot-f > :last-child {
  justify-self: start; /* 第二个子元素靠右对齐 */
}
  .item-hot-f>span{
    color: var(--list-item-text-rate-text);
    font-size: var(--list-rate-item-in-span-font-size);
    margin-left: var(--list-rate-item-in-span-margin-left1);
    font-family:Verdana, Geneva, Tahoma, sans-serif;
  }
  .rate_item_in>div>span.hot{
    color: var(--span-hot-color);
  }

  .del_span {
    width: 18px;
    height: 18px;
    color: var(--list-item-text-del-span);
    display: none;
  }
  .rate_item:hover .del_span {
    display: block;
  }
  .del_span:hover {
    width: 20px;
    height: 20px;
    color: var(--list-item-text-del-span-hover);
    cursor: pointer;
  }
  .rateList-enter-active,
  .rateList-leave-active {
    transition: all 0.5s ease;
  }
  .rateList-enter-from,
  .rateList-leave-to {
    opacity: 0;
    transform: translateX(30px);
  }
</style>
