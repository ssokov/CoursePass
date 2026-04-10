import { computed } from 'vue'
import { useLocalStorage } from '@vueuse/core'
import { api } from '../api/client'
import type { IStudentDraft } from '../api/api'

const token = useLocalStorage<string>('token', '')

export function useAuth() {
  const isAuthenticated = computed(() => Boolean(token.value))

  async function login(login: string, password: string) {
    const result = await api.auth.login({ studentLogin: { login, password } })
    token.value = result.accessToken
  }

  async function register(draft: IStudentDraft) {
    const result = await api.auth.registerStudent({ studentDraft: draft })
    token.value = result.accessToken
  }

  function logout() {
    token.value = ''
  }

  return { token, isAuthenticated, login, register, logout }
}
