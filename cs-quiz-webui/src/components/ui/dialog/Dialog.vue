<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  open: boolean
  title?: string
  description?: string
  class?: HTMLAttributes['class']
}>()

const emit = defineEmits<{
  close: []
}>()

function onBackdrop() {
  emit('close')
}
</script>

<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      role="dialog"
      aria-modal="true"
    >
      <div class="absolute inset-0 bg-black/50" @click="onBackdrop" />
      <div
        :class="
          cn(
            'relative z-10 w-full max-w-lg max-h-[90vh] overflow-y-auto rounded-lg border border-border bg-background p-5 shadow-lg',
            props.class,
          )
        "
      >
        <div v-if="title || description" class="mb-4 space-y-1">
          <h2 v-if="title" class="text-lg font-semibold tracking-tight">{{ title }}</h2>
          <p v-if="description" class="text-sm text-muted-foreground">{{ description }}</p>
        </div>
        <slot />
      </div>
    </div>
  </Teleport>
</template>
