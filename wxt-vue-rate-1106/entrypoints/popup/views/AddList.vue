<template>
  <div class="add_list">
    <Nav darp-e="no"></Nav>
    <el-row justify="space-between" align="middle">
      <el-col :span="18"> <h3 style="color: var(--add-list-title-color);">选中货币加入计算列表</h3></el-col>
      <el-col :span="6"><span 
        :class="selectedClassObject"
        @click="lookSelected"
        >{{selected_text}}</span></el-col>
    </el-row>
    <div class="add_list_search">
      <el-input size="default" v-model="searchKey" placeholder="搜索货币:国家/代号" @input="searchInput">
          <template #append><span>
            <el-icon><Search /></el-icon>
          
              </span>
            </template>
        </el-input>
      </div>

        <div class="add_list_all">
          <template v-for="(item,index) in country_list_status" :key="item.tcur">
            <template v-if="item.status==='ALREADY'">
          <div class="add_list_country_item"  @click="changStatus(item)">
            <span ><img :src="item.flagURL" alt="" width="40" height="40"></span>
            <span style="color:var(--add-list-item-text-color);font-size:14px;justify-self:start;padding-left: 20px;">{{item.name}}</span>
            <span>
              <el-icon v-if="item.selected" :size="20" style="color: var(--add-list-selected-color);"><Select /></el-icon>
            
            </span>
            </div>
          </template>
          </template>
          </div>


        <div class="add_country_button">
          <el-button type="danger"  circle @click="$router.push('/')">返回</el-button>
        </div>
  </div>
  </template>

  <script setup>
    import { ref,onMounted,onBeforeMount ,reactive,computed} from 'vue'
    import Nav from './Nav.vue'

    const searchKey = ref('') //搜索关键字
    const country_list = ref([]) //store中的所有汇率列表
    const country_selected = ref([]) //选中的汇率列表
    const country_list_status =ref([]) //增加是否选中的状态后的列表-正式要显示的列表,去除了状态不对
    const selected_status = ref(false) //是否显示已选中的列表
    const selected_text = ref('查看已选中')
    const selectedClassObject = computed(() => ({
          'selected_country': selected_status.value,
          'selected_country_not': !selected_status.value,
        }))

    onBeforeMount(()=>{
      chrome.storage.local.get(["my_rate_zz", "my_rate_zz_countries"],(result)=>{
        let { my_rate_zz, my_rate_zz_countries } = result;
        let c = [];//要返回的列表
        let my_rate_zz_keys = Object.keys(my_rate_zz);
        my_rate_zz_keys.forEach((el) => {
         
            let rate_country = {};
            rate_country.flagURL = 'icons/flags/4x3/'+my_rate_zz[el].tcur.slice(0,2).toLowerCase()+'.svg';
            rate_country.name = my_rate_zz[el].name; 
            rate_country.ratenm = my_rate_zz[el].ratenm;
            rate_country.scur = my_rate_zz[el].scur;
            rate_country.status = my_rate_zz[el].status;
            rate_country.tcur = my_rate_zz[el].tcur;
            rate_country.symbol = my_rate_zz[el].symbol;
            rate_country.update = my_rate_zz[el].update;
            rate_country.selected = my_rate_zz_countries.includes(el);        
            c.push(rate_country);
        

        });
        country_list_status.value = c;
        country_selected.value = my_rate_zz_countries;
        // console.log(1,country_selected.value)
          
      })
    })

    const searchInput = () => {
      // console.log(searchKey.value)
      country_list_status.value.forEach((el)=>{
        if(el.name.includes(searchKey.value)|| 
           el.name.includes(searchKey.value.toUpperCase())){
          el.status = 'ALREADY'
        }else{
          el.status = 'noShow'
        }
      })
    }
    const changStatus=(item)=>{

      item.selected = !item.selected;
      if(item.selected){
        country_selected.value.push(item.tcur);
      }else{
        country_selected.value.splice(country_selected.value.indexOf(item.tcur),1);
      }

      let c =country_selected.value.map((el)=>{
        return el
    
      })
      chrome.storage.local.set({my_rate_zz_countries:c},function(){
        // console.log('修改成功')
      })
    }
    const lookSelected=()=>{
      searchKey.value='' //清空搜索框
      if(selected_status.value){
        //处于查看选中的状态
        country_list_status.value.forEach((el)=>{
          el.status='ALREADY'
        })
        selected_text.value = '查看选中'
        selected_status.value = false;
          
      }else{
        //处于查看所有的状态
        country_list_status.value.forEach((el)=>{
        if(el.selected){
          el.status = 'ALREADY'
        }else{
          el.status = 'noShow'
        }
        selected_text.value = '查看所有'
        selected_status.value = true;
      })
    }
  }
  </script>

  <style>
    .add_list{
      margin-bottom: 30px;
      padding: 0 8px;
    }
    .add_list .el-input-group__append {
    background-color: rgba(0, 0, 0, 0) ;
    box-shadow: none ;
}
  .add_list .el-input__wrapper {
    background-color: rgba(0, 0, 0, 0);
    box-shadow: none;
    /* border: 1px solid #ffffff; */
    border-radius: 5px;
  }
  .add_list .el-input__wrapper input::placeholder{
    color: var(--add-list-search-placeholder-color);
  }
  .add_list .el-input__inner {
    color: var(--add-list-search-text-color);
    font-family:Verdana, Geneva, Tahoma, sans-serif;
    font-size: 16px;
  }
    .add_list_search span{
      height: 20px;
      width: 20px;
      font-size: 16px;
      color:  var(--add-list-search-icons-color);
    }
  .add_list_all{
    margin-top: 10px;
  }
    .add_list_country_item{
      width: 100%;
      height: 58px;
      display: grid;
      grid-template-columns: var(  --list-add-grid-template-columns);
      place-items: center center;
      background: var(--add-list-item-backgroud-color);
      margin: 5px 0;
    }
    .add_list_country_item:hover{
      background: var(--add-list-item-backgroud-color-hover);
    }
    .selected_country_not{
      color: var(--add-list-selected-color);
      font-size: 14px; cursor: pointer;
      font-style: oblique;
      font-family:'Gill Sans', 'Gill Sans MT', Calibri, 'Trebuchet MS', sans-serif;
    }
    .selected_country{
      color: var(--add-list-selected-not-color);
      font-size: 14px; cursor: pointer;
      font-style:initial;
      font-family:'Times New Roman', Times, serif
    }
  </style>