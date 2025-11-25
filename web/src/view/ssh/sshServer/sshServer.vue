
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
      
            <el-form-item label="服务器IP地址" prop="ip">
  <el-input v-model="searchInfo.ip" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="SSH端口" prop="port">
  <el-input v-model.number="searchInfo.port" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="登录账号" prop="username">
  <el-input v-model="searchInfo.username" placeholder="搜索条件" />
</el-form-item>
            <el-form-item label="服务器标签名" prop="label">
  <el-input v-model="searchInfo.label" placeholder="搜索条件" />
</el-form-item>
            <el-form-item label="服务器地区" prop="region">
  <el-input v-model="searchInfo.region" placeholder="搜索条件" />
</el-form-item>
            

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
          <el-form-item label="登录密码" prop="password">
  <el-input v-model="searchInfo.password" placeholder="搜索条件" />
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
            
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @expand-change="onExpandChange"
        @selection-change="handleSelectionChange"
        @sort-change="sortChange"
        >
        <el-table-column type="expand">
          <template #default="scope">
            <div style="margin-bottom:8px; display:flex; gap:8px; align-items:center">
              <span>nvidia-smi 尾部</span>
              <el-input-number v-model="scope.row.__nvsmiTail" :min="10" :max="2000" />
              <el-switch v-model="scope.row.__nvsmiAuto" active-text="自动刷新" />
              <el-input-number v-model="scope.row.__nvsmiInterval" :min="2" :max="120" />
              <span>秒</span>
              <el-button size="small" :loading="scope.row.__nvsmiLoading" @click="refreshNvsmi(scope.row)">刷新</el-button>
            </div>
            <div class="logbox"><pre>{{ scope.row.__nvsmiText || '暂无输出' }}</pre></div>
          </template>
        </el-table-column>
        <el-table-column type="selection" width="55" />
        
        <el-table-column sortable align="left" label="日期" prop="CreatedAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        
            <el-table-column sortable align="left" label="服务器IP地址" prop="ip" width="200" />

            <el-table-column align="left" label="SSH端口" prop="port" width="120" />

            <el-table-column align="left" label="登录账号" prop="username" width="120" />
            <el-table-column align="left" label="服务器标签名" prop="label" width="150" />
            <el-table-column align="left" label="服务器地区" prop="region" width="150" />
            <el-table-column align="left" label="状态" prop="status" width="100" />

        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            <el-button  type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
            <el-button  type="primary" link icon="edit" class="table-button" @click="updateSshServerFunc(scope.row)">编辑</el-button>
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
          </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看">
            <el-descriptions :column="1" border>
                    <el-descriptions-item label="服务器IP地址">
    {{ detailForm.ip }}
</el-descriptions-item>
                    <el-descriptions-item label="SSH端口">
    {{ detailForm.port }}
</el-descriptions-item>
                    <el-descriptions-item label="登录账号">
    {{ detailForm.username }}
</el-descriptions-item>
                    <el-descriptions-item label="服务器标签名">
    {{ detailForm.label }}
</el-descriptions-item>
                    <el-descriptions-item label="服务器地区">
    {{ detailForm.region }}
</el-descriptions-item>
                    <el-descriptions-item label="状态">
    {{ detailForm.status }}
</el-descriptions-item>
            </el-descriptions>
        </el-drawer>

  </div>
</template>

<script setup>
import {
  createSshServer,
  deleteSshServer,
  deleteSshServerByIds,
  updateSshServer,
  findSshServer,
  getSshServerList,
  getSshNvidiaSmiText
} from '@/api/ssh/sshServer'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"




defineOptions({
    name: 'SshServer'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
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
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
               port : [{
                   required: true,
                   message: '请输入端口',
                   trigger: ['input','blur'],
               },
              ],
               username : [{
                   required: true,
                   message: '请输入账号',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
               password : [{
                   required: true,
                   message: '请输入密码',
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
            ip: 'ip',
            status: 'status',
            label: 'label',
            region: 'region',
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
  const table = await getSshServerList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
    // 初始化扩展字段
    tableData.value.forEach(row => {
      row.__nvsmiText = ''
      row.__nvsmiLoading = false
      row.__nvsmiTail = 50
      row.__nvsmiAuto = false
      row.__nvsmiInterval = 5
      row.__nvsmiTimer = null
    })
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
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
            deleteSshServerFunc(row)
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
      const res = await deleteSshServerByIds({ IDs })
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
const updateSshServerFunc = async(row) => {
    const res = await findSshServer({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteSshServerFunc = async (row) => {
    const res = await deleteSshServer({ ID: row.ID })
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
        ip: '',
        port: 22,
        username: 'root',
        label: '',
        region: '',
        password: '',
        status: '',
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
  const res = await findSshServer({ ID: row.ID })
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


// 展开事件与自动刷新
const onExpandChange = (row, expandedRows) => {
  const expanded = expandedRows.some(r => r.ID === row.ID)
  if (expanded) {
    refreshNvsmi(row)
    setupAuto(row)
  } else {
    teardownAuto(row)
  }
}
const refreshNvsmi = async (row) => {
  row.__nvsmiLoading = true
  try {
    const res = await getSshNvidiaSmiText({ ID: row.ID, tail: row.__nvsmiTail })
    row.__nvsmiLoading = false
    if (res.code === 0) {
      row.__nvsmiText = (res.data && res.data.text) ? res.data.text : ''
    } else {
      ElMessage.error(res.msg || '获取失败')
    }
  } catch (e) {
    row.__nvsmiLoading = false
    ElMessage.error('获取失败')
  }
}
const setupAuto = (row) => {
  teardownAuto(row)
  if (row.__nvsmiAuto) {
    const interval = Math.max(2, Math.min(120, Number(row.__nvsmiInterval || 5)))
    row.__nvsmiTimer = setInterval(() => refreshNvsmi(row), interval * 1000)
  }
}
const teardownAuto = (row) => {
  if (row.__nvsmiTimer) { clearInterval(row.__nvsmiTimer); row.__nvsmiTimer = null }
}

</script>

<style scoped>
.logbox { white-space: pre-wrap; background: #111; color: #d5d5d5; padding: 8px; border-radius: 4px; font-size: 12px; line-height: 1.4; }
</style>

