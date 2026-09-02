import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import LoginView from '@/views/LoginView.vue'
import SignupView from '@/views/SignupView.vue'
import AccountView from '@/views/AccountView.vue'
import ProfileView from '@/views/ProfileView.vue'
import AboutMeView from '@/views/AboutMeView.vue'
import ManagementView from '@/views/ManagementView.vue'
import AddQuizView from '@/views/AddQuizView.vue'
import QuestionsListView from '@/views/QuestionsListView.vue'
import QuizPlayView from '@/views/QuizPlayView.vue'
import { getToken, getStoredUser } from '@/lib/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    { path: '/login', name: 'login', component: LoginView, meta: { guest: true } },
    { path: '/signup', name: 'signup', component: SignupView, meta: { guest: true } },
    { path: '/account', name: 'account', component: AccountView, meta: { requiresAuth: true } },
    { path: '/profile', name: 'profile-self', component: ProfileView, meta: { requiresAuth: true } },
    { path: '/profile/:username', name: 'profile', component: ProfileView },
    { path: '/about-me', name: 'about-me', component: AboutMeView },
    {
      path: '/management',
      name: 'management',
      component: ManagementView,
      meta: { requiresAuth: true, requiresAdmin: true },
    },
    {
      path: '/management/add-quiz',
      name: 'add-quiz',
      component: AddQuizView,
      meta: { requiresAuth: true, requiresAdmin: true },
    },
    {
      path: '/management/quiz/:quizName/add-question',
      name: 'add-question',
      redirect: (to) => ({
        name: 'questions-list',
        params: { quizName: to.params.quizName },
      }),
    },
    {
      path: '/management/quiz/:quizName/questions-list',
      name: 'questions-list',
      component: QuestionsListView,
      meta: { requiresAuth: true, requiresAdmin: true },
    },
    {
      path: '/quiz/:slug',
      name: 'quiz-play',
      component: QuizPlayView,
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach((to) => {
  const token = getToken()
  const user = getStoredUser()

  if (to.meta.requiresAuth && !token) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.requiresAdmin && !user?.is_admin) {
    return { name: 'home' }
  }
  if (to.meta.guest && token) {
    return { name: 'home' }
  }
  return true
})

export default router
