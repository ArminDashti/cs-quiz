import { computed, ref } from 'vue'
import {
  clearSession,
  fetchMe,
  getStoredUser,
  getToken,
  setSession,
  type User,
} from '@/lib/auth'

const user = ref<User | null>(getStoredUser())
const ready = ref(false)

export function useAuth() {
  const isAuthenticated = computed(() => Boolean(user.value && getToken()))
  const isAdmin = computed(() => Boolean(user.value?.is_admin))

  async function hydrate() {
    if (!getToken()) {
      user.value = null
      ready.value = true
      return
    }
    try {
      const me = await fetchMe()
      user.value = me
      setSession(getToken()!, me)
    } catch {
      clearSession()
      user.value = null
    } finally {
      ready.value = true
    }
  }

  function applySession(token: string, next: User) {
    setSession(token, next)
    user.value = next
  }

  function logout() {
    clearSession()
    user.value = null
  }

  function setUser(next: User) {
    user.value = next
    const token = getToken()
    if (token) setSession(token, next)
  }

  return {
    user,
    ready,
    isAuthenticated,
    isAdmin,
    hydrate,
    applySession,
    logout,
    setUser,
  }
}
