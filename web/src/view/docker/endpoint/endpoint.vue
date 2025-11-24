<template>
  <div>
    <el-card shadow="never">
      <template #header>
        <div style="display:flex;align-items:center;gap:8px">
          <span>Docker SDK 配置</span>
          <el-input v-model="searchInfo.label" placeholder="标签" style="width:160px" />
          <el-input v-model="searchInfo.endpoint" placeholder="地址" style="width:260px" />
          <el-select v-model="searchInfo.status" placeholder="状态" style="width:140px" clearable>
            <el-option label="正常" value="正常" />
            <el-option label="异常" value="异常" />
          </el-select>
          <el-button type="primary" @click="getList">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
          <el-button type="primary" @click="openDialog('create')">新增</el-button>
          <el-button type="danger" :disabled="multipleSelection.length===0" @click="handleBatchDelete">批量删除</el-button>
        </div>
      </template>
      <el-table :data="tableData" height="calc(100vh - 280px)" @selection-change="onSelectionChange">
        <el-table-column type="selection" width="48" />
        <el-table-column prop="label" label="标签" width="160" />
        <el-table-column prop="endpoint" label="地址" min-width="260" />
        <el-table-column prop="useTls" label="TLS" width="80">
          <template #default="scope">{{ scope.row.useTls ? '是' : '否' }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" />
        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-button size="small" @click="openDialog('update', scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="getList" @size-change="getList" />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogType==='create'?'新增':'编辑'" width="600px">
      <el-form :model="formData" label-width="120px">
        <el-form-item label="标签" prop="label">
          <el-input v-model="formData.label" />
        </el-form-item>
        <el-form-item label="连接地址" prop="endpoint">
          <el-input v-model="formData.endpoint" placeholder="例如 https://1.2.3.4:2376 或 tcp://..." />
        </el-form-item>
        <el-form-item label="使用TLS" prop="useTls">
          <el-switch v-model="formData.useTls" />
        </el-form-item>
        <el-form-item label="CA证书" prop="caCert">
          <el-input v-model="formData.caCert" type="textarea" autosize placeholder="PEM 格式" />
        </el-form-item>
        <el-form-item label="客户端证书" prop="clientCert">
          <el-input v-model="formData.clientCert" type="textarea" autosize placeholder="PEM 格式" />
        </el-form-item>
        <el-form-item label="客户端私钥" prop="clientKey">
          <el-input v-model="formData.clientKey" type="textarea" autosize placeholder="PEM 格式" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-input v-model="formData.status" disabled placeholder="保存后自动检测" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible=false">取消</el-button>
        <el-button type="primary" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </div>
 </template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getDockerEndpointList, createDockerEndpoint, updateDockerEndpoint, deleteDockerEndpoint } from '@/api/docker/endpoint'

const searchInfo = ref({ label: '', endpoint: '', status: '' })
const tableData = ref([])
const multipleSelection = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const getList = async () => {
  try {
    const res = await getDockerEndpointList({ ...searchInfo.value, page: page.value, pageSize: pageSize.value })
    if (res.code === 0) {
      tableData.value = res.data.list
      total.value = res.data.total
    } else {
      ElMessage.error(res.msg || '获取失败')
    }
  } catch (e) { ElMessage.error('获取失败') }
}

const resetSearch = () => { searchInfo.value = { label: '', endpoint: '', status: '' }; page.value = 1; getList() }

const dialogVisible = ref(false)
const dialogType = ref('create')
const formData = ref({ label: '', endpoint: '', useTls: true, caCert: '', clientCert: '', clientKey: '', status: '' })
const openDialog = (type, row) => {
  dialogType.value = type
  dialogVisible.value = true
  formData.value = type === 'create' ? { label: '', endpoint: '', useTls: true, caCert: '', clientCert: '', clientKey: '', status: '' } : { ...row }
}

const submit = async () => {
  try {
    const api = dialogType.value === 'create' ? createDockerEndpoint : updateDockerEndpoint
    const res = await api(formData.value)
    if (res.code === 0) {
      ElMessage.success('保存成功')
      dialogVisible.value = false
      getList()
    } else { ElMessage.error(res.msg || '保存失败') }
  } catch (e) { ElMessage.error('保存失败') }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确认删除该配置？', '提示')
    const res = await deleteDockerEndpoint({ ID: row.id })
    if (res.code === 0) { ElMessage.success('删除成功'); getList() } else { ElMessage.error(res.msg || '删除失败') }
  } catch (e) {}
}

const onSelectionChange = (rows) => { multipleSelection.value = rows }
const handleBatchDelete = async () => {
  if (multipleSelection.value.length === 0) return
  try {
    await ElMessageBox.confirm('确认批量删除选中配置？', '提示')
    const ids = multipleSelection.value.map(i => i.id)
    const params = new URLSearchParams()
    ids.forEach(id => params.append('IDs[]', id))
    const res = await fetch('/docker/deleteDockerEndpointByIds?' + params.toString(), { method: 'DELETE' }).then(r => r.json())
    if (res.code === 0) { ElMessage.success('批量删除成功'); getList() } else { ElMessage.error(res.msg || '批量删除失败') }
  } catch (e) {}
}

onMounted(() => { getList() })
</script>
