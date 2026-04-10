<script setup lang="ts">
import { ref, reactive } from 'vue'
import { watchDebounced } from '@vueuse/core'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { useToast } from '../composables/useToast'
import { api, RpcError } from '../api/client'
import type { IFieldError } from '../api/api'

const router = useRouter()
const { login, register } = useAuth()
const { show: showToast } = useToast()

const mode = ref<'login' | 'register'>('login')

const loginForm = reactive({ login: '', password: '' })
const registerForm = reactive({
  login: '',
  email: '',
  password: '',
  firstName: '',
  lastName: '',
})

const errors = ref<IFieldError[]>([])
const globalError = ref('')
const loading = ref(false)
const touched = reactive<Record<string, boolean>>({})

function markTouched(field: string) {
  touched[field] = true
}

function resetRegister() {
  errors.value = []
  globalError.value = ''
  Object.keys(touched).forEach((k) => delete touched[k])
}

function fieldErrorMessage(fe: IFieldError): string {
  const c = fe.constraint
  switch (fe.error) {
    case 'required': return 'Required'
    case 'min': return `Minimum ${c?.min ?? '?'} characters`
    case 'max': return `Maximum ${c?.max ?? '?'} characters`
    case 'format': return 'Invalid email address'
    default: return fe.error
  }
}

function fieldError(field: string): string | undefined {
  if (mode.value === 'register' && !touched[field]) return undefined
  const fe = errors.value.find((e) => e.field === field)
  return fe ? fieldErrorMessage(fe) : undefined
}

// Live validation while typing in the register form
watchDebounced(
  () => ({ ...registerForm }),
  async () => {
    if (mode.value !== 'register') return
    if (!Object.values(touched).some(Boolean)) return
    try {
      const fieldErrors = await api.auth.validateStudent({
        studentDraft: {
          login: registerForm.login,
          email: registerForm.email,
          password: registerForm.password,
          firstName: registerForm.firstName,
          lastName: registerForm.lastName,
        },
      })
      errors.value = fieldErrors ?? []
    } catch {
      // silent — validation errors don't block the UI
    }
  },
  { debounce: 400, deep: true },
)

