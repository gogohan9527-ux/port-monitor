import { defineStore } from 'pinia'
import http from '@/api/http'

export interface Settings {
  intervalMs: number
  highRiskPorts: number[]
  listenAddr?: string
}

export const useSettingsStore = defineStore('settings', {
  state: (): Settings => ({
    intervalMs: 1000,
    highRiskPorts: []
  }),
  actions: {
    async load() {
      const { data } = await http.get<Settings>('/settings')
      this.intervalMs = data.intervalMs
      this.highRiskPorts = data.highRiskPorts ?? []
    },
    async save(intervalMs: number, highRiskPorts: number[]) {
      const { data } = await http.put<Settings>('/settings', { intervalMs, highRiskPorts })
      this.intervalMs = data.intervalMs
      this.highRiskPorts = data.highRiskPorts ?? []
    }
  }
})
