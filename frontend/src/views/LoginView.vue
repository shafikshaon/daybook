<template>
  <div class="login-view">
    <div class="container">
      <div class="row justify-content-center align-items-center">
        <div class="col-12 col-md-6 col-lg-4">
          <div class="card shadow-lg">
            <div class="card-body p-5">
              <div class="text-center mb-4">
                <h2 class="text-purple fw-bold">Daybook</h2>
                <p class="text-muted">Personal Finance Tracker</p>
              </div>

              <h4 class="mb-4">Sign In</h4>

              <div v-if="error" class="alert alert-danger" role="alert">
                {{ error }}
              </div>

              <div v-if="showDefaultCredentials" class="alert alert-info" role="alert">
                <strong>Default Admin:</strong><br>
                Username: <code>admin</code><br>
                Password: <code>admin</code>
              </div>

              <form @submit.prevent="handleLogin">
                <div class="mb-3">
                  <label class="form-label">Username</label>
                  <input
                    type="text"
                    class="form-control"
                    v-model="form.username"
                    required
                    placeholder="Enter username"
                    autocomplete="username"
                  />
                </div>

                <div class="mb-3">
                  <label class="form-label">Password</label>
                  <input
                    type="password"
                    class="form-control"
                    v-model="form.password"
                    required
                    placeholder="Enter password"
                    autocomplete="current-password"
                  />
                </div>

                <div class="mb-3 form-check">
                  <input
                    type="checkbox"
                    class="form-check-input"
                    id="rememberMe"
                    v-model="form.rememberMe"
                  />
                  <label class="form-check-label" for="rememberMe">
                    Remember me
                  </label>
                </div>

                <button
                  type="submit"
                  class="btn btn-primary w-100 mb-3"
                  :disabled="loading"
                >
                  <span v-if="loading">
                    <span class="spinner-border spinner-border-sm me-2"></span>
                    Signing in...
                  </span>
                  <span v-else>Sign In</span>
                </button>
              </form>

              <div class="text-center">
                <p class="text-muted mb-0">
                  Don't have an account?
                  <router-link to="/signup" class="text-primary fw-bold">
                    Sign Up
                  </router-link>
                </p>
              </div>
            </div>
          </div>

          <div class="text-center mt-3">
            <small class="text-muted">
              © 2025 Daybook. All rights reserved.
            </small>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  username: '',
  password: '',
  rememberMe: false
})

const loading = ref(false)
const error = ref('')
const showDefaultCredentials = ref(true)

const handleLogin = async () => {
  try {
    loading.value = true
    error.value = ''

    const result = await authStore.login(form.value.username, form.value.password)

    // Only redirect if login was successful
    if (result && result.success) {
      // Wait a bit for token to be set in localStorage
      await new Promise(resolve => setTimeout(resolve, 100))

      // Redirect to dashboard on successful login
      router.push('/dashboard')
    } else {
      error.value = 'Login failed. Please check your credentials.'
    }
  } catch (err) {
    console.error('Login error:', err)
    error.value = err.message || 'Login failed. Please try again.'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  // If already logged in, redirect to dashboard
  if (authStore.isLoggedIn) {
    router.push('/dashboard')
  }
})
</script>

<style scoped>
.login-view {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: calc(100vh - 60px);
  width: 100%;
  overflow-x: hidden;
  padding: 60px 20px;
  display: flex;
  align-items: center;
}

.container {
  max-width: 100%;
  overflow-x: hidden;
}

.row {
  min-height: auto;
}

.card {
  border: none;
  border-radius: 15px;
  width: 100%;
}

.text-purple {
  color: #6f42c1;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  padding: 12px;
  font-weight: 600;
}

.btn-primary:hover {
  background: linear-gradient(135deg, #764ba2 0%, #667eea 100%);
}

code {
  background-color: #f8f9fa;
  padding: 2px 6px;
  border-radius: 4px;
  color: #e83e8c;
}
</style>
