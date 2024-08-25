import { ref ,watch} from 'vue'
import { defineStore } from 'pinia'

export const useMessageStore = defineStore('message', ()=>{
    const message = ref([])
    // const setMessage = (msg)=>{
    //     message.value.push(msg)
    // }
    const layout = ref('small')
    const $reset=()=>{
        message.value = []
    }
    watch(
        message,
        (state) => {
          // 每当状态发生变化时，将整个 state 持久化到本地存储。
            // console.log(state)
        },
        { deep: true }
      )
    return {
        message,$reset,layout
    }
  })