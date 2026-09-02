<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ChevronLeft, ChevronRight } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { playQuiz, submitQuiz, type PlayQuestion, type Quiz } from '@/lib/auth'

const route = useRoute()
const slug = computed(() => String(route.params.slug ?? ''))

const quiz = ref<Quiz | null>(null)
const questions = ref<PlayQuestion[]>([])
const answers = ref<Record<string, number>>({})
const currentIndex = ref(0)
const result = ref<{ correct: number; total: number } | null>(null)
const errorMessage = ref<string | null>(null)
const loading = ref(true)
const submitting = ref(false)

const current = computed(() => questions.value[currentIndex.value] ?? null)
const isFirst = computed(() => currentIndex.value <= 0)
const isLast = computed(() => currentIndex.value >= questions.value.length - 1)
const options = (q: PlayQuestion) => [q.option_a, q.option_b, q.option_c, q.option_d]

async function load() {
  loading.value = true
  errorMessage.value = null
  result.value = null
  answers.value = {}
  currentIndex.value = 0
  try {
    const data = await playQuiz(slug.value)
    quiz.value = data.quiz
    questions.value = data.questions
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Failed to load quiz'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  void load()
})

watch(slug, () => {
  void load()
})

function goPrev() {
  if (!isFirst.value) currentIndex.value -= 1
}

function goNext() {
  if (!isLast.value) currentIndex.value += 1
}

async function onSubmit() {
  errorMessage.value = null
  submitting.value = true
  try {
    const payload = questions.value.map((q) => ({
      question_id: q.id,
      selected_index: answers.value[q.id] ?? -1,
    }))
    const res = await submitQuiz(slug.value, payload)
    result.value = { correct: res.correct, total: res.total }
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Submit failed'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-2xl space-y-4 px-4 py-10">
    <div v-if="loading" class="text-muted-foreground">Loading…</div>
    <template v-else-if="quiz">
      <Card>
        <CardHeader>
          <CardTitle class="text-2xl tracking-tight">{{ quiz.name }}</CardTitle>
          <CardDescription>{{ quiz.description }}</CardDescription>
        </CardHeader>
      </Card>

      <template v-if="!result && current">
        <Card>
          <CardHeader>
            <CardDescription>
              Question {{ currentIndex + 1 }} of {{ questions.length }}
            </CardDescription>
            <CardTitle class="text-base leading-snug">
              {{ currentIndex + 1 }}. {{ current.prompt }}
            </CardTitle>
          </CardHeader>
          <CardContent class="space-y-2">
            <label
              v-for="(opt, i) in options(current)"
              :key="i"
              class="flex cursor-pointer items-center gap-2 rounded-md border border-border px-3 py-2.5 text-sm transition-colors hover:bg-accent"
              :class="answers[current.id] === i ? 'border-primary bg-accent' : ''"
            >
              <input v-model.number="answers[current.id]" type="radio" :name="current.id" :value="i" />
              <span>{{ opt }}</span>
            </label>
          </CardContent>
        </Card>

        <div class="flex items-center justify-between gap-3">
          <Button variant="outline" size="sm" class="gap-1.5" :disabled="isFirst" aria-label="Previous question" @click="goPrev">
            <ChevronLeft aria-hidden="true" />
            Back
          </Button>
          <Button
            v-if="!isLast"
            variant="outline"
            size="sm"
            class="gap-1.5"
            aria-label="Next question"
            @click="goNext"
          >
            Forward
            <ChevronRight aria-hidden="true" />
          </Button>
          <Button
            v-else
            size="sm"
            :disabled="submitting || questions.length === 0"
            @click="onSubmit"
          >
            {{ submitting ? 'Submitting…' : 'Submit answers' }}
          </Button>
        </div>
      </template>

      <p v-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>

      <Card v-if="result">
        <CardHeader>
          <CardTitle class="tracking-tight">Result</CardTitle>
          <CardDescription>
            You scored {{ result.correct }} / {{ result.total }}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Button variant="outline" @click="load">Try again</Button>
        </CardContent>
      </Card>
    </template>
    <p v-else-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>
  </div>
</template>
