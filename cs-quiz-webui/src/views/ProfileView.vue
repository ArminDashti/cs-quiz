<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { fetchProfile, resolveAssetUrl, type PublicProfile } from '@/lib/auth'
import { useAuth } from '@/lib/useAuth'

const route = useRoute()
const router = useRouter()
const { user } = useAuth()

const profile = ref<PublicProfile | null>(null)
const errorMessage = ref<string | null>(null)

const username = computed(() => {
  const param = route.params.username
  if (typeof param === 'string' && param) return param
  return user.value?.username ?? ''
})

const avatarUrl = computed(() => resolveAssetUrl(profile.value?.avatar_url))

async function load() {
  errorMessage.value = null
  profile.value = null
  const name = username.value
  if (!name) {
    errorMessage.value = 'Set a username on your account to view a public profile.'
    return
  }
  if (route.name === 'profile-self' && user.value?.username) {
    await router.replace(`/profile/${user.value.username}`)
    return
  }
  try {
    profile.value = await fetchProfile(name)
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Failed to load profile'
  }
}

onMounted(() => {
  void load()
})

watch(
  () => route.fullPath,
  () => {
    void load()
  },
)
</script>

<template>
  <div class="mx-auto max-w-2xl space-y-4 px-4 py-10">
    <Card>
      <CardHeader>
        <CardTitle class="tracking-tight">Profile</CardTitle>
        <CardDescription>Public scores for this user.</CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <p v-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>
        <template v-if="profile">
          <div class="flex items-center gap-3">
            <img
              v-if="avatarUrl"
              :src="avatarUrl"
              alt=""
              class="h-14 w-14 rounded-full object-cover ring-1 ring-border"
            />
            <div>
              <p class="text-lg font-semibold tracking-tight">@{{ profile.username }}</p>
              <p class="text-sm text-muted-foreground">Recent quiz scores</p>
            </div>
          </div>
          <ul class="space-y-2 text-sm">
            <li
              v-for="score in profile.scores"
              :key="score.id"
              class="flex justify-between rounded-md border border-border bg-card px-3 py-2.5"
            >
              <span class="font-medium">{{ score.quiz_slug ?? score.quiz_id }}</span>
              <span class="text-muted-foreground">{{ score.correct }} / {{ score.total }}</span>
            </li>
            <li v-if="profile.scores.length === 0" class="text-muted-foreground">
              No scores yet.
            </li>
          </ul>
        </template>
      </CardContent>
    </Card>
  </div>
</template>
