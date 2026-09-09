<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { Play, Settings } from '@lucide/vue'
import HeroBackdrop from '@/components/HeroBackdrop.vue'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { fetchQuizzes, type Quiz } from '@/lib/auth'
import { useAuth } from '@/lib/useAuth'

const { isAuthenticated, isAdmin } = useAuth()
const quizzes = ref<Quiz[]>([])
const errorMessage = ref<string | null>(null)

onMounted(async () => {
  try {
    quizzes.value = await fetchQuizzes()
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Failed to load quizzes'
  }
})
</script>

<template>
  <div class="home-page">
    <section class="hero relative isolate min-h-[calc(100dvh-8.5rem)] overflow-hidden">
      <div class="pointer-events-none absolute inset-0 opacity-70 dark:opacity-55">
        <HeroBackdrop />
      </div>
      <div
        class="pointer-events-none absolute inset-0 bg-gradient-to-b from-background/75 via-background/55 to-background"
      />
      <div
        class="hero-content relative z-10 mx-auto flex min-h-[calc(100dvh-8.5rem)] max-w-3xl flex-col items-center justify-center px-4 py-16 text-center"
      >
        <p class="hero-kicker mb-3 text-sm font-medium uppercase tracking-[0.2em] text-primary">
          Computer science
        </p>
        <h1 class="hero-title text-5xl font-semibold tracking-tight sm:text-6xl">CS Quiz</h1>
        <p class="hero-lead mt-4 max-w-xl text-base text-muted-foreground sm:text-lg">
          Practice C#, .NET, Python, JavaScript, and SQL Server with multiple-choice quizzes.
        </p>
        <div class="hero-actions mt-8 flex flex-wrap items-center justify-center gap-3">
          <template v-if="isAuthenticated">
            <RouterLink to="/quiz/csharp">
              <Button size="lg">Start C# quiz</Button>
            </RouterLink>
            <RouterLink to="/quiz/js">
              <Button size="lg" variant="outline">Try JavaScript</Button>
            </RouterLink>
          </template>
          <template v-else>
            <RouterLink to="/signup">
              <Button size="lg">Sign up</Button>
            </RouterLink>
            <RouterLink to="/login">
              <Button size="lg" variant="outline">Log in</Button>
            </RouterLink>
          </template>
        </div>
      </div>
    </section>

    <section class="mx-auto max-w-4xl space-y-6 px-4 pb-14 pt-4">
      <div>
        <h2 class="text-xl font-semibold tracking-tight">Ways to play</h2>
        <p class="mt-1 text-sm text-muted-foreground">
          Pick a topic and submit answers for a scored attempt.
        </p>
      </div>

      <p v-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>

      <div class="grid gap-4 sm:grid-cols-2">
        <Card v-for="quiz in quizzes" :key="quiz.id">
          <CardHeader>
            <CardTitle>{{ quiz.name }}</CardTitle>
            <CardDescription>
              <span v-if="quiz.category" class="mb-1 block text-xs font-medium uppercase tracking-wide text-primary">
                {{ quiz.category }}
              </span>
              {{ quiz.description }}
            </CardDescription>
          </CardHeader>
          <CardContent class="space-y-3">
            <p class="text-xs text-muted-foreground">
              Created at {{ new Date(quiz.created_at).toLocaleString() }}
            </p>
            <div v-if="quiz.attachments?.length" class="flex flex-wrap gap-2 text-sm">
              <a
                v-for="(url, i) in quiz.attachments"
                :key="`${quiz.id}-att-${i}`"
                :href="url"
                target="_blank"
                rel="noreferrer"
                class="text-primary underline"
              >
                Attachment {{ i + 1 }}
              </a>
            </div>
            <RouterLink :to="`/quiz/${quiz.slug}`">
              <Button class="gap-1.5" :disabled="!isAuthenticated">
                <Play :size="14" aria-hidden="true" />
                {{ isAuthenticated ? 'Start quiz' : 'Log in to play' }}
              </Button>
            </RouterLink>
          </CardContent>
        </Card>
      </div>

      <div v-if="isAdmin" class="pt-2">
        <RouterLink to="/management">
          <Button variant="outline" class="gap-1.5">
            <Settings :size="14" aria-hidden="true" />
            Open management
          </Button>
        </RouterLink>
      </div>
    </section>
  </div>
</template>

<style scoped>
.hero-kicker,
.hero-title,
.hero-lead,
.hero-actions {
  animation: hero-rise 0.6s ease-out both;
}

.hero-title {
  animation-delay: 0.08s;
}

.hero-lead {
  animation-delay: 0.16s;
}

.hero-actions {
  animation-delay: 0.24s;
}

@keyframes hero-rise {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (prefers-reduced-motion: reduce) {
  .hero-kicker,
  .hero-title,
  .hero-lead,
  .hero-actions {
    animation: none;
  }
}
</style>
