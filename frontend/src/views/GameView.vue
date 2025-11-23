<template>
  <div class="game-container">
    <div class="container" style="max-width: 900px;">
      <div class="card">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px;">
          <h2>Вопрос {{ currentQuestionNumber }} / {{ totalQuestions }}</h2>
          <div class="timer" :class="{ 'timer-warning': timeLeft < 20, 'timer-danger': timeLeft < 10 }">
            ⏱️ {{ formatTime(timeLeft) }}
          </div>
        </div>
        
        <div v-if="loading" style="text-align: center; padding: 40px;">
          <p>Загрузка вопроса...</p>
        </div>
        
        <div v-else-if="question" style="position: relative; z-index: 10;">
          <!-- Кадр -->
          <div v-if="question.type === 'frame' && question.image_url" style="margin-bottom: 24px;">
            <img 
              :src="getMediaUrl(question.image_url)" 
              alt="Кадр из фильма"
              @error="handleImageError"
              @load="handleImageLoad"
              style="width: 100%; max-height: 400px; object-fit: contain; border-radius: 8px;"
            />
            <div v-if="imageError" style="padding: 20px; background: #fff3cd; border-radius: 8px; color: #856404;">
              <p>⚠️ Изображение не загружено. Убедитесь, что файл загружен в MinIO.</p>
              <p style="font-size: 12px; margin-top: 8px;">Путь: {{ question.image_url }}</p>
            </div>
          </div>
          
          <!-- Видео -->
          <div v-if="question.type === 'video' && question.video_url" style="margin-bottom: 24px;">
            <video 
              :src="getMediaUrl(question.video_url)" 
              @error="handleVideoError"
              controls
              autoplay
              style="width: 100%; max-height: 400px; border-radius: 8px;"
            ></video>
            <div v-if="videoError" style="padding: 20px; background: #fff3cd; border-radius: 8px; color: #856404; margin-top: 8px;">
              <p>⚠️ Видео не загружено. Убедитесь, что файл загружен в MinIO.</p>
              <p style="font-size: 12px; margin-top: 8px;">Путь: {{ question.video_url }}</p>
            </div>
          </div>
          
          <!-- Цитата -->
          <div v-if="question.type === 'quote' && question.text" style="margin-bottom: 24px;">
            <div style="background: #f5f5f5; padding: 24px; border-radius: 8px; border-left: 4px solid #667eea;">
              <p style="font-size: 20px; font-style: italic; color: #333;">"{{ question.text }}"</p>
            </div>
          </div>
          
          
          <!-- Варианты ответов -->
          <div v-if="question.options && Array.isArray(question.options) && question.options.length > 0" 
               style="display: grid; gap: 12px; margin-bottom: 24px; position: relative; z-index: 100;">
            <button
              v-for="(option, index) in question.options"
              :key="`option-${index}-${option}`"
              @click="handleOptionClick(option)"
              :disabled="answered || loading"
              class="btn-option"
              :class="{ 'btn-option-selected': selectedAnswer === option }"
              type="button"
            >
              {{ option }}
            </button>
          </div>
          
          <!-- Поле ввода для цитат (если нет вариантов) -->
          <div v-else style="margin-bottom: 24px;">
            <input
              v-model="textAnswer"
              type="text"
              placeholder="Введите название фильма"
              @keyup.enter="submitAnswer(textAnswer)"
              :disabled="answered"
              style="margin-bottom: 16px;"
            />
            <button 
              @click="handleSubmitClick"
              :disabled="answered || !textAnswer"
              class="btn-primary"
              style="width: 100%;"
            >
              Ответить
            </button>
          </div>
          
          <div v-if="answerResult" style="margin-top: 24px; padding: 16px; border-radius: 8px;" 
               :style="{ background: answerResult.correct ? '#d4edda' : '#f8d7da', color: answerResult.correct ? '#155724' : '#721c24' }">
            <p style="font-weight: 600; margin-bottom: 8px;">
              {{ answerResult.correct ? '✅ Правильно!' : '❌ Неправильно' }}
            </p>
            <p>Начислено очков: {{ answerResult.score }}</p>
          </div>
        </div>
        
        <div v-if="sessionStatus === 'finished'" style="text-align: center; padding: 40px;">
          <h2 style="margin-bottom: 24px;">Раунд завершен!</h2>
          <div style="margin-bottom: 24px;">
            <p style="font-size: 24px; margin-bottom: 8px;">Ваш результат:</p>
            <p style="font-size: 36px; font-weight: bold; color: #667eea;">{{ sessionSummary?.score || 0 }} очков</p>
            <p style="margin-top: 16px;">Правильных ответов: {{ sessionSummary?.correct_answers || 0 }} / {{ sessionSummary?.total_questions || 0 }}</p>
          </div>
          <div style="display: flex; gap: 16px; justify-content: center;">
            <button @click="goHome" class="btn-primary">На главную</button>
            <button @click="startNewGame" class="btn-secondary">Новая игра</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'

