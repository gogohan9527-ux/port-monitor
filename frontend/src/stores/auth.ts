import { defineStore } from 'pinia'
import http from '@/api/http'

const TOKEN_KEY = 'pm_token'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem(TOKEN_KEY) || ''
  }),
  actions: {
    async login(username: string, password: string) {
      const { data } = await http.post<{ token: string }>('/login', { username, password })
      this.token = data.token
      localStorage.setItem(TOKEN_KEY, data.token)
    },
    clear() {
      this.token = ''
      localStorage.removeItem(TOKEN_KEY)
    }
  }
})
