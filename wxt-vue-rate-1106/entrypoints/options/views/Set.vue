<template>
    <div class="mt-5 mb-3">
    <el-tabs type="border-card" class="demo-tabs">
        <el-tab-pane>
          <template #label>
            <span class="custom-tabs-label">
              <el-icon><calendar /></el-icon>
              <span>站点启用</span>
            </span>
          </template>
          <Switch></Switch>
        </el-tab-pane>
        <el-tab-pane>
          <template #label>
            <span class="custom-tabs-label">
                <el-icon><Tools /></el-icon>
              <span>设置</span>
            </span>
          </template> 
          <Config></Config>
        </el-tab-pane>
        <el-tab-pane>
          <template #label>
            <span class="custom-tabs-label">
                <el-icon><ChatLineSquare /></el-icon>
              <span>日志</span>
            </span>
          </template>
          <LogList></LogList>
        </el-tab-pane>
      </el-tabs>
    </div>
</template>
<script setup>
    import Switch from './Switch.vue'
    import Config from './Config.vue'
    import LogList from './Log.vue'
</script>

<style>

</style>