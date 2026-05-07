<template>
  <el-dialog v-model="visible" title="设置" width="480px">
    <el-form label-width="120px" size="default">
      <el-form-item label="刷新间隔 (ms)">
        <el-input-number v-model="intervalMs" :min="200" :max="600000" :step="100" />
        <div class="hint">最小 200ms。后端立即生效。</div>
      </el-form-item>
      <el-form-item label="高危端口">
        <el-input
          v-model="riskText"
          type="textarea"
          :rows="4"
          placeholder="逗号或空格分隔，例如 22, 3389, 6379"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="saving" @click="save">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useSettingsStore } from '@/stores/settings'

const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{ (e: 'update:modelValue', v: boolean): void }>()

const settings = useSettingsStore()
const visible = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

const intervalMs = ref(settings.intervalMs)
const riskText = ref((settings.highRiskPorts || []).join(', '))
const saving = ref(false)

watch(visible, (v) => {
  if (v) {
    intervalMs.value = settings.intervalMs
    riskText.value = (settings.highRiskPorts || []).join(', ')
  }
})

async function save() {
  const ports = riskText.value
    .split(/[\s,，]+/)
    .map((s) => parseInt(s, 10))
    .filter((n) => Number.isFinite(n) && n > 0 && n < 65536)
  saving.value = true
  try {
    await settings.save(intervalMs.value, ports)
    ElMessage.success('设置已保存')
    visible.value = false
  } catch (e: any) {
    ElMessage.error(e?.response?.data || e?.message || '保存失败')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.hint { color: #909399; font-size: 12px; margin-top: 4px; }
</style>
