<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Monitor, Moon, Sun } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  changePassword,
  deleteAccount,
  resolveAssetUrl,
  updateAccount,
  uploadAvatar,
} from '@/lib/auth'
import { getTheme, setTheme, type Theme } from '@/lib/theme'
import { useAuth } from '@/lib/useAuth'

const router = useRouter()
const { user, setUser, logout } = useAuth()

const firstName = ref('')
const lastName = ref('')
const username = ref('')
const theme = ref<Theme>(getTheme())
const currentPassword = ref('')
const newPassword = ref('')
const message = ref<string | null>(null)
const errorMessage = ref<string | null>(null)
const saving = ref(false)

const avatarUrl = computed(() => resolveAssetUrl(user.value?.avatar_url))

onMounted(() => {
  if (!user.value) return
  firstName.value = user.value.first_name
  lastName.value = user.value.last_name
  username.value = user.value.username ?? ''
})

async function saveProfile() {
  message.value = null
  errorMessage.value = null
  saving.value = true
  try {
    const updated = await updateAccount({
      first_name: firstName.value.trim(),
      last_name: lastName.value.trim(),
      username: username.value.trim(),
    })
    setUser(updated)
    message.value = 'Profile saved.'
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Save failed'
  } finally {
    saving.value = false
  }
}

async function onAvatarChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  errorMessage.value = null
  try {
    const updated = await uploadAvatar(file)
    setUser(updated)
    message.value = 'Avatar updated.'
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Avatar upload failed'
  } finally {
    input.value = ''
  }
}

function onThemeChange(next: Theme) {
  theme.value = setTheme(next)
}

async function onChangePassword() {
  message.value = null
  errorMessage.value = null
  try {
    await changePassword({
      current_password: currentPassword.value,
      new_password: newPassword.value,
    })
    currentPassword.value = ''
    newPassword.value = ''
    message.value = 'Password changed.'
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Password change failed'
  }
}

async function onDeleteAccount() {
  if (!confirm('Delete your account permanently?')) return
  try {
    await deleteAccount()
    logout()
    await router.push('/')
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Delete failed'
  }
}
</script>

<template>
  <div class="mx-auto max-w-xl space-y-4 px-4 py-10">
    <Card>
      <CardHeader>
        <CardTitle class="tracking-tight">Account</CardTitle>
        <CardDescription>Update your profile, avatar, and password.</CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="flex items-center gap-3">
          <img
            v-if="avatarUrl"
            :src="avatarUrl"
            alt=""
            class="h-14 w-14 rounded-full object-cover ring-1 ring-border"
          />
          <div
            v-else
            class="flex h-14 w-14 items-center justify-center rounded-full bg-muted text-sm text-muted-foreground"
          >
            —
          </div>
          <label class="text-sm">
            <span class="mb-1 block text-muted-foreground">Avatar</span>
            <input type="file" accept="image/*" @change="onAvatarChange" />
          </label>
        </div>
        <label class="block space-y-1 text-sm">
          <span>First name</span>
          <input v-model="firstName" class="w-full rounded-md border border-input bg-background px-3 py-2" />
        </label>
        <label class="block space-y-1 text-sm">
          <span>Last name</span>
          <input v-model="lastName" class="w-full rounded-md border border-input bg-background px-3 py-2" />
        </label>
        <label class="block space-y-1 text-sm">
          <span>Username</span>
          <input v-model="username" class="w-full rounded-md border border-input bg-background px-3 py-2" />
        </label>
        <div class="space-y-1 text-sm">
          <span>Theme</span>
          <div class="flex gap-2">
            <Button
              size="sm"
              class="gap-1.5"
              :variant="theme === 'light' ? 'default' : 'outline'"
              aria-label="Light theme"
              @click="onThemeChange('light')"
            >
              <Sun :size="14" aria-hidden="true" />
              Light
            </Button>
            <Button
              size="sm"
              class="gap-1.5"
              :variant="theme === 'dark' ? 'default' : 'outline'"
              aria-label="Dark theme"
              @click="onThemeChange('dark')"
            >
              <Moon :size="14" aria-hidden="true" />
              Dark
            </Button>
            <Button
              size="sm"
              class="gap-1.5"
              :variant="theme === 'system' ? 'default' : 'outline'"
              aria-label="System theme"
              @click="onThemeChange('system')"
            >
              <Monitor :size="14" aria-hidden="true" />
              System
            </Button>
          </div>
        </div>
        <Button :disabled="saving" @click="saveProfile">
          {{ saving ? 'Saving…' : 'Save profile' }}
        </Button>

        <div class="space-y-3 border-t border-border pt-4">
          <h3 class="font-medium tracking-tight">Change password</h3>
          <input
            v-model="currentPassword"
            type="password"
            placeholder="Current password"
            class="w-full rounded-md border border-input bg-background px-3 py-2"
          />
          <input
            v-model="newPassword"
            type="password"
            placeholder="New password"
            class="w-full rounded-md border border-input bg-background px-3 py-2"
          />
          <Button variant="outline" @click="onChangePassword">Update password</Button>
        </div>

        <p v-if="message" class="text-sm text-green-700 dark:text-green-400">{{ message }}</p>
        <p v-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>

        <Button variant="ghost" class="text-red-600 dark:text-red-400" @click="onDeleteAccount">
          Delete account
        </Button>
      </CardContent>
    </Card>
  </div>
</template>
