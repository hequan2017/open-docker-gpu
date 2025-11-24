
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="服务器IP地址:" prop="ip">
    <el-input v-model="formData.ip" :clearable="true" placeholder="请输入服务器IP地址" />
</el-form-item>
        <el-form-item label="SSH端口:" prop="port">
    <el-input v-model.number="formData.port" :clearable="true" placeholder="请输入SSH端口" />
</el-form-item>
        <el-form-item label="登录账号:" prop="username">
    <el-input v-model="formData.username" :clearable="true" placeholder="请输入登录账号" />
</el-form-item>
        <el-form-item label="服务器标签名:" prop="label">
    <el-input v-model="formData.label" :clearable="true" placeholder="请输入服务器标签名" />
</el-form-item>
        <el-form-item label="服务器地区:" prop="region">
    <el-input v-model="formData.region" :clearable="true" placeholder="请输入服务器地区" />
</el-form-item>
        <el-form-item label="登录密码:" prop="password">
    <el-input v-model="formData.password" type="password" :clearable="true" placeholder="请输入登录密码" />
</el-form-item>
        <el-form-item label="状态:" prop="status">
    <el-input v-model="formData.status" disabled placeholder="自动检测结果（只读）" />
</el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
  createSshServer,
  updateSshServer,
  findSshServer
} from '@/api/ssh/sshServer'

defineOptions({
    name: 'SshServerForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const formData = ref({
            ip: '',
            port: 22,
            username: 'root',
            label: '',
            region: '',
            password: '',
            status: '',
        })
// 验证规则
const rule = reactive({
               ip : [{
                   required: true,
                   message: '请输入IP地址',
                   trigger: ['input','blur'],
               }],
               port : [{
                   required: true,
                   message: '请输入端口',
                   trigger: ['input','blur'],
               }],
               username : [{
                   required: true,
                   message: '请输入账号',
                   trigger: ['input','blur'],
               }],
               label : [{
                    required: false,
                    message: '请输入服务器标签名',
                    trigger: ['input','blur'],
                }],
               region : [{
                    required: false,
                    message: '请输入服务器地区',
                    trigger: ['input','blur'],
                }],
               password : [{
                   required: true,
                   message: '请输入密码',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findSshServer({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return btnLoading.value = false
            let res
           switch (type.value) {
             case 'create':
               res = await createSshServer(formData.value)
               break
             case 'update':
               res = await updateSshServer(formData.value)
               break
             default:
               res = await createSshServer(formData.value)
               break
           }
           btnLoading.value = false
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
