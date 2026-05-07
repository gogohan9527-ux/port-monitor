<template>
  <div class="stat-card" :class="`tone-${tone}`">
    <div class="icon">
      <el-icon :size="22"><component :is="icon" /></el-icon>
    </div>
    <div class="body">
      <div class="title">{{ title }}</div>
      <div class="value">{{ value }}</div>
      <div class="hint">{{ hint }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
defineProps<{
  title: string
  value: number | string
  hint: string
  icon: Component
  tone: 'blue' | 'green' | 'purple' | 'orange'
}>()
</script>

<style scoped>
.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 18px 20px;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 10px;
  flex: 1;
}
.icon {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
}
.body { flex: 1; min-width: 0; }
.title { color: #606266; font-size: 13px; }
.value { font-size: 28px; font-weight: 600; line-height: 1.2; margin: 4px 0; color: #303133; }
.hint  { color: #909399; font-size: 12px; }

.tone-blue   .icon { background: #eff6ff; color: #2563eb; }
.tone-green  .icon { background: #ecfdf5; color: #10b981; }
.tone-purple .icon { background: #f5f3ff; color: #8b5cf6; }
.tone-orange .icon { background: #fff7ed; color: #f97316; }
</style>
