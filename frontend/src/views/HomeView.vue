<template>
  <div class="home-container">
    <nav style="background: white; padding: 20px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); margin-bottom: 40px;">
      <div class="container" style="display: flex; justify-content: space-between; align-items: center;">
        <h1 style="color: #667eea; margin: 0;">Movie Quizzer</h1>
        <div style="display: flex; gap: 16px; align-items: center;">
          <span style="color: #666;">{{ authStore.user?.nickname }}</span>
          <router-link to="/profile" style="text-decoration: none; color: #667eea; font-weight: 600;">Профиль</router-link>
          <router-link to="/leaderboard" style="text-decoration: none; color: #667eea; font-weight: 600;">Рейтинг</router-link>
          <button @click="handleLogout" class="btn-secondary">Выйти</button>
        </div>
      </div>
    </nav>
    
    <div class="container">
      <div style="text-align: center; margin-bottom: 48px;">
        <h2 style="font-size: 36px; margin-bottom: 16px; color: white;">Выберите режим игры</h2>
        <p style="font-size: 18px; color: rgba(255,255,255,0.9);">Проверьте свои знания о фильмах!</p>
      </div>
      
      <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 24px;">
        <div class="card game-mode-card" @click="startGame('frame')">
          <div style="font-size: 48px; margin-bottom: 16px;">🎬</div>
          <h3 style="margin-bottom: 12px;">По кадру</h3>
          <p style="color: #666; margin-bottom: 24px;">Угадайте фильм по изображению кадра</p>
          <button class="btn-primary" style="width: 100%;">Играть</button>
        </div>
        
        <div class="card game-mode-card" @click="startGame('video')">
          <div style="font-size: 48px; margin-bottom: 16px;">🎥</div>
          <h3 style="margin-bottom: 12px;">По видеофрагменту</h3>
          <p style="color: #666; margin-bottom: 24px;">Угадайте фильм по видеоролику</p>
          <button class="btn-primary" style="width: 100%;">Играть</button>
        </div>
        
        <div class="card game-mode-card" @click="startGame('quote')">
          <div style="font-size: 48px; margin-bottom: 16px;">💬</div>
          <h3 style="margin-bottom: 12px;">По цитате</h3>
          <p style="color: #666; margin-bottom: 24px;">Угадайте фильм по известной цитате</p>
          <button class="btn-primary" style="width: 100%;">Играть</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

function startGame(mode: string) {
  router.push(`/game/${mode}`)
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
}

.game-mode-card {
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.game-mode-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}
</style>

