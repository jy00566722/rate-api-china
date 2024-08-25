<script lang="ts" setup>
import { ref,onMounted } from 'vue';

defineProps({
  msg: String,
});

const rate = ref(0);
onMounted(()=>{
  //从storage中获取rate数据
  chrome.storage.local.get(['rate'], function(result) {
    console.log('Value currently is ' + result.rate);
    rate.value = result.rate;
  });
})
</script>

<template>
  <h2>{{ msg }}</h2>
  <h2>1BTC = {{rate}}元</h2>

</template>


<style scoped>
.read-the-docs {
  color: #888;
}
</style>
