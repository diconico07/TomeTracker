import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/planning',
      name: 'planning',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/PlanningView.vue'),
    },
    {
      path: '/missing',
      name: 'missing',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/MissingView.vue'),
    },
    {
      path: '/series/:id',
      name: 'series',
      component: () => import('../views/SeriesView.vue'),
    }
  ],
})

export default router
