<template>
    <div class="feedback">
    <Nav darp-e="no"></Nav>
    <h1>问题及建议
      <el-icon style="color: var(--feedback-icons-color);"><MagicStick /></el-icon>
    </h1>
    <h3>你的反馈非常重要!</h3>
    <div class="feedback-form">
    <el-form
    label-position="top"
    label-width="100px"
    :model="form"
    style="max-width: 460px"
  >
    <el-form-item label="您的称谓及联系方式">
      <el-input v-model="form.name" />
    </el-form-item>
    <el-form-item label="有问题的URL">
      <el-input v-model="form.url" 
      :rows="2"
      maxlength="200"
      show-word-limit
      type="textarea"/>
    </el-form-item>
    <el-form-item label="问题及意见">
      <el-input v-model="form.context" 
      :rows="3"
      maxlength="500"
      show-word-limit
      type="textarea"
      />
    </el-form-item>
    <el-row>
        <el-col :span="12"><el-button type="primary" text plain @click="postFeedBack">提交</el-button></el-col>
        <el-col :span="12"><el-button type="primary" text plain @click="$router.push('/')">返回</el-button></el-col>
      </el-row>


  </el-form>
</div>
</div>
</template>

<script setup>
    import { ref,onMounted,onBeforeMount ,reactive} from 'vue'
    import Nav from './Nav.vue'
    import { ElMessage } from 'element-plus'

    const form = reactive({
        name:"",//用户名称
        url:"", //问题URL
        context:"",
        time:""
    })
    const postFeedBack=async ()=>{
      if(form.name==""||form.context==""){
        ElMessage.error('请填写完整信息')
        return
      }
      form.time = new Date().toLocaleString()
      let url = 'https://rate.lizudi.top/v2/feedback'
      const a = await fetch(url,{
                method:'POST',
                body:JSON.stringify(form),
                headers: {
                  'content-type': 'application/json'
                }
          })
      const b = await a.json()
      if(b.code===20000&b.msg==='反馈成功'){
          form.name = ''
          form.url = ''
          form.context= ''
          form.time = ''
          ElMessage({message:'反馈成功!',type:'success'})
        }else{
          form.name = ''
          form.url = ''
          form.context= ''
          form.time = ''
          ElMessage({message:'反馈失败,请重试!',type:'warning'})
        }
    }

</script>

<style>
    .feedback{
        background-color: var(--main-backgroud-color);
        color: var(--feedback-title-color);
    }
    .feedback h3{
        color: var(--feedback-title-color);
    }
    .feedback-form {
        padding: 2px 20px 30px 20px;
    }
    .feedback-form .el-form-item__label{
        color: var(--feedback-text-color);
    }
</style>