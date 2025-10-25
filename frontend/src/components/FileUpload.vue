<template>
  <div class="file-upload-component">
    <div class="upload-label" v-if="label">
      <label>{{ label }}</label>
      <span v-if="required" class="required-star">*</span>
    </div>

    <!-- Upload Area -->
    <div
      class="upload-area"
      :class="{ 'dragover': isDragging, 'has-error': error }"
      @click="triggerFileInput"
      @dragover.prevent="handleDragOver"
      @dragleave.prevent="handleDragLeave"
      @drop.prevent="handleDrop"
    >
      <input
        ref="fileInput"
        type="file"
        :multiple="multiple"
        :accept="acceptedTypes"
        @change="handleFileSelect"
        style="display: none"
      />

      <div class="upload-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" fill="currentColor" viewBox="0 0 16 16">
          <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5"/>
          <path d="M7.646 1.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1-.708.708L8.5 2.707V11.5a.5.5 0 0 1-1 0V2.707L5.354 4.854a.5.5 0 1 1-.708-.708z"/>
        </svg>
      </div>

      <div class="upload-text">
        <p class="upload-main-text">
          <span class="upload-link">Click to upload</span> or drag and drop
        </p>
        <p class="upload-sub-text">
          {{ acceptedTypesText }} (Max {{ maxSizeMB }}MB per file)
        </p>
        <p class="upload-sub-text" v-if="multiple">
          Maximum {{ maxFiles }} files
        </p>
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <!-- Upload Progress -->
    <div v-if="uploading" class="upload-progress">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: uploadProgress + '%' }"></div>
      </div>
      <p class="progress-text">Uploading... {{ uploadProgress }}%</p>
    </div>

    <!-- Uploaded Files List -->
    <div v-if="uploadedFiles.length > 0" class="uploaded-files">
      <h4>Attached Files ({{ uploadedFiles.length }})</h4>
      <div class="files-list">
        <div
          v-for="(file, index) in uploadedFiles"
          :key="index"
          class="file-item"
        >
          <div class="file-icon">
            <img v-if="isImage(file)" :src="file.fileUrl" alt="preview" class="file-preview-img" />
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="32" height="32" fill="currentColor" viewBox="0 0 16 16">
              <path d="M5 4a.5.5 0 0 0 0 1h6a.5.5 0 0 0 0-1zm-.5 2.5A.5.5 0 0 1 5 6h6a.5.5 0 0 1 0 1H5a.5.5 0 0 1-.5-.5M5 8a.5.5 0 0 0 0 1h6a.5.5 0 0 0 0-1zm0 2a.5.5 0 0 0 0 1h3a.5.5 0 0 0 0-1z"/>
              <path d="M2 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2zm10-1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1"/>
            </svg>
          </div>

          <div class="file-info">
            <p class="file-name" :title="file.originalName">{{ file.originalName }}</p>
            <p class="file-size">{{ formatFileSize(file.fileSize) }}</p>
          </div>

          <button
            type="button"
            class="file-delete-btn"
            @click="removeFile(index)"
            :disabled="uploading"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
              <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0z"/>
              <path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4zM2.5 3h11V2h-11z"/>
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { api } from '@/services/api-backend'

const props = defineProps({
  label: {
    type: String,
    default: 'Upload Files'
  },
  modelValue: {
    type: Array,
    default: () => []
  },
  multiple: {
    type: Boolean,
    default: true
  },
  maxFiles: {
    type: Number,
    default: 10
  },
  maxSize: {
    type: Number,
    default: 10 * 1024 * 1024 // 10MB
  },
  acceptedTypes: {
    type: String,
    default: '.jpg,.jpeg,.png,.gif,.pdf,.doc,.docx,.xls,.xlsx,.txt,.csv'
  },
  required: {
    type: Boolean,
    default: false
  },
  autoUpload: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'upload-success', 'upload-error', 'file-removed'])

