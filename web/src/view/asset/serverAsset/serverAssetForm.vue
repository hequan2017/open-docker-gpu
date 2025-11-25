
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="服务器名称:" prop="name">
    <el-input v-model="formData.name" :clearable="true" placeholder="请输入服务器名称" />
</el-form-item>
        <el-form-item label="主机名:" prop="hostname">
    <el-input v-model="formData.hostname" :clearable="true" placeholder="请输入主机名" />
</el-form-item>
        <el-form-item label="IP地址:" prop="ip">
    <el-input v-model="formData.ip" :clearable="false" placeholder="请输入IP地址" />
</el-form-item>
        <el-form-item label="端口:" prop="port">
    <el-input v-model.number="formData.port" :clearable="true" placeholder="请输入端口" />
</el-form-item>
        <el-form-item label="操作系统类型:" prop="osType">
    <el-tree-select v-model="formData.osType" placeholder="请选择操作系统类型" :data="os_typeOptions" style="width:100%" filterable :clearable="true" check-strictly></el-tree-select>
</el-form-item>
        <el-form-item label="CPU核数:" prop="cpuCores">
    <el-input v-model.number="formData.cpuCores" :clearable="true" placeholder="请输入CPU核数" />
</el-form-item>
        <el-form-item label="内存容量GB:" prop="memoryGB">
    <el-input v-model.number="formData.memoryGB" :clearable="true" placeholder="请输入内存容量GB" />
</el-form-item>
        <el-form-item label="磁盘容量GB:" prop="diskGB">
    <el-input v-model.number="formData.diskGB" :clearable="true" placeholder="请输入磁盘容量GB" />
</el-form-item>
        <el-form-item label="GPU数量:" prop="gpuCount">
    <el-input v-model.number="formData.gpuCount" :clearable="true" placeholder="请输入GPU数量" />
</el-form-item>
        <el-form-item label="资产状态:" prop="status">
    <el-tree-select v-model="formData.status" placeholder="请选择资产状态" :data="asset_statusOptions" style="width:100%" filterable :clearable="true" check-strictly></el-tree-select>
</el-form-item>
        <el-form-item label="地区:" prop="region">
    <el-input v-model="formData.region" :clearable="true" placeholder="请输入地区" />
</el-form-item>
        <el-form-item label="标签:" prop="label">
    <el-input v-model="formData.label" :clearable="true" placeholder="请输入标签" />
</el-form-item>
        <el-form-item label="备注:" prop="remark">
    <RichEdit v-model="formData.remark"/>
</el-form-item>
        <el-form-item label="关联SSH服务器:" prop="sshServerId">
    <el-select v-model="formData.sshServerId" placeholder="请选择关联SSH服务器" filterable style="width:100%" :clearable="true">
        <el-option v-for="(item,key) in dataSource.sshServerId" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="关联Docker端点:" prop="endpointId">
    <el-select v-model="formData.endpointId" placeholder="请选择关联Docker端点" filterable style="width:100%" :clearable="true">
        <el-option v-for="(item,key) in dataSource.endpointId" :key="key" :label="item.label" :value="item.value" />
    </el-select>
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
    getServerAssetDataSource,
  createServerAsset,
  updateServerAsset,
  findServerAsset
} from '@/api/asset/serverAsset'

defineOptions({
    name: 'ServerAssetForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'
// 富文本组件
import RichEdit from '@/components/richtext/rich-edit.vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const os_typeOptions = ref([])
const asset_statusOptions = ref([])
const formData = ref({
            name: '',
            hostname: '',
            ip: '',
            port: 0,
            osType: '',
            cpuCores: 0,
            memoryGB: 0,
            diskGB: 0,
            gpuCount: 0,
            status: '',
            region: '',
            label: '',
            remark: '',
            sshServerId: undefined,
            endpointId: undefined,
        })
// 验证规则
const rule = reactive({
               name : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               ip : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getServerAssetDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findServerAsset({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    os_typeOptions.value = await getDictFunc('os_type')
    asset_statusOptions.value = await getDictFunc('asset_status')
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
               res = await createServerAsset(formData.value)
               break
             case 'update':
               res = await updateServerAsset(formData.value)
               break
             default:
               res = await createServerAsset(formData.value)
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
