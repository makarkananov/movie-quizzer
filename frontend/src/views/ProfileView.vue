<template>
  <div class="profile-container">
    <nav style="background: white; padding: 20px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); margin-bottom: 40px;">
      <div class="container" style="display: flex; justify-content: space-between; align-items: center;">
        <h1 style="color: #667eea; margin: 0;">Movie Quizzer</h1>
        <div style="display: flex; gap: 16px; align-items: center;">
          <router-link to="/home" style="text-decoration: none; color: #667eea; font-weight: 600;">Главная</router-link>
          <router-link to="/leaderboard" style="text-decoration: none; color: #667eea; font-weight: 600;">Рейтинг</router-link>
          <button @click="handleLogout" class="btn-secondary">Выйти</button>
        </div>
      </div>
    </nav>
    
    <div class="container" style="max-width: 900px;">
      <div class="card" style="margin-bottom: 24px;">
        <h2 style="margin-bottom: 24px;">Профиль</h2>
        <div v-if="authStore.user" style="margin-bottom: 32px;">
          <p style="font-size: 18px; margin-bottom: 8px;"><strong>Никнейм:</strong> {{ authStore.user.nickname }}</p>
          <p style="font-size: 18px;"><strong>Email:</strong> {{ authStore.user.email }}</p>
        </div>
      </div>
      
      <div class="card" style="margin-bottom: 24px;">
        <h2 style="margin-bottom: 24px;">Статистика</h2>
        <div v-if="loading">Загрузка...</div>
        <div v-else-if="profile" style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 24px;">
          <div style="text-align: center; padding: 24px; background: #f5f5f5; border-radius: 8px;">
            <div style="font-size: 36px; font-weight: bold; color: #667eea; margin-bottom: 8px;">
              {{ profile.total_sessions }}
            </div>
            <div style="color: #666;">Сыграно раундов</div>
          </div>
          <div style="text-align: center; padding: 24px; background: #f5f5f5; border-radius: 8px;">
            <div style="font-size: 36px; font-weight: bold; color: #667eea; margin-bottom: 8px;">
              {{ profile.total_answers }}
            </div>
            <div style="color: #666;">Всего ответов</div>
          </div>
          <div style="text-align: center; padding: 24px; background: #f5f5f5; border-radius: 8px;">
            <div style="font-size: 36px; font-weight: bold; color: #667eea; margin-bottom: 8px;">
              {{ profile.correct_answers }}
            </div>
            <div style="color: #666;">Правильных ответов</div>
          </div>
          <div style="text-align: center; padding: 24px; background: #f5f5f5; border-radius: 8px;">
            <div style="font-size: 36px; font-weight: bold; color: #667eea; margin-bottom: 8px;">
              {{ profile.total_score }}
            </div>
            <div style="color: #666;">Всего очков</div>
          </div>
        </div>
        <div v-if="profile" style="margin-top: 24px; padding: 24px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 8px; color: white; text-align: center;">
          <div style="font-size: 48px; font-weight: bold; margin-bottom: 8px;">
            {{ profile.accuracy_percent.toFixed(1) }}%
          </div>
          <div style="font-size: 18px;">Точность ответов</div>
        </div>
      </div>
      
      <div class="card">
        <h2 style="margin-bottom: 24px;">Достижения</h2>
        <div v-if="achievementsLoading">Загрузка...</div>
        <div v-else-if="achievements.length === 0" style="text-align: center; padding: 40px; color: #666;">
          Достижения пока недоступны
        </div>
        <div v-else style="display: grid; gap: 16px;">
          <div 
            v-for="achievement in achievements" 
            :key="achievement.id"
            style="padding: 20px; border-radius: 8px; border: 2px solid #e0e0e0;"
            :style="{ 
              background: achievement.earned ? '#d4edda' : '#f5f5f5',
              borderColor: achievement.earned ? '#28a745' : '#e0e0e0'
            }"
          >
            <div style="display: flex; align-items: center; gap: 16px;">
              <div style="font-size: 32px;">{{ achievement.earned ? '🏆' : '🔒' }}</div>
              <div style="flex: 1;">
                <h3 style="margin-bottom: 4px;">{{ achievement.title }}</h3>
                <p style="color: #666; margin: 0;">{{ achievement.description }}</p>
              </div>
            </div>
          </div>
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

const profile = ref<any>(null)
const achievements = ref<any[]>([])
const loading = ref(false)
const achievementsLoading = ref(false)

async function loadProfile() {
  loading.value = true
  try {
    const response = await api.get('/profile')
    profile.value = response.data
  } catch (error) {
    console.error('Failed to load profile:', error)
  } finally {
    loading.value = false
  }
}

async function loadAchievements() {
  achievementsLoading.value = true
  try {
    const response = await api.get('/profile/achievements')
    achievements.value = response.data
  } catch (error) {
    console.error('Failed to load achievements:', error)
  } finally {
    achievementsLoading.value = false
  }
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}

onMounted(() => {
  loadProfile()
  loadAchievements()
})
</script>

<style scoped>
.profile-container {
  min-height: 100vh;
  padding-bottom: 40px;
}
</style>