const fileInput = ref(null)
const uploadedFiles = ref([...props.modelValue])
const isDragging = ref(false)
const uploading = ref(false)
const uploadProgress = ref(0)
const error = ref('')

const maxSizeMB = computed(() => Math.round(props.maxSize / (1024 * 1024)))

const acceptedTypesText = computed(() => {
  const types = props.acceptedTypes.split(',').map(t => t.trim().toUpperCase().replace('.', ''))
  if (types.length <= 3) {
    return types.join(', ')
  }
  return types.slice(0, 3).join(', ') + ' and more'
})

// Watch for external changes to modelValue
watch(() => props.modelValue, (newVal) => {
  uploadedFiles.value = [...newVal]
}, { deep: true })

// Watch uploadedFiles and emit changes
watch(uploadedFiles, (newVal) => {
  emit('update:modelValue', newVal)
}, { deep: true })

const triggerFileInput = () => {
  if (!uploading.value) {
    fileInput.value.click()
  }
}

const handleDragOver = () => {
  isDragging.value = true
}

const handleDragLeave = () => {
  isDragging.value = false
}

const handleDrop = (e) => {
  isDragging.value = false
  const files = Array.from(e.dataTransfer.files)
  processFiles(files)
}

const handleFileSelect = (e) => {
  const files = Array.from(e.target.files)
  processFiles(files)
  // Reset input value to allow selecting the same file again
  e.target.value = ''
}

const processFiles = async (files) => {
  error.value = ''

  // Check total number of files
  if (!props.multiple && files.length > 1) {
    error.value = 'Only one file allowed'
    return
  }

  if (uploadedFiles.value.length + files.length > props.maxFiles) {
    error.value = `Maximum ${props.maxFiles} files allowed`
    return
  }

  // Validate each file
  const validFiles = []
  for (const file of files) {
    // Check file size
    if (file.size > props.maxSize) {
      error.value = `File "${file.name}" exceeds maximum size of ${maxSizeMB.value}MB`
      continue
    }

    // Check file type
    const fileExt = '.' + file.name.split('.').pop().toLowerCase()
    if (!props.acceptedTypes.includes(fileExt)) {
      error.value = `File type "${fileExt}" not allowed for "${file.name}"`
      continue
    }

    validFiles.push(file)
  }

  if (validFiles.length === 0) {
    return
  }

  if (props.autoUpload) {
    await uploadFiles(validFiles)
  }
}

const uploadFiles = async (files) => {
  uploading.value = true
  uploadProgress.value = 0
  error.value = ''

  try {
    const formData = new FormData()
    files.forEach(file => {
      formData.append('files', file)
    })

    const response = await api.post('/uploads', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        uploadProgress.value = Math.round((progressEvent.loaded * 100) / progressEvent.total)
      }
    })

    // API interceptor extracts data, so response.data is already the nested data object
    if (response.data && response.data.files) {
      uploadedFiles.value = [...uploadedFiles.value, ...response.data.files]
      emit('upload-success', response.data.files)
    }
  } catch (err) {
    error.value = err.response?.data?.error || err.message || 'Upload failed. Please try again.'
    emit('upload-error', err)
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}

const removeFile = async (index) => {
  const file = uploadedFiles.value[index]

  try {
    const filename = file.fileName || file.filename

    if (filename) {
      await api.delete(`/uploads/${filename}`)
    }

    uploadedFiles.value.splice(index, 1)
    emit('file-removed', file)
  } catch (err) {
    console.error('Failed to delete file:', err)
    error.value = err.message || 'Failed to delete file'
  }
}

