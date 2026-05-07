import { defineStore } from 'pinia'

export interface PortInfo {
  port: number
  protocol: string
  status: string
  family: string
  process: string
  pid: number
  address: string
  startedAt: string
  highRisk: boolean
}

export const usePortsStore = defineStore('ports', {
  state: () => ({
    list: [] as PortInfo[],
    timestamp: 0
  }),
  getters: {
    total: (s) => s.list.length,
    tcpCount: (s) => s.list.filter((p) => p.protocol.startsWith('TCP')).length,
    udpCount: (s) => s.list.filter((p) => p.protocol.startsWith('UDP')).length,
    riskCount: (s) => s.list.filter((p) => p.highRisk).length
  },
  actions: {
    setList(list: PortInfo[], ts?: number) {
      this.list = list
      if (ts) this.timestamp = ts
    }
  }
})
