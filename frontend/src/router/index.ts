import { createRouter, createWebHashHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import CoursesView from '../views/CoursesView.vue'
import ExamView from '../views/ExamView.vue'
import HistoryView from '../views/HistoryView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/login', component: LoginView, meta: { public: true } },
    { path: '/', redirect: '/courses' },
    { path: '/courses', component: CoursesView },
    { path: '/exam/:courseId', component: ExamView },
    { path: '/history', component: HistoryView },
    { path: '/profile', component: ProfileView },
  ],
})

router.beforeEach((to) => {
  const token = localStorage.getItem('token')
  if (!to.meta.public && !token) {
    return '/login'
  }
  if (to.path === '/login' && token) {
    return '/courses'
  }
})

export default router
