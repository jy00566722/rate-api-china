<template>
  <div class="home-page">
    <Nav :darp-e="drapE"></Nav>
    <Note></Note>
    <List></List>
    <!-- 增加国家的按键 -->
    <div class="add_country_button">
      <el-button type="warning"  circle @click="$router.push('/add_list')">
        <el-icon><Plus /></el-icon>
      </el-button>
    </div>
  </div>
</template>

  <script setup>
    import { ref,onMounted,onBeforeMount,computed } from 'vue'
    import { useMessageStore } from '../pinia'
    import Nav from './Nav.vue'
    import Note from './Note.vue'
    import List from './List.vue'

    const messageStore = useMessageStore()
    // const drapE = ref('no')
    onBeforeMount(() => {

    })
    //用计算属性来获取messageStore中layout的值,然后设置给drapE
    const drapE = computed(() => {
      // console.log(messageStore.layout)
      return messageStore.layout==='small'?'no':'yes'
    })

  </script>

<style>
  .home-page{
    margin-bottom: 20px;
  }
  .add_country_button{
    position:fixed;
    bottom: 10px;
    right: 10px;
    z-index:99;
  }
  </style>

