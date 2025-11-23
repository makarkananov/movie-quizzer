<template>
  <div class="login-container">
    <div class="card" style="max-width: 400px; margin: 100px auto;">
      <h1 style="text-align: center; margin-bottom: 32px; color: #667eea;">Movie Quizzer</h1>
      <h2 style="text-align: center; margin-bottom: 24px;">Вход</h2>
      
      <form @submit.prevent="handleLogin">
        <div style="margin-bottom: 20px;">
          <label style="display: block; margin-bottom: 8px; font-weight: 600;">Email</label>
          <input
            v-model="email"
            type="email"
            required
            placeholder="your@email.com"
          />
        </div>
        
        <div style="margin-bottom: 24px;">
          <label style="display: block; margin-bottom: 8px; font-weight: 600;">Пароль</label>
          <input
            v-model="password"
            type="password"
            required
            placeholder="••••••••"
          />
        </div>
        
        <button type="submit" class="btn-primary" style="width: 100%; margin-bottom: 16px;" :disabled="loading">
          {{ loading ? 'Вход...' : 'Войти' }}
        </button>
        
        <div style="text-align: center; color: #666;">
          Нет аккаунта? 
          <router-link to="/register" style="color: #667eea; text-decoration: none; font-weight: 600;">
            Зарегистрироваться
          </router-link>
        </div>
        
        <div v-if="error" style="margin-top: 16px; padding: 12px; background: #fee; color: #c33; border-radius: 8px; text-align: center;">
          {{ error }}
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  loading.value = true
  error.value = ''
  
  try {
    await authStore.login(email.value, password.value)
    router.push('/home')
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Ошибка входа'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
}
</style>