const isImage = (file) => {
  const imageTypes = ['.jpg', '.jpeg', '.png', '.gif']
  const ext = '.' + (file.originalName || file.fileName || '').split('.').pop().toLowerCase()
  return imageTypes.includes(ext)
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

// Expose methods to parent component
defineExpose({
  uploadFiles,
  clearFiles: () => { uploadedFiles.value = [] }
})
</script>

<style scoped>
.file-upload-component {
  width: 100%;
}

.upload-label {
  margin-bottom: 8px;
  font-size: 15px;
  font-weight: 600;
  color: #0a2540;
}

.required-star {
  color: #df1b41;
  margin-left: 4px;
}

.upload-area {
  border: 2px dashed #c7d2e0;
  border-radius: 12px;
  padding: 32px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background-color: #f6f9fc;
}

.upload-area:hover {
  border-color: #635bff;
  background-color: rgba(99, 91, 255, 0.02);
}

.upload-area.dragover {
  border-color: #635bff;
  background-color: rgba(99, 91, 255, 0.08);
}

.upload-area.has-error {
  border-color: #df1b41;
  background-color: rgba(223, 27, 65, 0.02);
}

.upload-icon {
  color: #8898aa;
  margin-bottom: 16px;
}

.upload-area:hover .upload-icon {
  color: #635bff;
}

.upload-text {
  color: #425466;
}

.upload-main-text {
  font-size: 16px;
  margin-bottom: 8px;
  font-weight: 500;
}

.upload-link {
  color: #635bff;
  font-weight: 600;
}

.upload-sub-text {
  font-size: 14px;
  color: #8898aa;
  margin: 4px 0;
}

.error-message {
  margin-top: 8px;
  padding: 12px;
  background-color: rgba(223, 27, 65, 0.08);
  border: 1px solid #df1b41;
  border-radius: 8px;
  color: #df1b41;
  font-size: 14px;
}

.upload-progress {
  margin-top: 16px;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background-color: #e3e8ee;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #635bff 0%, #5851d9 100%);
  transition: width 0.3s ease;
}

.progress-text {
  margin-top: 8px;
  font-size: 14px;
  color: #635bff;
  font-weight: 600;
}

.uploaded-files {
  margin-top: 24px;
}

.uploaded-files h4 {
  font-size: 16px;
  font-weight: 600;
  color: #0a2540;
  margin-bottom: 12px;
}

.files-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background-color: #ffffff;
  border: 1px solid #e3e8ee;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.file-item:hover {
  border-color: #635bff;
  box-shadow: 0 2px 8px rgba(99, 91, 255, 0.1);
}

.file-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: #8898aa;
}

.file-preview-img {
  width: 48px;
  height: 48px;
  object-fit: cover;
  border-radius: 6px;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-size: 14px;
  font-weight: 500;
  color: #0a2540;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  font-size: 13px;
  color: #8898aa;
}

.file-delete-btn {
  background: none;
  border: none;
  color: #df1b41;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.file-delete-btn:hover {
  background-color: rgba(223, 27, 65, 0.08);
}

.file-delete-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Dark Mode */
.dark-mode .upload-label {
  color: #f9fafb;
}

.dark-mode .upload-area {
  background-color: #111827;
  border-color: #374151;
}

.dark-mode .upload-area:hover {
  background-color: rgba(99, 91, 255, 0.05);
}

.dark-mode .upload-text {
  color: #e5e7eb;
}

.dark-mode .upload-sub-text {
  color: #9ca3af;
}

.dark-mode .file-item {
  background-color: #1f2937;
  border-color: #374151;
}

.dark-mode .file-name {
  color: #f9fafb;
}

.dark-mode .error-message {
  background-color: rgba(223, 27, 65, 0.15);
}

/* Mobile Styles */
@media (max-width: 767px) {
  .upload-area {
    padding: 24px 16px;
  }

  .upload-icon svg {
    width: 40px;
    height: 40px;
  }

  .upload-main-text {
    font-size: 15px;
  }

  .upload-sub-text {
    font-size: 13px;
  }

  .file-item {
    padding: 10px;
  }

  .file-icon {
    width: 40px;
    height: 40px;
  }

  .file-icon svg {
    width: 24px;
    height: 24px;
  }

  .file-preview-img {
    width: 40px;
    height: 40px;
  }
}
</style>
