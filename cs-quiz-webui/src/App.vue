<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { Monitor, Moon, Sun } from '@lucide/vue'
import AppFooter from '@/components/AppFooter.vue'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/lib/useAuth'
import { resolveAssetUrl } from '@/lib/auth'
import { getTheme, setTheme, type Theme } from '@/lib/theme'

const router = useRouter()
const { user, isAuthenticated, isAdmin, logout } = useAuth()

const theme = ref<Theme>(getTheme())
const avatarUrl = computed(() => resolveAssetUrl(user.value?.avatar_url))

const navLinkClass =
  'rounded-md px-2.5 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-muted hover:text-foreground'
const navLinkActiveClass = 'bg-muted text-foreground'

function syncTheme() {
  theme.value = getTheme()
}

function onLogout() {
  logout()
  void router.push('/')
}

function onTheme(next: Theme) {
  theme.value = setTheme(next)
}

let themeObserver: MutationObserver | null = null

onMounted(() => {
  themeObserver = new MutationObserver(syncTheme)
  themeObserver.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['class'],
  })
})

onBeforeUnmount(() => {
  themeObserver?.disconnect()
  themeObserver = null
})
</script>

<template>
  <div class="flex h-full w-full flex-col bg-background text-foreground">
    <header
      class="sticky top-0 z-40 flex shrink-0 items-center justify-between gap-4 border-b border-border/80 bg-background/95 px-4 py-3.5 backdrop-blur supports-[backdrop-filter]:bg-background/80"
    >
      <nav class="flex min-w-0 flex-wrap items-center gap-1 sm:gap-1.5" aria-label="Main">
        <RouterLink
          class="mr-2 shrink-0 text-base font-semibold tracking-tight text-foreground"
          to="/"
        >
          CS Quiz
        </RouterLink>
        <RouterLink :class="navLinkClass" :active-class="navLinkActiveClass" to="/quiz/csharp">
          C#
        </RouterLink>
        <RouterLink :class="navLinkClass" :active-class="navLinkActiveClass" to="/quiz/dotnet">
          .NET
        </RouterLink>
        <RouterLink :class="navLinkClass" :active-class="navLinkActiveClass" to="/quiz/python">
          Python
        </RouterLink>
        <RouterLink :class="navLinkClass" :active-class="navLinkActiveClass" to="/quiz/js">
          JS
        </RouterLink>
        <RouterLink
          :class="navLinkClass"
          :active-class="navLinkActiveClass"
          to="/quiz/sql-server"
        >
          SQL Server
        </RouterLink>
        <RouterLink
          v-if="isAdmin"
          :class="navLinkClass"
          :active-class="navLinkActiveClass"
          to="/management"
        >
          Management
        </RouterLink>
        <RouterLink :class="navLinkClass" :active-class="navLinkActiveClass" to="/about-me">
          About Me
        </RouterLink>
      </nav>
      <div class="flex shrink-0 items-center gap-1.5">
        <div class="flex items-center gap-0.5 rounded-md border border-border p-0.5" role="group" aria-label="Theme">
          <Button
            variant="ghost"
            size="sm"
            class="h-8 w-8 px-0"
            :class="theme === 'light' ? 'bg-muted' : ''"
            aria-label="Light theme"
            :aria-pressed="theme === 'light'"
            @click="onTheme('light')"
          >
            <Sun aria-hidden="true" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            class="h-8 w-8 px-0"
            :class="theme === 'dark' ? 'bg-muted' : ''"
            aria-label="Dark theme"
            :aria-pressed="theme === 'dark'"
            @click="onTheme('dark')"
          >
            <Moon aria-hidden="true" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            class="h-8 w-8 px-0"
            :class="theme === 'system' ? 'bg-muted' : ''"
            aria-label="System theme"
            :aria-pressed="theme === 'system'"
            @click="onTheme('system')"
          >
            <Monitor aria-hidden="true" />
          </Button>
        </div>
        <template v-if="isAuthenticated">
          <RouterLink
            v-if="user?.username"
            :class="navLinkClass"
            :to="`/profile/${user.username}`"
          >
            Profile
          </RouterLink>
          <RouterLink v-else :class="navLinkClass" to="/profile">Profile</RouterLink>
          <RouterLink to="/account">
            <Button variant="outline" size="sm" class="gap-2">
              <img
                v-if="avatarUrl"
                :src="avatarUrl"
                alt=""
                class="h-5 w-5 rounded-full object-cover"
              />
              Account
            </Button>
          </RouterLink>
          <Button variant="ghost" size="sm" @click="onLogout">Log out</Button>
        </template>
        <template v-else>
          <RouterLink to="/login">
            <Button variant="outline" size="sm">Log in</Button>
          </RouterLink>
          <RouterLink to="/signup">
            <Button size="sm">Sign up</Button>
          </RouterLink>
        </template>
      </div>
    </header>
    <div class="min-h-0 flex-1 overflow-auto">
      <div class="flex min-h-full flex-col">
        <main class="flex-1">
          <RouterView />
        </main>
        <AppFooter />
      </div>
    </div>
  </div>
</template>
