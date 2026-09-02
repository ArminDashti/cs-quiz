<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { Pencil, Plus, Trash2 } from '@lucide/vue'
import QuestionFormModal from '@/components/QuestionFormModal.vue'
import { Button } from '@/components/ui/button'
import { adminDeleteQuestion, adminListQuestions, type AdminQuestion } from '@/lib/auth'

const route = useRoute()
const quizName = computed(() => String(route.params.quizName ?? ''))
const questions = ref<AdminQuestion[]>([])
const errorMessage = ref<string | null>(null)
const selected = ref<Set<string>>(new Set())
const modalOpen = ref(false)
const editing = ref<AdminQuestion | null>(null)

const labels = ['A', 'B', 'C', 'D'] as const

const allSelected = computed(
  () => questions.value.length > 0 && selected.value.size === questions.value.length,
)

async function load() {
  errorMessage.value = null
  selected.value = new Set()
  try {
    questions.value = await adminListQuestions(quizName.value)
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Failed to load'
  }
}

onMounted(() => {
  void load()
})

watch(quizName, () => {
  void load()
})

function toggleAll() {
  if (allSelected.value) {
    selected.value = new Set()
    return
  }
  selected.value = new Set(questions.value.map((q) => q.id))
}

function toggleOne(id: string) {
  const next = new Set(selected.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  selected.value = next
}

function openCreate() {
  editing.value = null
  modalOpen.value = true
}

function openEdit(q: AdminQuestion) {
  editing.value = q
  modalOpen.value = true
}

function typeLabel(t: string): string {
  if (!t || t === 'multiple_choice') return 'Multiple choice'
  return t
}

function difficultyLabel(d: string): string {
  if (d === 'easy') return 'Easy'
  if (d === 'hard') return 'Hard'
  return 'Medium'
}

async function removeSelected() {
  if (selected.value.size === 0) return
  if (!confirm(`Delete ${selected.value.size} question(s)?`)) return
  try {
    for (const id of selected.value) {
      await adminDeleteQuestion(quizName.value, id)
    }
    await load()
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Delete failed'
  }
}
</script>

<template>
  <div class="mx-auto max-w-6xl space-y-4 px-4 py-10">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <h1 class="text-2xl font-semibold tracking-tight">Questions — {{ quizName }}</h1>
      <div class="flex flex-wrap gap-2">
        <Button
          variant="outline"
          size="sm"
          class="gap-1.5 text-red-600 dark:text-red-400"
          :disabled="selected.size === 0"
          @click="removeSelected"
        >
          <Trash2 :size="14" aria-hidden="true" />
          Delete selected
        </Button>
        <Button size="sm" class="gap-1.5" @click="openCreate">
          <Plus :size="14" aria-hidden="true" />
          Add question
        </Button>
      </div>
    </div>

    <p v-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>

    <div class="overflow-x-auto rounded-md border border-border">
      <table class="w-full min-w-[800px] border-collapse text-sm">
        <thead class="bg-muted/50 text-left">
          <tr>
            <th class="px-3 py-2 font-medium">
              <input
                type="checkbox"
                :checked="allSelected"
                aria-label="Select all"
                @change="toggleAll"
              />
            </th>
            <th class="px-3 py-2 font-medium">Number</th>
            <th class="px-3 py-2 font-medium">Questions</th>
            <th class="px-3 py-2 font-medium">Type</th>
            <th class="px-3 py-2 font-medium">Option</th>
            <th class="px-3 py-2 font-medium">Difficulty level</th>
            <th class="px-3 py-2 font-medium">Edit</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(q, index) in questions" :key="q.id" class="border-t border-border">
            <td class="px-3 py-2 align-middle">
              <input
                type="checkbox"
                :checked="selected.has(q.id)"
                :aria-label="`Select question ${index + 1}`"
                @change="toggleOne(q.id)"
              />
            </td>
            <td class="px-3 py-2 align-middle">{{ index + 1 }}</td>
            <td class="max-w-md px-3 py-2 align-middle">{{ q.prompt }}</td>
            <td class="px-3 py-2 align-middle whitespace-nowrap">{{ typeLabel(q.question_type) }}</td>
            <td class="px-3 py-2 align-middle">{{ labels[q.correct_index] ?? '—' }}</td>
            <td class="px-3 py-2 align-middle">{{ difficultyLabel(q.difficulty) }}</td>
            <td class="px-3 py-2 align-middle">
              <Button size="sm" variant="outline" class="gap-1.5" @click="openEdit(q)">
                <Pencil :size="14" aria-hidden="true" />
                Edit
              </Button>
            </td>
          </tr>
          <tr v-if="questions.length === 0">
            <td colspan="7" class="px-3 py-6 text-center text-muted-foreground">No questions yet.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <QuestionFormModal
      v-model:open="modalOpen"
      :quiz-name="quizName"
      :question="editing"
      @saved="load"
    />
  </div>
</template>