const route = useRoute()
const router = useRouter()

const mode = route.params.mode as string
const sessionId = ref<number | null>(null)
const question = ref<any>(null)
const currentQuestionNumber = ref(1)
const totalQuestions = ref(10)
const timeLeft = ref(60)
const timer = ref<number | null>(null)
const loading = ref(false)
const answered = ref(false)
const selectedAnswer = ref('')
const textAnswer = ref('')
const answerResult = ref<any>(null)
const sessionStatus = ref('in_progress')
const sessionSummary = ref<any>(null)
const startTime = ref<number>(Date.now())
const imageError = ref(false)
const videoError = ref(false)

function formatTime(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

function getMediaUrl(url: string | null): string {
  if (!url) return ''
  if (url.startsWith('http')) return url
  // Убираем начальный слеш если есть
  const cleanUrl = url.startsWith('/') ? url.slice(1) : url
  // URL формат: /api/media/{file}
  // file = путь из БД (например, "frames/matrix.jpg")
  return `/api/media/${cleanUrl}`
}

function handleImageError(event: Event) {
  console.error('Ошибка загрузки изображения:', event)
  imageError.value = true
}

function handleImageLoad() {
  imageError.value = false
}

function handleVideoError(event: Event) {
  console.error('Ошибка загрузки видео:', event)
  videoError.value = true
}

function handleOptionClick(option: string) {
  console.log('🔵🔵🔵 handleOptionClick ВЫЗВАН! 🔵🔵🔵', option)
  console.log('Состояние ДО проверок:', {
    answered: answered.value,
    loading: loading.value,
    sessionId: sessionId.value,
    hasQuestion: !!question.value,
    questionId: question.value?.id
  })
  
  // Временно отключаем проверки для теста - чтобы увидеть, вызывается ли функция вообще
  if (answered.value) {
    console.warn('⚠️ Уже отвечено, но продолжаем для теста')
    // return  // Закомментировано для теста
  }
  
  if (loading.value) {
    console.warn('⚠️ Идет загрузка, но продолжаем для теста')
    // return  // Закомментировано для теста
  }
  
  if (!sessionId.value || !question.value) {
    console.error('❌ Нет sessionId или question - это критично!', {
      sessionId: sessionId.value,
      hasQuestion: !!question.value
    })
    // Не показываем alert, просто логируем - это должно быть исправлено в startSession
    return
  }
  
  console.log('✅ Все проверки пройдены, отправляем ответ:', option)
  selectedAnswer.value = option
  
  // Вызываем напрямую
  submitAnswer(option).catch(err => {
    console.error('❌ Ошибка в submitAnswer:', err)
  })
}

function handleSubmitClick() {
  console.log('Кнопка нажата, textAnswer:', textAnswer.value)
  if (textAnswer.value && !answered.value) {
    submitAnswer(textAnswer.value)
  } else {
    console.warn('textAnswer пустой или уже отвечено, кнопка должна быть disabled')
  }
}

async function startSession() {
  loading.value = true
  try {
    const response = await api.post('/game/sessions', { mode })
    
    // Отладка: логируем весь ответ
    console.log('=== Ответ от сервера ===')
    console.log('Полный response:', response)
    console.log('response.data:', response.data)
    console.log('response.data.session:', response.data.session)
    console.log('response.data.question:', response.data.question)
    
    // Проверяем структуру ответа
    if (!response.data || !response.data.session) {
      console.error('❌ Неверная структура ответа:', response.data)
      alert('Ошибка: неверный формат ответа от сервера')
      router.push('/home')
      return
    }
    
    // Пробуем разные варианты доступа к ID
    const sessionIdValue = response.data.session.id || response.data.session.ID
    if (!sessionIdValue) {
      console.error('❌ session.id не найден. Доступные поля:', Object.keys(response.data.session))
      alert(`Ошибка: не найден ID сессии. Поля: ${Object.keys(response.data.session).join(', ')}`)
      router.push('/home')
      return
    }
    
    sessionId.value = sessionIdValue
    question.value = response.data.question
    totalQuestions.value = response.data.session.total_questions || response.data.session.TotalQuestions || 10
    
    console.log('✅ Сессия создана:', {
      sessionId: sessionId.value,
      totalQuestions: totalQuestions.value,
      questionId: question.value?.id,
      questionType: question.value?.type
    })
    
    // Отладка
    console.log('Question loaded:', question.value)
    console.log('Options:', question.value?.options)
    console.log('Options length:', question.value?.options?.length)
    
    startTimer()
  } catch (error: any) {
    console.error('❌ Failed to start session:', error)
    console.error('Детали ошибки:', {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status
    })
    alert(`Ошибка создания сессии: ${error.response?.data?.message || error.message}`)
    router.push('/home')
  } finally {
    loading.value = false
  }
}

function startTimer() {
  timeLeft.value = 60
  startTime.value = Date.now()
  timer.value = window.setInterval(() => {
    timeLeft.value--
    if (timeLeft.value <= 0) {
      handleTimeout()
    }
  }, 1000)
}

function stopTimer() {
  if (timer.value) {
    clearInterval(timer.value)
    timer.value = null
  }
}

async function handleTimeout() {
  stopTimer()
  if (!answered.value && question.value) {
    await submitAnswer('', true)
  }
}

async function submitAnswer(answer: string, timeout = false) {
  console.log('=== submitAnswer вызван ===', { answer, timeout })
  
  // Дополнительная проверка на всякий случай
  if (answered.value) {
    console.warn('Уже отвечено, игнорируем')
    return
  }
  
  if (!sessionId.value) {
    console.error('Нет sessionId')
    return
  }
  
  if (!question.value) {
    console.error('Нет question')
    return
  }
  
  console.log('Устанавливаем answered = true')
  answered.value = true
  stopTimer()
  
  const elapsedMs = Date.now() - startTime.value
  console.log('Отправляем запрос:', {
    sessionId: sessionId.value,
    questionId: question.value.id,
    answer,
    elapsedMs
  })
  
  try {
    const response = await api.post(`/game/sessions/${sessionId.value}/answers`, {
      question_id: question.value.id,
      answer: timeout ? '' : answer,
      elapsed_ms: elapsedMs
    })
    
    console.log('Ответ получен:', response.data)
    answerResult.value = response.data
    sessionStatus.value = response.data.session_status
    
    if (response.data.session_status === 'finished') {
      await loadSessionSummary()
    } else if (response.data.next_question) {
      setTimeout(async () => {
        await loadNextQuestion()
      }, 2000)
    } else {
      // Если нет next_question, загружаем следующий вручную
      setTimeout(async () => {
        await loadNextQuestion()
      }, 2000)
    }
  } catch (error: any) {
    console.error('Ошибка отправки ответа:', error)
    console.error('Детали ошибки:', {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status
    })
    alert(`Ошибка отправки ответа: ${error.response?.data?.message || error.message}`)
    answered.value = false
    startTimer()
  }
}

async function loadNextQuestion() {
  if (!sessionId.value) return
  
  loading.value = true
  answered.value = false
  selectedAnswer.value = ''
  textAnswer.value = ''
  answerResult.value = null
  imageError.value = false
  videoError.value = false
  
  try {
    const response = await api.get(`/game/sessions/${sessionId.value}/next`)
    question.value = response.data
    currentQuestionNumber.value++
    startTimer()
  } catch (error: any) {
    if (error.response?.status === 404) {
      sessionStatus.value = 'finished'
      await loadSessionSummary()
    }
  } finally {
    loading.value = false
  }
}

async function loadSessionSummary() {
  if (!sessionId.value) return
  
  try {
    const response = await api.get(`/game/sessions/${sessionId.value}`)
    sessionSummary.value = response.data
  } catch (error) {
    console.error('Failed to load summary:', error)
  }
}

function goHome() {
  router.push('/home')
}

async function startNewGame() {
  // Сбрасываем все состояние перед началом новой игры
  stopTimer()
  sessionId.value = null
  question.value = null
  currentQuestionNumber.value = 1
  totalQuestions.value = 10
  timeLeft.value = 60
  answered.value = false
  selectedAnswer.value = ''
  textAnswer.value = ''
  answerResult.value = null
  sessionStatus.value = 'in_progress'
  sessionSummary.value = null
  imageError.value = false
  videoError.value = false
  loading.value = false
  
  // Начинаем новую сессию
  await startSession()
}

onMounted(() => {
  startSession()
})

onUnmounted(() => {
  stopTimer()
})
</script>

<style scoped>
.game-container {
  min-height: 100vh;
  padding: 40px 0;
}

.timer {
  font-size: 24px;
  font-weight: bold;
  padding: 12px 24px;
  background: #e8f5e9;
  border-radius: 8px;
  color: #2e7d32;
}

.timer-warning {
  background: #fff3cd;
  color: #856404;
}

.timer-danger {
  background: #f8d7da;
  color: #721c24;
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

.btn-option {
  padding: 16px;
  background: #f5f5f5;
  border: 2px solid #e0e0e0;
  text-align: left;
  font-size: 16px;
  transition: all 0.2s;
  cursor: pointer;
  position: relative;
  z-index: 1;
}

.btn-option:hover:not(:disabled) {
  background: #e8e8e8;
  border-color: #667eea;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.btn-option:active:not(:disabled) {
  transform: translateY(0);
}

.btn-option:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  /* НЕ используем pointer-events: none, чтобы клики все равно логировались */
}

.btn-option-selected {
  background: #667eea;
  color: white;
  border-color: #667eea;
}
</style>

