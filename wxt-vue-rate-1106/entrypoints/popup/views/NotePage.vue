<template >
    <div class="note-page">
    <Nav darp-e="no"></Nav>
    <el-icon :size="40" style="color: var(--notePage-title-color);"><Message /></el-icon>
    <h2 style="color: var(--notePage-title-color);margin:0px 1px">消 息 列 表</h2>
    <div class="note-list">
        <div class="note-item" v-for="item in messageStore.message" :key="item.index">
            <el-alert
            :title="item.context"
            :type="item.type"
            show-icon
            >
            </el-alert>
        </div>
    </div>
    <el-button type="primary" text plain @click="$router.push('/')">返回</el-button>
    </div>
</template>

<script setup>
import { ref,onMounted,onBeforeMount } from 'vue'
import Nav from './Nav.vue'
import { useMessageStore } from '../pinia'

const messageStore = useMessageStore()

</script>

<style>
    .note-page{
        /* background-color: aquamarine; */
        /* margin-bottom: 30px; */
    }
    .note-list{
        padding: 10px 10px;
    }
    .note-item{
        margin: 10px 5px 10px 5px;
    }
    
</style>
