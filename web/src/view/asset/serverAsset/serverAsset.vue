
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="创建日期" prop="createdAtRange">
      <template #label>
        <span>
          创建日期
          <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
            <el-icon><QuestionFilled /></el-icon>
          </el-tooltip>
        </span>
      </template>

      <el-date-picker
            v-model="searchInfo.createdAtRange"
            class="!w-380px"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
       </el-form-item>
      
            <el-form-item label="服务器名称" prop="name">
  <el-input v-model="searchInfo.name" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="主机名" prop="hostname">
  <el-input v-model="searchInfo.hostname" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="IP地址" prop="ip">
  <el-input v-model="searchInfo.ip" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="操作系统类型" prop="osType">
    <el-tree-select v-model="searchInfo.osType" placeholder="请选择操作系统类型" :data="os_typeOptions" style="width:100%" filterable :clearable="true" check-strictly ></el-tree-select>
</el-form-item>
            
            <el-form-item label="GPU数量" prop="gpuCount">
  <el-input v-model.number="searchInfo.gpuCount" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="资产状态" prop="status">
    <el-tree-select v-model="searchInfo.status" placeholder="请选择资产状态" :data="asset_statusOptions" style="width:100%" filterable :clearable="true" check-strictly ></el-tree-select>
</el-form-item>
            

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
          <el-form-item label="端口" prop="port">
  <el-input v-model.number="searchInfo.port" placeholder="搜索条件" />
</el-form-item>
          
          <el-form-item label="CPU核数" prop="cpuCores">
  <el-input v-model.number="searchInfo.cpuCores" placeholder="搜索条件" />
</el-form-item>
          
          <el-form-item label="内存容量GB" prop="memoryGB">
  <el-input v-model.number="searchInfo.memoryGB" placeholder="搜索条件" />
</el-form-item>
          
          <el-form-item label="磁盘容量GB" prop="diskGB">
  <el-input v-model.number="searchInfo.diskGB" placeholder="搜索条件" />
</el-form-item>
          
          <el-form-item label="地区" prop="region">
  <el-input v-model="searchInfo.region" placeholder="搜索条件" />
</el-form-item>
          
          <el-form-item label="标签" prop="label">
  <el-input v-model="searchInfo.label" placeholder="搜索条件" />
</el-form-item>
          
          <el-form-item label="备注" prop="remark">
  <el-input v-model="searchInfo.remark" placeholder="搜索条件" />
</el-form-item>
          
          <el-form-item label="关联SSH服务器" prop="sshServerId">
  <el-select v-model="searchInfo.sshServerId" filterable placeholder="请选择关联SSH服务器" :clearable="true">
    <el-option v-for="(item,key) in dataSource.sshServerId" :key="key" :label="item.label" :value="item.value" />
  </el-select>
</el-form-item>
          
          <el-form-item label="关联Docker端点" prop="endpointId">
  <el-select v-model="searchInfo.endpointId" filterable placeholder="请选择关联Docker端点" :clearable="true">
    <el-option v-for="(item,key) in dataSource.endpointId" :key="key" :label="item.label" :value="item.value" />
  </el-select>
</el-form-item>
          
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery=true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery=false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button  type="primary" icon="plus" @click="openDialog()">新增</el-button>
            <el-button  icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
            <ExportTemplate  template-id="asset_ServerAsset" />
            <ExportExcel  template-id="asset_ServerAsset" filterDeleted/>
            <ImportExcel  template-id="asset_ServerAsset" @on-success="getTableData" />
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        @sort-change="sortChange"
        >
        <el-table-column type="selection" width="55" />
        
        <el-table-column sortable align="left" label="日期" prop="CreatedAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        
            <el-table-column sortable align="left" label="服务器名称" prop="name" width="120" />

            <el-table-column align="left" label="主机名" prop="hostname" width="120" />

            <el-table-column align="left" label="IP地址" prop="ip" width="120" />

            <el-table-column align="left" label="端口" prop="port" width="120" />

            <el-table-column align="left" label="操作系统类型" prop="osType" width="120">
    <template #default="scope">
    {{ filterDict(scope.row.osType,os_typeOptions) }}
    </template>
</el-table-column>
            <el-table-column sortable align="left" label="CPU核数" prop="cpuCores" width="120" />

            <el-table-column sortable align="left" label="内存容量GB" prop="memoryGB" width="120" />

            <el-table-column sortable align="left" label="磁盘容量GB" prop="diskGB" width="120" />

            <el-table-column sortable align="left" label="GPU数量" prop="gpuCount" width="120" />

            <el-table-column sortable align="left" label="资产状态" prop="status" width="120">
    <template #default="scope">
    {{ filterDict(scope.row.status,asset_statusOptions) }}
    </template>
</el-table-column>
            <el-table-column sortable align="left" label="地区" prop="region" width="120" />

            <el-table-column sortable align="left" label="标签" prop="label" width="120" />

            <el-table-column align="left" label="关联SSH服务器" prop="sshServerId" width="120">
    <template #default="scope">
        <span>{{ filterDataSource(dataSource.sshServerId,scope.row.sshServerId) }}</span>
    </template>
