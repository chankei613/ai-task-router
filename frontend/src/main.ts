import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router'
import 'virtual:uno.css'
import '@unocss/reset/tailwind.css'
import './assets/globals.css'
import App from './App.vue'
import ModelsView from './pages/ModelsView.vue'
import RouteView from './pages/RouteView.vue'
import RouteDetailView from './pages/RouteDetailView.vue'
import HelpView from './pages/HelpView.vue'
import SettingsView from './pages/SettingsView.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/models' },
    { path: '/models', component: ModelsView },
    { path: '/route', component: RouteView },
    { path: '/routes/:id', component: RouteDetailView },
    { path: '/help', component: HelpView },
    { path: '/settings', component: SettingsView },
  ],
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
