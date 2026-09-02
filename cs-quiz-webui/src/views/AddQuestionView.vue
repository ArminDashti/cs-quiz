<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { adminCreateQuestion } from '@/lib/auth'

const route = useRoute()
const router = useRouter()

const quizName = computed(() => String(route.params.quizName ?? ''))

const prompt = ref('')
const optionA = ref('')
const optionB = ref('')
const optionC = ref('')
const optionD = ref('')
const correctIndex = ref(0)
const errorMessage = ref<string | null>(null)
const submitting = ref(false)

async function onSubmit() {
  errorMessage.value = null
  submitting.value = true
  try {
    await adminCreateQuestion(quizName.value, {
      prompt: prompt.value.trim(),
      option_a: optionA.value.trim(),
      option_b: optionB.value.trim(),
      option_c: optionC.value.trim(),
      option_d: optionD.value.trim(),
      correct_index: correctIndex.value,
    })
    await router.push(`/management/quiz/${quizName.value}/questions-list`)
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Create failed'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-lg px-4 py-10">
    <Card>
      <CardHeader>
        <CardTitle class="tracking-tight">Add question</CardTitle>
        <CardDescription>Quiz: {{ quizName }}</CardDescription>
      </CardHeader>
      <CardContent>
        <form class="space-y-3" @submit.prevent="onSubmit">
          <label class="block space-y-1 text-sm">
            <span>Prompt</span>
            <textarea
              v-model="prompt"
              required
              rows="3"
              class="w-full rounded-md border border-input bg-background px-3 py-2"
            />
          </label>
          <label class="block space-y-1 text-sm">
            <span>Option A</span>
            <input v-model="optionA" required class="w-full rounded-md border border-input bg-background px-3 py-2" />
          </label>
          <label class="block space-y-1 text-sm">
            <span>Option B</span>
            <input v-model="optionB" required class="w-full rounded-md border border-input bg-background px-3 py-2" />
          </label>
          <label class="block space-y-1 text-sm">
            <span>Option C</span>
            <input v-model="optionC" required class="w-full rounded-md border border-input bg-background px-3 py-2" />
          </label>
          <label class="block space-y-1 text-sm">
            <span>Option D</span>
            <input v-model="optionD" required class="w-full rounded-md border border-input bg-background px-3 py-2" />
          </label>
          <label class="block space-y-1 text-sm">
            <span>Correct option</span>
            <select
              v-model.number="correctIndex"
              class="w-full rounded-md border border-input bg-background px-3 py-2"
            >
              <option :value="0">A</option>
              <option :value="1">B</option>
              <option :value="2">C</option>
              <option :value="3">D</option>
            </select>
          </label>
          <p v-if="errorMessage" class="text-sm text-red-600">{{ errorMessage }}</p>
          <Button type="submit" :disabled="submitting">
            {{ submitting ? 'Saving…' : 'Add question' }}
          </Button>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