</el-table-column>
            <el-table-column align="left" label="关联Docker端点" prop="endpointId" width="120">
    <template #default="scope">
        <span>{{ filterDataSource(dataSource.endpointId,scope.row.endpointId) }}</span>
    </template>
</el-table-column>
        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            <el-button  type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
            <el-button  type="primary" link icon="edit" class="table-button" @click="updateServerAssetFunc(scope.row)">编辑</el-button>
            <el-button   type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
            </template>
        </el-table-column>
        </el-table>
        <div class="gva-pagination">
            <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            />
        </div>
    </div>
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
       <template #header>
              <div class="flex justify-between items-center">
                <span class="text-lg">{{type==='create'?'新增':'编辑'}}</span>
                <div>
                  <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
                  <el-button @click="closeDialog">取 消</el-button>
                </div>
              </div>
            </template>

          <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
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
          </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看">
            <el-descriptions :column="1" border>
                    <el-descriptions-item label="服务器名称">
    {{ detailForm.name }}
</el-descriptions-item>
                    <el-descriptions-item label="主机名">
    {{ detailForm.hostname }}
</el-descriptions-item>
                    <el-descriptions-item label="IP地址">
    {{ detailForm.ip }}
</el-descriptions-item>
                    <el-descriptions-item label="端口">
    {{ detailForm.port }}
</el-descriptions-item>
                    <el-descriptions-item label="操作系统类型">
    {{ detailForm.osType }}
</el-descriptions-item>
                    <el-descriptions-item label="CPU核数">
    {{ detailForm.cpuCores }}
</el-descriptions-item>
                    <el-descriptions-item label="内存容量GB">
    {{ detailForm.memoryGB }}
</el-descriptions-item>
                    <el-descriptions-item label="磁盘容量GB">
    {{ detailForm.diskGB }}
</el-descriptions-item>
                    <el-descriptions-item label="GPU数量">
    {{ detailForm.gpuCount }}
</el-descriptions-item>
                    <el-descriptions-item label="资产状态">
    {{ detailForm.status }}
</el-descriptions-item>
                    <el-descriptions-item label="地区">
    {{ detailForm.region }}
</el-descriptions-item>
                    <el-descriptions-item label="标签">
    {{ detailForm.label }}
</el-descriptions-item>
                    <el-descriptions-item label="备注">
    <RichView v-model="detailForm.remark" />
</el-descriptions-item>
                    <el-descriptions-item label="关联SSH服务器">
    <template #default="scope">
        <span>{{ filterDataSource(dataSource.sshServerId,detailForm.sshServerId) }}</span>
    </template>
</el-descriptions-item>
                    <el-descriptions-item label="关联Docker端点">
    <template #default="scope">
        <span>{{ filterDataSource(dataSource.endpointId,detailForm.endpointId) }}</span>
    </template>
</el-descriptions-item>
            </el-descriptions>
        </el-drawer>

  </div>
</template>

<script setup>
import {
    getServerAssetDataSource,
  createServerAsset,
  deleteServerAsset,
  deleteServerAssetByIds,
  updateServerAsset,
  findServerAsset,
  getServerAssetList
} from '@/api/asset/serverAsset'
// 富文本组件
import RichEdit from '@/components/richtext/rich-edit.vue'
import RichView from '@/components/richtext/rich-view.vue'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"

// 导出组件
import ExportExcel from '@/components/exportExcel/exportExcel.vue'
// 导入组件
import ImportExcel from '@/components/exportExcel/importExcel.vue'
// 导出模板组件
import ExportTemplate from '@/components/exportExcel/exportTemplate.vue'


defineOptions({
    name: 'ServerAsset'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
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
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getServerAssetDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()



// 验证规则
const rule = reactive({
               name : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
               ip : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
// 排序
const sortChange = ({ prop, order }) => {
  const sortMap = {
    CreatedAt:"created_at",
    ID:"id",
            name: 'name',
            cpuCores: 'cpu_cores',
            memoryGB: 'memory_gb',
            diskGB: 'disk_gb',
            gpuCount: 'gpu_count',
            status: 'status',
            region: 'region',
            label: 'label',
  }

  let sort = sortMap[prop]
  if(!sort){
   sort = prop.replace(/[A-Z]/g, match => `_${match.toLowerCase()}`)
  }

  searchInfo.value.sort = sort
  searchInfo.value.order = order
  getTableData()
}
// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await getServerAssetList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
    os_typeOptions.value = await getDictFunc('os_type')
    asset_statusOptions.value = await getDictFunc('asset_status')
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteServerAssetFunc(row)
        })
    }

// 多选删除
const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
      const IDs = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          IDs.push(item.ID)
        })
      const res = await deleteServerAssetByIds({ IDs })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === IDs.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
      })
    }

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateServerAssetFunc = async(row) => {
    const res = await findServerAsset({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteServerAssetFunc = async (row) => {
    const res = await deleteServerAsset({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: '删除成功'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
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
        }
}
// 弹窗确定
const enterDialog = async () => {
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
                closeDialog()
                getTableData()
              }
      })
}

const detailForm = ref({})

// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findServerAsset({ ID: row.ID })
  if (res.code === 0) {
    detailForm.value = res.data
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = {}
}


</script>

<style>

</style>
