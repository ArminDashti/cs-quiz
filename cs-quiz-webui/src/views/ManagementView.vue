<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { Ban, CircleCheck, List, Pencil, Plus, Trash2 } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Dialog } from '@/components/ui/dialog'
import {
  adminDeleteQuiz,
  adminListQuizzes,
  adminUpdateQuiz,
  fetchInviteCode,
  updateInviteCode,
  type Quiz,
} from '@/lib/auth'

const quizzes = ref<Quiz[]>([])
const inviteCode = ref('')
const message = ref<string | null>(null)
const errorMessage = ref<string | null>(null)

const editOpen = ref(false)
const editing = ref<Quiz | null>(null)
const editName = ref('')
const editCategory = ref('')
const editDescription = ref('')
const editAttachments = ref('')
const editSaving = ref(false)

function parseAttachments(raw: string): string[] {
  return raw
    .split(/\r?\n|,/)
    .map((s) => s.trim())
    .filter(Boolean)
}

function formatAttachments(list: string[] | undefined): string {
  return (list ?? []).join('\n')
}

async function load() {
  errorMessage.value = null
  try {
    quizzes.value = await adminListQuizzes()
    const invite = await fetchInviteCode()
    inviteCode.value = invite.invite_code
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Failed to load'
  }
}

onMounted(() => {
  void load()
})

async function saveInvite() {
  message.value = null
  errorMessage.value = null
  try {
    const result = await updateInviteCode(inviteCode.value.trim())
    inviteCode.value = result.invite_code
    message.value = 'Invite code updated.'
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Update failed'
  }
}

function openEdit(quiz: Quiz) {
  editing.value = quiz
  editName.value = quiz.name
  editCategory.value = quiz.category || ''
  editDescription.value = quiz.description
  editAttachments.value = formatAttachments(quiz.attachments)
  editOpen.value = true
}

async function saveEdit() {
  if (!editing.value) return
  editSaving.value = true
  errorMessage.value = null
  try {
    await adminUpdateQuiz(editing.value.slug, {
      name: editName.value.trim(),
      category: editCategory.value.trim(),
      description: editDescription.value.trim(),
      attachments: parseAttachments(editAttachments.value),
    })
    editOpen.value = false
    await load()
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Update failed'
  } finally {
    editSaving.value = false
  }
}

async function removeQuiz(slug: string) {
  if (!confirm(`Delete quiz "${slug}"?`)) return
  try {
    await adminDeleteQuiz(slug)
    await load()
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Delete failed'
  }
}

async function toggleEnabled(quiz: Quiz) {
  errorMessage.value = null
  try {
    await adminUpdateQuiz(quiz.slug, { enabled: !quiz.enabled })
    await load()
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Update failed'
  }
}

function formatDate(iso: string): string {
  try {
    return new Date(iso).toLocaleString()
  } catch {
    return iso
  }
}
</script>

