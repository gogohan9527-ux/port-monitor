<template>
  <div class="table-wrap">
    <el-table :data="paged" stripe size="default" style="width: 100%">
      <el-table-column prop="port" label="端口号" width="100" sortable />
      <el-table-column prop="protocol" label="协议" width="100">
        <template #default="{ row }">
          <el-tag :type="row.protocol.startsWith('TCP') ? 'success' : 'warning'" effect="plain" size="small">
            {{ row.protocol }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.highRisk ? 'danger' : 'info'" effect="light" size="small">
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="family" label="类型" width="100" />
      <el-table-column prop="process" label="进程" min-width="160" show-overflow-tooltip />
      <el-table-column prop="pid" label="PID" width="100" />
      <el-table-column prop="address" label="地址" min-width="160" show-overflow-tooltip />
      <el-table-column prop="startedAt" label="启动时间" width="180" />
    </el-table>

    <div class="pager">
      <el-pagination
        v-model:current-page="page"
        :page-size="pageSize"
        :total="rows.length"
        layout="prev, pager, next, jumper, total"
        background
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { PortInfo } from '@/stores/ports'

const props = defineProps<{ rows: PortInfo[] }>()
const page = ref(1)
const pageSize = 15

const paged = computed(() => {
  const start = (page.value - 1) * pageSize
  return props.rows.slice(start, start + pageSize)
})

watch(() => props.rows.length, () => {
  const max = Math.max(1, Math.ceil(props.rows.length / pageSize))
  if (page.value > max) page.value = max
})
</script>

<style scoped>
.table-wrap {
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 10px;
  padding: 8px 12px 16px;
}
.pager {
  display: flex;
  justify-content: flex-end;
  padding-top: 12px;
}
</style>
