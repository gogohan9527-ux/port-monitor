<template>
  <div class="page">
    <div class="header">
      <div>
        <h1>端口管理</h1>
        <div class="subtitle">查看本地端口监控情况</div>
      </div>
      <el-button text @click="logout">
        <el-icon><SwitchButton /></el-icon>
        <span style="margin-left:4px">退出</span>
      </el-button>
    </div>

    <div class="stats">
      <StatCard title="监听总端口数" :value="ports.total" hint="本机正在使用的全部端口" :icon="Monitor" tone="blue" />
      <StatCard title="TCP 端口" :value="ports.tcpCount" hint="使用中的 TCP 端口" :icon="Connection" tone="green" />
      <StatCard title="UDP 端口" :value="ports.udpCount" hint="使用中的 UDP 端口" :icon="Promotion" tone="purple" />
      <StatCard title="高危端口" :value="ports.riskCount" :hint="refreshHint" :icon="Warning" tone="orange" />
    </div>

    <PortFilter v-model="filter" @open-settings="settingsOpen = true" />

    <PortTable :rows="filteredRows" />

    <SettingsDialog v-model="settingsOpen" />
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Monitor, Connection, Promotion, Warning, SwitchButton } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import StatCard from '@/components/StatCard.vue'
import PortFilter, { type FilterValue } from '@/components/PortFilter.vue'
import PortTable from '@/components/PortTable.vue'
import SettingsDialog from '@/components/SettingsDialog.vue'
import { useAuthStore } from '@/stores/auth'
import { usePortsStore, type PortInfo } from '@/stores/ports'
import { useSettingsStore } from '@/stores/settings'
import http from '@/api/http'
import { connectWS, disconnectWS } from '@/api/ws'

const router = useRouter()
const auth = useAuthStore()
const ports = usePortsStore()
const settings = useSettingsStore()

const filter = reactive<FilterValue>({ q: '', protocol: '', status: '', family: '', onlyRisk: false })
const settingsOpen = ref(false)

const refreshHint = computed(() => {
  const ms = settings.intervalMs || 1000
  return ms >= 1000 ? `1次/${(ms / 1000).toFixed(ms % 1000 ? 1 : 0)}s` : `1次/${ms}ms`
})

function parsePortRange(q: string): { min: number; max: number } | null {
  const m = q.match(/^\s*(\d{1,5})?\s*-\s*(\d{1,5})?\s*$/)
  if (!m) return null
  if (m[1] === undefined && m[2] === undefined) return null
  let min = m[1] === undefined ? 0 : Number(m[1])
  let max = m[2] === undefined ? 65535 : Number(m[2])
  if (min > 65535 || max > 65535) return null
  if (min > max) [min, max] = [max, min]
  return { min, max }
}

const filteredRows = computed<PortInfo[]>(() => {
  const qRaw = filter.q.trim()
  const q = qRaw.toLowerCase()
  const range = qRaw ? parsePortRange(qRaw) : null
  return ports.list.filter((p) => {
    if (filter.onlyRisk && !p.highRisk) return false
    if (filter.protocol && !p.protocol.startsWith(filter.protocol)) return false
    if (filter.status && p.status !== filter.status) return false
    if (filter.family && p.family !== filter.family) return false
    if (range) {
      if (p.port < range.min || p.port > range.max) return false
    } else if (q) {
      const hay = `${p.port} ${p.process} ${p.address}`.toLowerCase()
      if (!hay.includes(q)) return false
    }
    return true
  })
})

function logout() {
  disconnectWS()
  auth.clear()
  router.replace({ name: 'login' })
}

onMounted(async () => {
  try {
    await settings.load()
  } catch {
    /* will retry via WS */
  }
  try {
    const { data } = await http.get<{ ports: PortInfo[]; timestamp: number }>('/ports')
    ports.setList(data.ports, data.timestamp)
  } catch (e: any) {
    ElMessage.warning('初始数据加载失败，等待 WebSocket 推送')
  }
  if (auth.token) connectWS(auth.token)
})

onBeforeUnmount(() => {
  disconnectWS()
})
</script>

<style scoped>
.page {
  padding: 20px 28px;
  max-width: 1400px;
  margin: 0 auto;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 18px;
}
.header h1 {
  margin: 0;
  font-size: 22px;
  font-weight: 600;
}
.subtitle {
  color: #909399;
  font-size: 13px;
  margin-top: 4px;
}
.stats {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
</style>
