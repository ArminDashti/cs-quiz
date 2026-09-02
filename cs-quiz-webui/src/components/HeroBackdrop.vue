<template>
  <div class="hero-backdrop absolute inset-0 overflow-hidden" aria-hidden="true">
    <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_hsl(var(--primary)/0.18),_transparent_55%)]" />
    <div
      class="absolute inset-0 opacity-[0.35] dark:opacity-[0.45]"
      style="
        background-image:
          linear-gradient(hsl(var(--border) / 0.55) 1px, transparent 1px),
          linear-gradient(90deg, hsl(var(--border) / 0.55) 1px, transparent 1px);
        background-size: 48px 48px;
        mask-image: radial-gradient(ellipse at center, black 20%, transparent 75%);
      "
    />
    <pre
      class="absolute left-1/2 top-1/2 max-w-xl -translate-x-1/2 -translate-y-1/2 select-none px-6 font-mono text-[11px] leading-relaxed text-primary/40 sm:text-xs dark:text-primary/35"
    ><code>async Task&lt;Score&gt; PlayAsync(string slug)
{
  var quiz = await LoadQuiz(slug);
  return Grade(quiz, answers);
}

const score = await fetch(`/api/v1/quizzes/${slug}/play`)
  .then((r) => r.json())</code></pre>
  </div>
</template>
