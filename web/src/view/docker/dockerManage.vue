<template>
  <div class="docker-manage">
    <el-row :gutter="16">
      <el-col :span="2">
        <el-card shadow="never" class="left-card">
          <div class="left-header">正常服务器</div>
          <el-scrollbar height="calc(100vh - 220px)">
            <el-menu :default-active="activeId" @select="onSelectServer">
              <el-menu-item v-for="s in servers" :key="s.ID" :index="String(s.ID)">
                <div class="server-item">
                  <div class="server-title">{{ s.label }}</div>
                </div>
              </el-menu-item>
            </el-menu>
          </el-scrollbar>
        </el-card>
      </el-col>
      <el-col :span="22">
        <el-card shadow="never">
          <div class="right-header">
            容器列表 <span v-if="activeId">(Endpoint: {{ activeId }})</span>
            <div style="float:right; display:flex; gap:8px">
              <el-select v-model="scope" size="small" style="width:160px" @change="fetchContainers">
                <el-option label="启动的容器" value="running" />
                <el-option label="关闭的容器" value="exited" />
                <el-option label="所有的容器" value="all" />
              </el-select>
              <el-button type="primary" size="small" :disabled="!activeId" @click="openCreate">新建容器</el-button>
              <el-button size="small" :disabled="!activeId" @click="fetchContainers">刷新</el-button>
            </div>
          </div>
          <el-table v-loading="loading" :data="containers" size="small" height="calc(100vh - 260px)">
          <el-table-column label="ID" width="100">
            <template #default="{ row }">
              <el-popover placement="top" trigger="click" :width="480">
                <template #reference>
                  <span class="id-cell one-line">{{ row.ID }}</span>
                </template>
                <div class="id-full">{{ row.ID }}</div>
              </el-popover>
            </template>
          </el-table-column>
            <el-table-column prop="Names" label="名称" width="120" />
            <el-table-column prop="Image" label="镜像" min-width="150" />
            <el-table-column prop="Command" label="命令" min-width="180" />
            <el-table-column prop="Status" label="状态" width="120" />
            <el-table-column prop="Ports" label="端口" min-width="120" />
            <el-table-column prop="RunningFor" label="运行时长" width="100" />
            <el-table-column prop="CreatedAt" label="创建时间" width="140" />
          <el-table-column label="操作" width="180">
            <template #default="scope">
              <el-button size="small" @click="handleStart(scope.row)">启动</el-button>
              <el-button size="small" type="warning" @click="handleStop(scope.row)">停止</el-button>
              <el-button size="small" type="danger" @click="handleRemove(scope.row)">删除</el-button>
              <el-button size="small" type="success" @click="openTerminal(scope.row)">终端</el-button>
            </template>
          </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
    <el-dialog v-model="createDialog" title="新建容器" width="520px">
      <el-form :model="createForm" label-width="120px">
        <el-form-item label="镜像">
          <el-input v-model="createForm.image" placeholder="如：nginx:latest" />
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="createForm.name" placeholder="容器名称，可空" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialog=false">取消</el-button>
        <el-button type="primary" @click="submitCreate">创建</el-button>
      </template>
    </el-dialog>
    <el-dialog v-model="terminalDialog" title="容器终端" width="900px" @open="focusTerm" @close="cleanupTerminal">
      <div ref="termRef" class="terminal"></div>
      <template #footer>
        <el-select v-model="shellChoice" size="small" style="width:140px;margin-right:8px">
          <el-option label="/bin/sh" value="/bin/sh" />
          <el-option label="/bin/bash" value="/bin/bash" />
        </el-select>
        <el-select v-model="sizePreset" size="small" style="width:140px;margin-right:8px">
          <el-option label="自适应" value="fit" />
          <el-option label="80×24" value="80x24" />
          <el-option label="100×30" value="100x30" />
          <el-option label="120×36" value="120x36" />
        </el-select>
        <el-button size="small" @click="copySelection" style="margin-right:8px">复制选中</el-button>
        <el-button size="small" @click="clearTerm" style="margin-right:8px">清屏</el-button>
        <el-button size="small" @click="terminalDialog=false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
  </template>

<script setup>
import { ref, onMounted, nextTick, watch, onBeforeUnmount } from 'vue'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { WebLinksAddon } from 'xterm-addon-web-links'
import 'xterm/css/xterm.css'
import { ElMessage } from 'element-plus'
import { getDockerServers, getDockerPs, createContainer, startContainer, stopContainer, removeContainer } from '@/api/docker/docker'