async function handleLogin() {
  errors.value = []
  globalError.value = ''

  const fieldErrors = await api.auth.validateStudentLogin({
    studentLogin: { login: loginForm.login, password: loginForm.password },
  })
  if (fieldErrors?.length) {
    errors.value = fieldErrors
    return
  }

  loading.value = true
  try {
    await login(loginForm.login, loginForm.password)
    router.push('/courses')
  } catch (e) {
    globalError.value = e instanceof RpcError ? e.message : 'Login failed'
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  errors.value = []
  globalError.value = ''

  // Mark all fields as touched so field errors show on submit
  Object.keys(registerForm).forEach((k) => (touched[k] = true))

  const fieldErrors = await api.auth.validateStudent({
    studentDraft: {
      login: registerForm.login,
      email: registerForm.email,
      password: registerForm.password,
      firstName: registerForm.firstName,
      lastName: registerForm.lastName,
    },
  })
  if (fieldErrors?.length) {
    errors.value = fieldErrors
    return
  }

  loading.value = true
  try {
    await register({
      login: registerForm.login,
      email: registerForm.email,
      password: registerForm.password,
      firstName: registerForm.firstName,
      lastName: registerForm.lastName,
    })
    showToast('Account created successfully!', 'ok')
    router.push('/courses')
  } catch (e) {
    const msg = e instanceof RpcError ? e.message : 'Registration failed'
    showToast(msg, 'bad')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-bg flex items-center justify-center p-4">
    <div class="panel w-full max-w-md">
      <h1 class="font-heading font-bold text-2xl mb-1">Backend Courses</h1>
      <p class="text-ink-soft text-sm mb-6">Sign in to access courses and track your progress.</p>

      <div class="flex gap-2 mb-6">
        <button
          class="flex-1 py-2 text-sm font-semibold rounded-btn transition-colors"
          :class="mode === 'login' ? 'bg-brand text-white' : 'bg-white text-ink border border-line'"
          @click="mode = 'login'; errors = []; globalError = ''"
        >
          Login
        </button>
        <button
          class="flex-1 py-2 text-sm font-semibold rounded-btn transition-colors"
          :class="mode === 'register' ? 'bg-brand text-white' : 'bg-white text-ink border border-line'"
          @click="mode = 'register'; resetRegister()"
        >
          Register
        </button>
      </div>

      <div v-if="globalError" class="mb-4 p-3 rounded-btn bg-bad-bg text-bad text-sm">
        {{ globalError }}
      </div>

      <form v-if="mode === 'login'" @submit.prevent="handleLogin" class="flex flex-col gap-4">
        <div>
          <label class="block text-sm font-medium mb-1">Login</label>
          <input
            v-model="loginForm.login"
            type="text"
            class="w-full border rounded-btn px-3 py-2 text-sm outline-none focus:border-brand transition-colors"
            :class="fieldError('login') ? 'border-bad' : 'border-line'"
            placeholder="your_login"
          />
          <p v-if="fieldError('login')" class="text-bad text-xs mt-1">{{ fieldError('login') }}</p>
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Password</label>
          <input
            v-model="loginForm.password"
            type="password"
            class="w-full border rounded-btn px-3 py-2 text-sm outline-none focus:border-brand transition-colors"
            :class="fieldError('password') ? 'border-bad' : 'border-line'"
            placeholder="••••••••"
          />
          <p v-if="fieldError('password')" class="text-bad text-xs mt-1">{{ fieldError('password') }}</p>
        </div>
        <button type="submit" class="btn-main w-full" :disabled="loading">
          {{ loading ? 'Signing in…' : 'Sign In' }}
        </button>
      </form>

      <form v-else @submit.prevent="handleRegister" class="flex flex-col gap-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium mb-1">First Name</label>
            <input
              v-model="registerForm.firstName"
              type="text"
              class="w-full border rounded-btn px-3 py-2 text-sm outline-none focus:border-brand transition-colors"
              :class="fieldError('firstName') ? 'border-bad' : 'border-line'"
              @input="markTouched('firstName')"
            />
            <p v-if="fieldError('firstName')" class="text-bad text-xs mt-1">{{ fieldError('firstName') }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">Last Name</label>
            <input
              v-model="registerForm.lastName"
              type="text"
              class="w-full border rounded-btn px-3 py-2 text-sm outline-none focus:border-brand transition-colors"
              :class="fieldError('lastName') ? 'border-bad' : 'border-line'"
              @input="markTouched('lastName')"
            />
            <p v-if="fieldError('lastName')" class="text-bad text-xs mt-1">{{ fieldError('lastName') }}</p>
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Login</label>
          <input
            v-model="registerForm.login"
            type="text"
            class="w-full border rounded-btn px-3 py-2 text-sm outline-none focus:border-brand transition-colors"
            :class="fieldError('login') ? 'border-bad' : 'border-line'"
            @input="markTouched('login')"
          />
          <p v-if="fieldError('login')" class="text-bad text-xs mt-1">{{ fieldError('login') }}</p>
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Email</label>
          <input
            v-model="registerForm.email"
            type="email"
            class="w-full border rounded-btn px-3 py-2 text-sm outline-none focus:border-brand transition-colors"
            :class="fieldError('email') ? 'border-bad' : 'border-line'"
            @input="markTouched('email')"
          />
          <p v-if="fieldError('email')" class="text-bad text-xs mt-1">{{ fieldError('email') }}</p>
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Password</label>
          <input
            v-model="registerForm.password"
            type="password"
            class="w-full border rounded-btn px-3 py-2 text-sm outline-none focus:border-brand transition-colors"
            :class="fieldError('password') ? 'border-bad' : 'border-line'"
            placeholder="••••••••"
            @input="markTouched('password')"
          />
          <p v-if="fieldError('password')" class="text-bad text-xs mt-1">{{ fieldError('password') }}</p>
        </div>
        <button type="submit" class="btn-main w-full" :disabled="loading">
          {{ loading ? 'Creating account…' : 'Create Account' }}
        </button>
      </form>
    </div>
  </div>
</template>
