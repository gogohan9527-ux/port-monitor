<template>
  <div class="filter-row">
    <el-input
      v-model="local.q"
      placeholder="搜索端口/进程（支持 8000-9000）"
      clearable
      style="width: 240px"
    >
      <template #prefix><el-icon><Search /></el-icon></template>
    </el-input>
    <el-select v-model="local.protocol" placeholder="协议" clearable style="width: 140px">
      <el-option label="全部" value="" />
      <el-option label="TCP" value="TCP" />
      <el-option label="UDP" value="UDP" />
    </el-select>
    <el-select v-model="local.status" placeholder="状态" clearable style="width: 140px">
      <el-option label="全部" value="" />
      <el-option label="监听" value="监听" />
      <el-option label="活动" value="活动" />
    </el-select>
    <el-select v-model="local.family" placeholder="类型" clearable style="width: 140px">
      <el-option label="全部" value="" />
      <el-option label="IPv4" value="IPv4" />
      <el-option label="IPv6" value="IPv6" />
    </el-select>

    <div class="spacer" />

    <span class="risk-label">仅显示高危端口</span>
    <el-switch v-model="local.onlyRisk" />

    <el-button :icon="Setting" circle plain @click="$emit('open-settings')" />
  </div>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'
import { Search, Setting } from '@element-plus/icons-vue'

export interface FilterValue {
  q: string
  protocol: string
  status: string
  family: string
  onlyRisk: boolean
}

const props = defineProps<{ modelValue: FilterValue }>()
const emit = defineEmits<{
  (e: 'update:modelValue', v: FilterValue): void
  (e: 'open-settings'): void
}>()

const local = reactive<FilterValue>({ ...props.modelValue })
watch(local, (v) => emit('update:modelValue', { ...v }), { deep: true })
watch(() => props.modelValue, (v) => Object.assign(local, v))
</script>

<style scoped>
.filter-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 10px;
  margin-bottom: 12px;
}
.spacer { flex: 1; }
.risk-label { color: #606266; font-size: 13px; }
</style>