<template>
  <div class="mx-auto max-w-6xl space-y-4 px-4 py-10">
    <div class="flex items-center justify-between gap-3">
      <h1 class="text-2xl font-semibold tracking-tight">Management</h1>
      <RouterLink to="/management/add-quiz">
        <Button class="gap-1.5">
          <Plus aria-hidden="true" />
          Add quiz
        </Button>
      </RouterLink>
    </div>

    <Card>
      <CardHeader>
        <CardTitle class="tracking-tight">Invite code</CardTitle>
        <CardDescription>Required for new signups.</CardDescription>
      </CardHeader>
      <CardContent class="flex flex-wrap gap-2">
        <input
          v-model="inviteCode"
          class="min-w-[12rem] flex-1 rounded-md border border-input bg-background px-3 py-2"
        />
        <Button variant="outline" @click="saveInvite">Save</Button>
      </CardContent>
    </Card>

    <p v-if="message" class="text-sm text-green-700 dark:text-green-400">{{ message }}</p>
    <p v-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">{{ errorMessage }}</p>

    <div class="overflow-x-auto rounded-md border border-border">
      <table class="w-full min-w-[960px] border-collapse text-sm">
        <thead class="bg-muted/50 text-left">
          <tr>
            <th class="px-3 py-2 font-medium">Title</th>
            <th class="px-3 py-2 font-medium">Category</th>
            <th class="px-3 py-2 font-medium">Description</th>
            <th class="px-3 py-2 font-medium">Attachments</th>
            <th class="px-3 py-2 font-medium">Created at</th>
            <th class="px-3 py-2 font-medium">Edit</th>
            <th class="px-3 py-2 font-medium">Questions</th>
            <th class="px-3 py-2 font-medium">Delete</th>
            <th class="px-3 py-2 font-medium">Disable</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="quiz in quizzes" :key="quiz.id" class="border-t border-border">
            <td class="px-3 py-2 align-middle font-medium">
              {{ quiz.name }}
              <span v-if="!quiz.enabled" class="ml-1 text-xs text-muted-foreground">(disabled)</span>
            </td>
            <td class="px-3 py-2 align-middle text-muted-foreground">
              {{ quiz.category || '—' }}
            </td>
            <td class="max-w-xs truncate px-3 py-2 align-middle text-muted-foreground">
              {{ quiz.description || '—' }}
            </td>
            <td class="max-w-[12rem] truncate px-3 py-2 align-middle text-muted-foreground">
              <template v-if="quiz.attachments?.length">
                <a
                  v-for="(url, i) in quiz.attachments"
                  :key="`${quiz.id}-${i}`"
                  :href="url"
                  class="mr-1 text-primary underline"
                  target="_blank"
                  rel="noreferrer"
                  >{{ i + 1 }}</a
                >
              </template>
              <span v-else>—</span>
            </td>
            <td class="whitespace-nowrap px-3 py-2 align-middle text-muted-foreground">
              {{ formatDate(quiz.created_at) }}
            </td>
            <td class="px-3 py-2 align-middle">
              <Button size="sm" variant="outline" class="gap-1.5" @click="openEdit(quiz)">
                <Pencil :size="14" aria-hidden="true" />
                Edit
              </Button>
            </td>
            <td class="px-3 py-2 align-middle">
              <RouterLink :to="`/management/quiz/${quiz.slug}/questions-list`">
                <Button size="sm" variant="outline" class="gap-1.5">
                  <List :size="14" aria-hidden="true" />
                  Questions
                </Button>
              </RouterLink>
            </td>
            <td class="px-3 py-2 align-middle">
              <Button
                size="sm"
                variant="ghost"
                class="gap-1.5 text-red-600 dark:text-red-400"
                @click="removeQuiz(quiz.slug)"
              >
                <Trash2 :size="14" aria-hidden="true" />
                Delete
              </Button>
            </td>
            <td class="px-3 py-2 align-middle">
              <Button size="sm" variant="outline" class="gap-1.5" @click="toggleEnabled(quiz)">
                <Ban v-if="quiz.enabled" :size="14" aria-hidden="true" />
                <CircleCheck v-else :size="14" aria-hidden="true" />
                {{ quiz.enabled ? 'Disable' : 'Enable' }}
              </Button>
            </td>
          </tr>
          <tr v-if="quizzes.length === 0">
            <td colspan="9" class="px-3 py-6 text-center text-muted-foreground">No quizzes yet.</td>
          </tr>
        </tbody>
      </table>
    </div>

    <Dialog
      :open="editOpen"
      title="Edit quiz"
      :description="editing ? `Slug: ${editing.slug}` : undefined"
      @close="editOpen = false"
    >
      <form class="space-y-3" @submit.prevent="saveEdit">
        <label class="block space-y-1 text-sm">
          <span>Title</span>
          <input
            v-model="editName"
            required
            class="w-full rounded-md border border-input bg-background px-3 py-2"
          />
        </label>
        <label class="block space-y-1 text-sm">
          <span>Category</span>
          <input
            v-model="editCategory"
            class="w-full rounded-md border border-input bg-background px-3 py-2"
          />
        </label>
        <label class="block space-y-1 text-sm">
          <span>Description</span>
          <textarea
            v-model="editDescription"
            rows="3"
            class="w-full rounded-md border border-input bg-background px-3 py-2"
          />
        </label>
        <label class="block space-y-1 text-sm">
          <span>Attachments (one URL per line)</span>
          <textarea
            v-model="editAttachments"
            rows="3"
            class="w-full rounded-md border border-input bg-background px-3 py-2"
            placeholder="https://example.com/file.pdf"
          />
        </label>
        <div class="flex justify-end gap-2">
          <Button type="button" variant="outline" @click="editOpen = false">Cancel</Button>
          <Button type="submit" :disabled="editSaving">
            {{ editSaving ? 'Saving…' : 'Save' }}
          </Button>
        </div>
      </form>
    </Dialog>
  </div>
</template>
