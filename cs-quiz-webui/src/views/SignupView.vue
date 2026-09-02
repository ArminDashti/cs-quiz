<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { register } from '@/lib/auth'
import { useAuth } from '@/lib/useAuth'

const router = useRouter()
const { applySession } = useAuth()

const firstName = ref('')
const lastName = ref('')
const email = ref('')
const password = ref('')
const inviteCode = ref('')
const errorMessage = ref<string | null>(null)
const submitting = ref(false)

async function onSubmit() {
  errorMessage.value = null
  submitting.value = true
  try {
    const result = await register({
      first_name: firstName.value.trim(),
      last_name: lastName.value.trim(),
      email: email.value.trim(),
      password: password.value,
      invite_code: inviteCode.value.trim(),
    })
    applySession(result.token, result.user)
    await router.push('/account')
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Sign up failed'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="mx-auto flex min-h-full max-w-md items-center px-4 py-14">
    <Card class="w-full">
      <CardHeader>
        <CardTitle class="tracking-tight">Create account</CardTitle>
        <CardDescription>Names, email, password, and invite code required.</CardDescription>
      </CardHeader>
      <CardContent>
        <form class="space-y-3" @submit.prevent="onSubmit">
          <div class="grid grid-cols-2 gap-3">
            <label class="space-y-1 text-sm">
              <span>First name</span>
              <input
                v-model="firstName"
                required
                class="w-full rounded-md border border-input bg-background px-3 py-2"
              />
            </label>
            <label class="space-y-1 text-sm">
              <span>Last name</span>
              <input
                v-model="lastName"
                required
                class="w-full rounded-md border border-input bg-background px-3 py-2"
              />
            </label>
          </div>
          <label class="block space-y-1 text-sm">
            <span>Email</span>
            <input
              v-model="email"
              type="email"
              required
              autocomplete="email"
              class="w-full rounded-md border border-input bg-background px-3 py-2"
            />
          </label>
          <label class="block space-y-1 text-sm">
            <span>Password (min 8)</span>
            <input
              v-model="password"
              type="password"
              required
              minlength="8"
              autocomplete="new-password"
              class="w-full rounded-md border border-input bg-background px-3 py-2"
            />
          </label>
          <label class="block space-y-1 text-sm">
            <span>Invite code</span>
            <input
              v-model="inviteCode"
              required
              class="w-full rounded-md border border-input bg-background px-3 py-2"
            />
          </label>
          <p v-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">
            {{ errorMessage }}
          </p>
          <Button class="w-full" type="submit" :disabled="submitting">
            {{ submitting ? 'Creating…' : 'Sign up' }}
          </Button>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