const servers = ref([])
const activeId = ref('')
const containers = ref([])
const loading = ref(false)
const scope = ref('running')

const fetchServers = async () => {
  try {
    const res = await getDockerServers()
    if (res.code === 0) {
      servers.value = res.data || []
    } else {
      ElMessage.error(res.msg || '获取服务器失败')
    }
  } catch (e) {
    ElMessage.error('获取服务器失败')
  }
}

const fetchContainers = async () => {
  if (!activeId.value) { containers.value = []; return }
  loading.value = true
  try {
    const res = await getDockerPs({ ID: activeId.value, scope: scope.value })
    loading.value = false
    if (res.code === 0) {
      containers.value = res.data || []
    } else {
      containers.value = []
      ElMessage.error(res.msg || '获取容器失败')
    }
  } catch (e) {
    loading.value = false
    ElMessage.error('获取容器失败')
  }
}

const onSelectServer = (index) => {
  activeId.value = index
  fetchContainers()
}

const createDialog = ref(false)
const createForm = ref({ image: '', name: '' })
const openCreate = () => { createDialog.value = true; createForm.value = { image: '', name: '' } }
const submitCreate = async () => {
  if (!activeId.value) { ElMessage.error('未选择Endpoint'); return }
  try {
    const res = await createContainer({ endpointId: activeId.value, image: createForm.value.image, name: createForm.value.name })
    if (res.code === 0) { ElMessage.success('创建成功'); createDialog.value = false; fetchContainers() } else { ElMessage.error(res.msg || '创建失败') }
  } catch (e) { ElMessage.error('创建失败') }
}

const handleStart = async (row) => {
  try {
    const res = await startContainer({ endpointId: activeId.value, ID: row.ID })
    if (res.code === 0) { ElMessage.success('已启动'); fetchContainers() } else { ElMessage.error(res.msg || '启动失败') }
  } catch (e) { ElMessage.error('启动失败') }
}
const handleStop = async (row) => {
  try {
    const res = await stopContainer({ endpointId: activeId.value, ID: row.ID })
    if (res.code === 0) { ElMessage.success('已停止'); fetchContainers() } else { ElMessage.error(res.msg || '停止失败') }
  } catch (e) { ElMessage.error('停止失败') }
}
const handleRemove = async (row) => {
  try {
    const res = await removeContainer({ endpointId: activeId.value, ID: row.ID, force: true })
    if (res.code === 0) { ElMessage.success('已删除'); fetchContainers() } else { ElMessage.error(res.msg || '删除失败') }
  } catch (e) { ElMessage.error('删除失败') }
}

const terminalDialog = ref(false)
const termRef = ref(null)
let termSocket = null
let term = null
let fitAddon = null
let heartbeatTimer = null
let reconnectTimer = null
let reconnectAttempts = 0
let currentRow = null
const shellChoice = ref('/bin/sh')
const sizePreset = ref('fit')
const openTerminal = (row) => {
  if (!activeId.value) { ElMessage.error('未选择Endpoint'); return }
  currentRow = row
  terminalDialog.value = true
  nextTick(() => connectTerminal(row))
}
const focusTerm = () => { /* xterm自身处理焦点 */ }
const connectTerminal = (row) => {
  const proto = location.protocol === 'https:' ? 'wss' : 'ws'
  const base = (import.meta && import.meta.env && import.meta.env.VITE_BASE_API) ? import.meta.env.VITE_BASE_API : ''
  const prefix = base.endsWith('/') ? base : (base ? base + '/' : '/')
  const url = `${proto}://${location.host}${prefix}docker/execWs?endpointId=${activeId.value}&ID=${row.ID}&shell=${encodeURIComponent(shellChoice.value)}`
  const el = termRef.value
  if (!el) return
  term = new Terminal({
    cursorBlink: true,
    fontFamily: 'monospace',
    fontSize: 14,
    cols: 80,
    rows: 24,
    convertEol: true,
    theme: { background: '#111', foreground: '#eee' }
  })
  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.loadAddon(new WebLinksAddon())
  term.open(el)
  termSocket = new WebSocket(url)
  termSocket.binaryType = 'arraybuffer'
  termSocket.onopen = () => {
    fitAddon.fit()
    sendResize()
    term.focus()
    term.write('\r\n连接已建立\r\n')
    startHeartbeat()
    reconnectAttempts = 0
  }
  termSocket.onmessage = (ev) => {
    if (typeof ev.data === 'string') {
      term.write(ev.data)
    } else if (ev.data instanceof ArrayBuffer) {
      term.write(new Uint8Array(ev.data))
    } else {
      const reader = new FileReader()
      reader.onload = () => { term.write(reader.result) }
      reader.readAsText(ev.data)
    }
  }
  termSocket.onclose = () => { ElMessage.info('终端已关闭'); scheduleReconnect() }
  termSocket.onerror = () => { ElMessage.error('终端连接异常'); scheduleReconnect() }
  term.onData((data) => {
    if (termSocket && termSocket.readyState === WebSocket.OPEN) {
      termSocket.send(data)
    }
  })
  term.attachCustomKeyEventHandler((e) => {
    if (e.ctrlKey && e.key === 'c' && term.hasSelection()) {
      const text = term.getSelection()
      if (navigator.clipboard) navigator.clipboard.writeText(text)
      return false
    }
    return true
  })
  window.addEventListener('resize', onWindowResize)
}
const onWindowResize = () => { if (fitAddon) { fitAddon.fit(); sendResize() } }
const sendResize = () => {
  if (!term || !termSocket || termSocket.readyState !== WebSocket.OPEN) return
  const cols = term.cols || 80
  const rows = term.rows || 24
  termSocket.send(JSON.stringify({ type: 'resize', cols, rows }))
}
const cleanupTerminal = () => {
  window.removeEventListener('resize', onWindowResize)
  stopHeartbeat()
  stopReconnect()
  if (termSocket) { try { termSocket.close() } catch(e){} termSocket = null }
  if (term) { try { term.dispose() } catch(e){} term = null }
  fitAddon = null
}
onBeforeUnmount(() => cleanupTerminal())
watch(terminalDialog, (v) => { if (!v) cleanupTerminal() })

