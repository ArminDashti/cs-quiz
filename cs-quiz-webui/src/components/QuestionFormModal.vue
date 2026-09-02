<script setup lang="ts">
import { ref, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import {
  adminCreateQuestion,
  adminUpdateQuestion,
  type AdminQuestion,
  type QuestionPayload,
} from '@/lib/auth'

const props = defineProps<{
  open: boolean
  quizName: string
  question?: AdminQuestion | null
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  saved: []
}>()

const prompt = ref('')
const optionA = ref('')
const optionB = ref('')
const optionC = ref('')
const optionD = ref('')
const correctIndex = ref(0)
const questionType = ref('multiple_choice')
const difficulty = ref<'easy' | 'medium' | 'hard'>('medium')
const errorMessage = ref<string | null>(null)
const submitting = ref(false)

const isEdit = () => Boolean(props.question?.id)

function resetFromQuestion() {
  const q = props.question
  if (q) {
    prompt.value = q.prompt
    optionA.value = q.option_a
    optionB.value = q.option_b
    optionC.value = q.option_c
    optionD.value = q.option_d
    correctIndex.value = q.correct_index
    questionType.value = q.question_type || 'multiple_choice'
    difficulty.value =
      q.difficulty === 'easy' || q.difficulty === 'hard' ? q.difficulty : 'medium'
  } else {
    prompt.value = ''
    optionA.value = ''
    optionB.value = ''
    optionC.value = ''
    optionD.value = ''
    correctIndex.value = 0
    questionType.value = 'multiple_choice'
    difficulty.value = 'medium'
  }
  errorMessage.value = null
}

watch(
  () => [props.open, props.question] as const,
  ([open]) => {
    if (open) resetFromQuestion()
  },
)

function close() {
  emit('update:open', false)
}

async function onSubmit() {
  errorMessage.value = null
  submitting.value = true
  const body: QuestionPayload = {
    prompt: prompt.value.trim(),
    option_a: optionA.value.trim(),
    option_b: optionB.value.trim(),
    option_c: optionC.value.trim(),
    option_d: optionD.value.trim(),
    correct_index: correctIndex.value,
    question_type: questionType.value.trim() || 'multiple_choice',
    difficulty: difficulty.value,
  }
  try {
    if (isEdit() && props.question) {
      await adminUpdateQuestion(props.quizName, props.question.id, body)
    } else {
      await adminCreateQuestion(props.quizName, body)
    }
    emit('saved')
    close()
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Save failed'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <Dialog
    :open="open"
    :title="isEdit() ? 'Edit question' : 'Add question'"
    :description="`Quiz: ${quizName}`"
    @close="close"
  >
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
      <label class="block space-y-1 text-sm">
        <span>Type</span>
        <input
          v-model="questionType"
          class="w-full rounded-md border border-input bg-background px-3 py-2"
        />
      </label>
      <label class="block space-y-1 text-sm">
        <span>Difficulty level</span>
        <select v-model="difficulty" class="w-full rounded-md border border-input bg-background px-3 py-2">
          <option value="easy">Easy</option>
          <option value="medium">Medium</option>
          <option value="hard">Hard</option>
        </select>
      </label>
      <p v-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>
      <div class="flex justify-end gap-2 pt-1">
        <Button type="button" variant="outline" @click="close">Cancel</Button>
        <Button type="submit" :disabled="submitting">
          {{ submitting ? 'Saving…' : isEdit() ? 'Save' : 'Add question' }}
        </Button>
      </div>
    </form>
  </Dialog>
</template>
