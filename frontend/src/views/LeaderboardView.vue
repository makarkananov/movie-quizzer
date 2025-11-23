<template>
  <div class="leaderboard-container">
    <nav style="background: white; padding: 20px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); margin-bottom: 40px;">
      <div class="container" style="display: flex; justify-content: space-between; align-items: center;">
        <h1 style="color: #667eea; margin: 0;">Movie Quizzer</h1>
        <div style="display: flex; gap: 16px; align-items: center;">
          <router-link to="/home" style="text-decoration: none; color: #667eea; font-weight: 600;">Главная</router-link>
          <router-link to="/profile" style="text-decoration: none; color: #667eea; font-weight: 600;">Профиль</router-link>
          <button @click="handleLogout" class="btn-secondary">Выйти</button>
        </div>
      </div>
    </nav>
    
    <div class="container" style="max-width: 900px;">
      <div class="card" style="margin-bottom: 24px;" v-if="myEntry">
        <h2 style="margin-bottom: 24px;">Ваша позиция</h2>
        <div style="display: flex; justify-content: space-around; align-items: center; padding: 24px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 8px; color: white;">
          <div style="text-align: center;">
            <div style="font-size: 48px; font-weight: bold; margin-bottom: 8px;">#{{ myEntry.position }}</div>
            <div>Позиция</div>
          </div>
          <div style="text-align: center;">
            <div style="font-size: 48px; font-weight: bold; margin-bottom: 8px;">{{ myEntry.score }}</div>
            <div>Очков</div>
          </div>
          <div style="text-align: center;">
            <div style="font-size: 48px; font-weight: bold; margin-bottom: 8px;">{{ myEntry.accuracy.toFixed(1) }}%</div>
            <div>Точность</div>
          </div>
        </div>
      </div>
      
      <div class="card">
        <h2 style="margin-bottom: 24px;">Глобальный рейтинг</h2>
        <div v-if="loading">Загрузка...</div>
        <div v-else-if="leaderboard.length === 0" style="text-align: center; padding: 40px; color: #666;">
          Рейтинг пуст
        </div>
        <div v-else>
          <table style="width: 100%; border-collapse: collapse;">
            <thead>
              <tr style="background: #f5f5f5; border-bottom: 2px solid #e0e0e0;">
                <th style="padding: 16px; text-align: left;">Позиция</th>
                <th style="padding: 16px; text-align: left;">Никнейм</th>
                <th style="padding: 16px; text-align: right;">Очки</th>
                <th style="padding: 16px; text-align: right;">Точность</th>
              </tr>
            </thead>
            <tbody>
              <tr 
                v-for="entry in leaderboard" 
                :key="entry.user_id"
                style="border-bottom: 1px solid #e0e0e0;"
                :style="{ background: entry.user_id === myEntry?.user_id ? '#e3f2fd' : 'white' }"
              >
                <td style="padding: 16px; font-weight: bold; color: #667eea;">
                  {{ entry.position === 1 ? '🥇' : entry.position === 2 ? '🥈' : entry.position === 3 ? '🥉' : '#' + entry.position }}
                </td>
                <td style="padding: 16px; font-weight: 600;">{{ entry.nickname }}</td>
                <td style="padding: 16px; text-align: right; font-weight: 600;">{{ entry.score }}</td>
                <td style="padding: 16px; text-align: right;">{{ entry.accuracy.toFixed(1) }}%</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import api from '../api'

const router = useRouter()
const authStore = useAuthStore()

const leaderboard = ref<any[]>([])
const myEntry = ref<any>(null)
const loading = ref(false)

async function loadLeaderboard() {
  loading.value = true
  try {
    const [globalRes, myRes] = await Promise.all([
      api.get('/leaderboard/global'),
      api.get('/leaderboard/me')
    ])
    leaderboard.value = globalRes.data
    myEntry.value = myRes.data
  } catch (error) {
    console.error('Failed to load leaderboard:', error)
  } finally {
    loading.value = false
  }
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}

onMounted(() => {
  loadLeaderboard()
})
</script>

<style scoped>
.leaderboard-container {
  min-height: 100vh;
  padding-bottom: 40px;
}

table tr:hover {
  background: #f9f9f9 !important;
}
</style>