const copySelection = async () => { if (!term) return; const text = term.getSelection(); if (!text) { ElMessage.info('无选中内容'); return } try { if (navigator.clipboard) await navigator.clipboard.writeText(text); ElMessage.success('已复制') } catch(e){ ElMessage.error('复制失败') } }
const clearTerm = () => { if (term) term.clear() }
const applySizePreset = () => {
  if (!term) return
  if (sizePreset.value === 'fit') { if (fitAddon) { fitAddon.fit(); sendResize() } return }
  const parts = sizePreset.value.split('x')
  const c = parseInt(parts[0]) || 80
  const r = parseInt(parts[1]) || 24
  if (typeof term.resize === 'function') { term.resize(c, r) }
  sendResize()
}
watch(sizePreset, () => applySizePreset())
watch(shellChoice, () => { if (terminalDialog.value && currentRow) { reconnectNow() } })
const startHeartbeat = () => { stopHeartbeat(); heartbeatTimer = setInterval(() => { if (termSocket && termSocket.readyState === WebSocket.OPEN) { termSocket.send(JSON.stringify({ type: 'ping' })) } }, 20000) }
const stopHeartbeat = () => { if (heartbeatTimer) { clearInterval(heartbeatTimer); heartbeatTimer = null } }
const scheduleReconnect = () => { if (!terminalDialog.value) return; stopReconnect(); const delay = Math.min(1000 * Math.pow(2, reconnectAttempts || 0), 10000); reconnectAttempts = (reconnectAttempts || 0) + 1; reconnectTimer = setTimeout(() => { reconnectNow() }, delay) }
const stopReconnect = () => { if (reconnectTimer) { clearTimeout(reconnectTimer); reconnectTimer = null } }
const reconnectNow = () => { if (!currentRow) return; if (termSocket) { try { termSocket.close() } catch(e){} termSocket = null } connectTerminal(currentRow) }

onMounted(() => {
  fetchServers()
})
</script>

<style scoped>
.docker-manage { padding: 8px; }
.left-card { min-height: calc(100vh - 180px); }
.left-header { font-weight: 600; margin-bottom: 8px; }
.right-header { font-weight: 600; margin-bottom: 8px; }
.server-item { display: flex; flex-direction: column; }
.server-title { font-weight: 600; }
.server-sub { font-size: 12px; color: var(--el-text-color-secondary); }
.id-cell { display:inline-block; max-width: 90px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; cursor: pointer; }
.id-full { font-family: monospace; word-break: break-all; }
.terminal { height: 540px; overflow: auto; background: #111; color: #eee; padding: 8px; font-family: monospace; white-space: pre-wrap; }
.xterm { height: 100%; }
.xterm-viewport { background: #111; }
</style>
