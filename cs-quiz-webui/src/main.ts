import { createApp } from 'vue'
import { setLucideProps } from '@lucide/vue'
import App from './App.vue'
import router from './router'
import { initTheme } from './lib/theme'
import { useAuth } from './lib/useAuth'
import './assets/index.css'

// Modern-minimal defaults: thin stroke, compact UI size
setLucideProps({
  size: 16,
  strokeWidth: 1.75,
})

initTheme()

const { hydrate } = useAuth()

void hydrate().then(() => {
  createApp(App).use(router).mount('#app')
})
