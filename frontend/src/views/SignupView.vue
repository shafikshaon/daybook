<template>
  <div class="signup-view">
    <div class="container">
      <div class="row justify-content-center align-items-center">
        <div class="col-12 col-md-6 col-lg-5">
          <div class="card shadow-lg">
            <div class="card-body p-5">
              <div class="text-center mb-4">
                <h2 class="text-purple fw-bold">Daybook</h2>
                <p class="text-muted">Personal Finance Tracker</p>
              </div>

              <h4 class="mb-4">Create Account</h4>

              <div v-if="error" class="alert alert-danger" role="alert">
                {{ error }}
              </div>

              <form @submit.prevent="handleSignup">
                <div class="mb-3">
                  <label class="form-label">Full Name</label>
                  <input
                    type="text"
                    class="form-control"
                    v-model="form.fullName"
                    required
                    placeholder="Enter your full name"
                  />
                </div>

                <div class="mb-3">
                  <label class="form-label">Username *</label>
                  <input
                    type="text"
                    class="form-control"
                    v-model="form.username"
                    required
                    placeholder="Choose a username"
                    autocomplete="username"
                    @blur="validateUsername"
                  />
                  <small class="text-muted">Only letters, numbers, and underscores</small>
                </div>

                <div class="mb-3">
                  <label class="form-label">Email *</label>
                  <input
                    type="email"
                    class="form-control"
                    v-model="form.email"
                    required
                    placeholder="Enter your email"
                    autocomplete="email"
                  />
                </div>

                <div class="mb-3">
                  <label class="form-label">Password *</label>
                  <input
                    type="password"
                    class="form-control"
                    v-model="form.password"
                    required
                    placeholder="Choose a password"
                    autocomplete="new-password"
                    @input="checkPasswordStrength"
                  />
                  <div v-if="form.password" class="mt-2">
                    <small :class="passwordStrengthClass">
                      Password strength: {{ passwordStrength }}
                    </small>
                  </div>
                </div>

                <div class="mb-3">
                  <label class="form-label">Confirm Password *</label>
                  <input
                    type="password"
                    class="form-control"
                    v-model="form.confirmPassword"
                    required
                    placeholder="Confirm your password"
                    autocomplete="new-password"
                  />
                  <small v-if="form.confirmPassword && form.password !== form.confirmPassword" class="text-danger">
                    Passwords do not match
                  </small>
                </div>

                <div class="mb-3 form-check">
                  <input
                    type="checkbox"
                    class="form-check-input"
                    id="terms"
                    v-model="form.acceptTerms"
                    required
                  />
                  <label class="form-check-label" for="terms">
                    I agree to the Terms and Conditions
                  </label>
                </div>

                <button
                  type="submit"
                  class="btn btn-primary w-100 mb-3"
                  :disabled="loading || form.password !== form.confirmPassword"
                >
                  <span v-if="loading">
                    <span class="spinner-border spinner-border-sm me-2"></span>
                    Creating account...
                  </span>
                  <span v-else>Create Account</span>
                </button>
              </form>

              <div class="text-center">
                <p class="text-muted mb-0">
                  Already have an account?
                  <router-link to="/login" class="text-primary fw-bold">
                    Sign In
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
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  fullName: '',
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  acceptTerms: false
})

const loading = ref(false)
const error = ref('')
const passwordStrength = ref('')

const passwordStrengthClass = computed(() => {
  switch (passwordStrength.value) {
    case 'Weak':
      return 'text-danger'
    case 'Medium':
      return 'text-warning'
    case 'Strong':
      return 'text-success'
    default:
      return 'text-muted'
  }
})

const checkPasswordStrength = () => {
  const password = form.value.password

  if (password.length < 6) {
    passwordStrength.value = 'Weak'
  } else if (password.length < 10) {
    passwordStrength.value = 'Medium'
  } else {
    passwordStrength.value = 'Strong'
  }
}

const validateUsername = () => {
  const username = form.value.username
  const regex = /^[a-zA-Z0-9_]+$/

  if (username && !regex.test(username)) {
    error.value = 'Username can only contain letters, numbers, and underscores'
  } else {
    error.value = ''
  }
}

const handleSignup = async () => {
  try {
    loading.value = true
    error.value = ''

    // Validate form
    if (form.value.password !== form.value.confirmPassword) {
      throw new Error('Passwords do not match')
    }

    if (form.value.password.length < 6) {
      throw new Error('Password must be at least 6 characters long')
    }

    if (!form.value.acceptTerms) {
      throw new Error('You must accept the terms and conditions')
    }

    await authStore.signup({
      fullName: form.value.fullName,
      username: form.value.username,
      email: form.value.email,
      password: form.value.password
    })

    // Redirect to dashboard on successful signup
    router.push('/dashboard')
  } catch (err) {
    error.value = err.message || 'Signup failed. Please try again.'
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
.signup-view {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: calc(100vh - 60px);
  width: 100%;
  overflow-x: hidden;
  padding: 40px 20px;
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

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
