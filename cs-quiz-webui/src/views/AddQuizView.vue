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
import { adminCreateQuiz } from '@/lib/auth'

const router = useRouter()
const name = ref('')
const slug = ref('')
const description = ref('')
const errorMessage = ref<string | null>(null)
const submitting = ref(false)

async function onSubmit() {
  errorMessage.value = null
  submitting.value = true
  try {
    const quiz = await adminCreateQuiz({
      name: name.value.trim(),
      slug: slug.value.trim() || undefined,
      description: description.value.trim(),
    })
    await router.push(`/management/quiz/${quiz.slug}/add-question`)
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
        <CardTitle class="tracking-tight">Add quiz</CardTitle>
        <CardDescription>Create a new quiz topic.</CardDescription>
      </CardHeader>
      <CardContent>
        <form class="space-y-3" @submit.prevent="onSubmit">
          <label class="block space-y-1 text-sm">
            <span>Name</span>
            <input
              v-model="name"
              required
              class="w-full rounded-md border border-input bg-background px-3 py-2"
            />
          </label>
          <label class="block space-y-1 text-sm">
            <span>Slug (optional)</span>
            <input
              v-model="slug"
              placeholder="auto from name"
              class="w-full rounded-md border border-input bg-background px-3 py-2"
            />
          </label>
          <label class="block space-y-1 text-sm">
            <span>Description</span>
            <textarea
              v-model="description"
              rows="3"
              class="w-full rounded-md border border-input bg-background px-3 py-2"
            />
          </label>
          <p v-if="errorMessage" class="text-sm text-red-600">{{ errorMessage }}</p>
          <Button type="submit" :disabled="submitting">
            {{ submitting ? 'Creating…' : 'Create quiz' }}
          </Button>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
