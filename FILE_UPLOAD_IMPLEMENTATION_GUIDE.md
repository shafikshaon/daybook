# File Upload Feature - Implementation Guide

## 📋 Overview

This guide explains how to use the new file upload feature in the Daybook application. File uploads are now supported for:
- **Transactions** ✅
- **Savings Goals** ✅
- **Savings Contributions** ✅
- **Fixed Deposits** ✅

## 🎯 Features

### Backend Features
- ✅ Multiple file uploads (up to 10 files per upload)
- ✅ File size limit: 10MB per file
- ✅ Supported file types: JPG, JPEG, PNG, GIF, PDF, DOC, DOCX, XLS, XLSX, TXT, CSV
- ✅ Secure file storage (user-specific directories)
- ✅ File deletion support
- ✅ File serving with authentication
- ✅ Auto-migration for new attachment fields

### Frontend Features
- ✅ Drag-and-drop support
- ✅ Multiple file selection
- ✅ Image preview
- ✅ File size formatting
- ✅ Upload progress indicator
- ✅ Error handling
- ✅ Dark mode support
- ✅ Mobile-responsive design

## 🔧 Backend Implementation

### 1. API Endpoints

All endpoints require authentication (Bearer token).

#### Upload Multiple Files
```http
POST /api/v1/uploads
Content-Type: multipart/form-data
Authorization: Bearer <token>

Body: files[] (multiple files)
```

**Response:**
```json
{
  "success": true,
  "message": "Files uploaded successfully",
  "data": {
    "files": [
      {
        "fileName": "invoice_1234567890_abc12345.pdf",
        "originalName": "invoice.pdf",
        "filePath": "./uploads/{userId}/invoice_1234567890_abc12345.pdf",
        "fileUrl": "/api/v1/uploads/{userId}/invoice_1234567890_abc12345.pdf",
        "fileSize": 245678,
        "mimeType": "application/pdf"
      }
    ],
    "uploadedCount": 1,
    "totalFiles": 1
  }
}
```

#### Upload Single File
```http
POST /api/v1/uploads/single
Content-Type: multipart/form-data
Authorization: Bearer <token>

Body: file (single file)
```

#### Get File
```http
GET /api/v1/uploads/:userId/:filename
Authorization: Bearer <token>
```

#### Delete File
```http
DELETE /api/v1/uploads/:filename
Authorization: Bearer <token>
```

#### Get File Info
```http
GET /api/v1/uploads/info/:filename
Authorization: Bearer <token>
```

### 2. Database Models

The following models now have an `attachments` field:

**Transaction Model:**
```go
Attachments []string `gorm:"type:text[]" json:"attachments"`
```

**SavingsGoal Model:**
```go
Attachments []string `gorm:"type:text[]" json:"attachments"`
```

**SavingsContribution Model:**
```go
Attachments []string `gorm:"type:text[]" json:"attachments"`
```

**FixedDeposit Model:**
```go
Attachments []string `gorm:"type:text[]" json:"attachments"`
```

### 3. File Storage Structure

```
backend/
└── uploads/
    └── {userId}/
        ├── invoice_1234567890_abc12345.pdf
        ├── receipt_1234567891_def67890.jpg
        └── statement_1234567892_ghi12345.pdf
```

## 💻 Frontend Implementation

### 1. Import the Component

```vue
<script setup>
import { FileUpload } from '@/components'
import { ref } from 'vue'

const attachments = ref([])
</script>
```

### 2. Use in Your Form

```vue
<template>
  <div>
    <!-- Your form fields -->
    <div class="mb-3">
      <label>Amount</label>
      <input type="number" v-model="form.amount" class="form-control" />
    </div>

    <!-- File Upload Component -->
    <div class="mb-3">
      <FileUpload
        v-model="attachments"
        label="Attachments"
        :multiple="true"
        :max-files="10"
        :max-size="10485760"
        accepted-types=".jpg,.jpeg,.png,.pdf,.doc,.docx"
        @upload-success="handleUploadSuccess"
        @upload-error="handleUploadError"
        @file-removed="handleFileRemoved"
      />
    </div>

    <button @click="submitForm" class="btn btn-primary">Submit</button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { FileUpload } from '@/components'

const attachments = ref([])
const form = ref({
  amount: 0,
  description: ''
})

const handleUploadSuccess = (files) => {
  console.log('Files uploaded:', files)
}

const handleUploadError = (error) => {
  console.error('Upload error:', error)
}

const handleFileRemoved = (file) => {
  console.log('File removed:', file)
}

const submitForm = async () => {
  // Extract file URLs from attachments
  const attachmentUrls = attachments.value.map(f => f.fileUrl)

  const payload = {
    ...form.value,
    attachments: attachmentUrls
  }

  // Submit to your API
  // await transactionStore.createTransaction(payload)
}
</script>
```

### 3. Component Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `label` | String | 'Upload Files' | Label text |
| `modelValue` | Array | `[]` | v-model binding for uploaded files |
| `multiple` | Boolean | `true` | Allow multiple files |
| `maxFiles` | Number | 10 | Maximum number of files |
| `maxSize` | Number | 10485760 | Max file size in bytes (10MB) |
| `acceptedTypes` | String | '.jpg,.jpeg,...' | Accepted file extensions |
| `required` | Boolean | `false` | Show required star |
| `autoUpload` | Boolean | `true` | Auto-upload on selection |

### 4. Component Events

| Event | Payload | Description |
|-------|---------|-------------|
| `update:modelValue` | Array | Emitted when files change |
| `upload-success` | Array | Emitted on successful upload |
| `upload-error` | Error | Emitted on upload failure |
| `file-removed` | Object | Emitted when file is removed |

### 5. Component Methods (via ref)

```vue
<script setup>
import { ref } from 'vue'

const fileUploadRef = ref(null)

// Clear all files
const clearFiles = () => {
  fileUploadRef.value?.clearFiles()
}

// Manually trigger upload
const uploadFiles = (files) => {
  fileUploadRef.value?.uploadFiles(files)
}
</script>

<template>
  <FileUpload ref="fileUploadRef" v-model="attachments" />
</template>
```

## 📝 Example Implementations

### Example 1: Transaction Form

```vue
<template>
  <form @submit.prevent="createTransaction">
    <div class="row">
      <div class="col-md-6 mb-3">
        <label class="form-label">Type</label>
        <select v-model="form.type" class="form-select" required>
          <option value="expense">Expense</option>
          <option value="income">Income</option>
        </select>
      </div>

      <div class="col-md-6 mb-3">
        <label class="form-label">Amount</label>
        <input type="number" v-model="form.amount" class="form-control" required />
      </div>

      <div class="col-12 mb-3">
        <label class="form-label">Description</label>
        <textarea v-model="form.description" class="form-control"></textarea>
      </div>

      <div class="col-12 mb-3">
        <FileUpload
          v-model="attachments"
          label="Receipt/Invoice Attachments"
          :multiple="true"
          :max-files="5"
          accepted-types=".jpg,.jpeg,.png,.pdf"
        />
      </div>

      <div class="col-12">
        <button type="submit" class="btn btn-primary">Create Transaction</button>
      </div>
    </div>
  </form>
</template>

<script setup>
import { ref } from 'vue'
import { FileUpload } from '@/components'
import axios from 'axios'

const form = ref({
  type: 'expense',
  amount: 0,
  description: ''
})
const attachments = ref([])

const createTransaction = async () => {
  try {
    const payload = {
      ...form.value,
      attachments: attachments.value.map(f => f.fileUrl)
    }

    const token = localStorage.getItem('token')
    const response = await axios.post('/api/v1/transactions', payload, {
      headers: { Authorization: `Bearer ${token}` }
    })

    console.log('Transaction created:', response.data)
    // Reset form
    form.value = { type: 'expense', amount: 0, description: '' }
    attachments.value = []
  } catch (error) {
    console.error('Error:', error)
  }
}
</script>
```

### Example 2: Fixed Deposit Form

```vue
<template>
  <div class="card">
    <div class="card-header">
      <h5>Create Fixed Deposit</h5>
    </div>
    <div class="card-body">
      <form @submit.prevent="createFD">
        <div class="mb-3">
          <label class="form-label">Institution</label>
          <input v-model="form.institution" class="form-control" required />
        </div>

        <div class="row">
          <div class="col-md-6 mb-3">
            <label class="form-label">Principal Amount</label>
            <input type="number" v-model="form.principal" class="form-control" required />
          </div>
          <div class="col-md-6 mb-3">
            <label class="form-label">Interest Rate (%)</label>
            <input type="number" step="0.01" v-model="form.interestRate" class="form-control" required />
          </div>
        </div>

        <div class="mb-3">
          <FileUpload
            v-model="fdAttachments"
            label="FD Certificate & Documents"
            :multiple="true"
            :max-files="3"
            accepted-types=".pdf,.jpg,.jpeg,.png"
          />
        </div>

        <button type="submit" class="btn btn-primary">Create Fixed Deposit</button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { FileUpload } from '@/components'

const form = ref({
  institution: '',
  principal: 0,
  interestRate: 0
})
const fdAttachments = ref([])

const createFD = async () => {
  const payload = {
    ...form.value,
    attachments: fdAttachments.value.map(f => f.fileUrl)
  }

  // Submit to API
  console.log('Creating FD:', payload)
}
</script>
```

### Example 3: Viewing Attachments

```vue
<template>
  <div class="attachments-viewer">
    <h6>Attachments ({{ transaction.attachments?.length || 0 }})</h6>

    <div v-if="transaction.attachments?.length" class="row g-2">
      <div v-for="(url, index) in transaction.attachments" :key="index" class="col-md-4">
        <div class="attachment-card">
          <img
            v-if="isImage(url)"
            :src="url"
            alt="Attachment"
            class="attachment-img"
            @click="viewAttachment(url)"
          />
          <div v-else class="attachment-file" @click="viewAttachment(url)">
            <i class="bi bi-file-earmark"></i>
            <span>{{ getFileName(url) }}</span>
          </div>
        </div>
      </div>
    </div>

    <p v-else class="text-muted">No attachments</p>
  </div>
</template>

<script setup>
const props = defineProps({
  transaction: Object
})

const isImage = (url) => {
  return /\.(jpg|jpeg|png|gif)$/i.test(url)
}

const getFileName = (url) => {
  return url.split('/').pop()
}

const viewAttachment = (url) => {
  window.open(url, '_blank')
}
</script>

<style scoped>
.attachment-card {
  border: 1px solid #e3e8ee;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s ease;
}

.attachment-card:hover {
  border-color: #635bff;
  box-shadow: 0 2px 8px rgba(99, 91, 255, 0.2);
}

.attachment-img {
  width: 100%;
  height: 150px;
  object-fit: cover;
}

.attachment-file {
  padding: 24px;
  text-align: center;
  background-color: #f6f9fc;
}

.attachment-file i {
  font-size: 32px;
  color: #8898aa;
  margin-bottom: 8px;
}
</style>
```

## 🚀 Getting Started

### 1. Backend Setup

The backend is ready to use! The upload handler and routes are already configured. Just make sure:

1. The `uploads` directory exists (already created with `.gitkeep`)
2. The database migrations have run (auto-migration enabled)

### 2. Frontend Integration

To integrate file uploads in your forms:

1. **Import the component:**
   ```js
   import { FileUpload } from '@/components'
   ```

2. **Add to your template:**
   ```vue
   <FileUpload v-model="attachments" label="Attachments" />
   ```

3. **Include in form submission:**
   ```js
   const payload = {
     ...formData,
     attachments: attachments.value.map(f => f.fileUrl)
   }
   ```

## 🎨 Styling

The component is fully styled and supports:
- ✅ Light mode
- ✅ Dark mode
- ✅ Mobile responsive (app-like design on mobile)
- ✅ Touch-optimized for mobile devices

## 🔐 Security

- ✅ JWT authentication required for all operations
- ✅ Users can only access their own files
- ✅ File type validation (backend + frontend)
- ✅ File size limits enforced
- ✅ Unique filename generation prevents conflicts
- ✅ Path traversal protection

## 📱 Mobile Support

The component is fully optimized for mobile devices:
- Touch-friendly drag and drop
- Proper tap targets (44px+)
- Responsive layout
- Compressed upload area for small screens
- Image preview optimization

## 🐛 Troubleshooting

### Files not uploading?
1. Check JWT token is valid
2. Verify file size < 10MB
3. Check file type is allowed
4. Ensure backend server is running

### Can't see uploaded files?
1. Check the `attachments` array in your model
2. Verify fileUrl is being saved correctly
3. Check browser console for errors

### Permission denied?
1. Ensure `uploads/` directory has write permissions
2. Check JWT token is included in requests

## 📚 Additional Resources

- Backend handler: `/backend/handlers/upload_handler.go`
- Frontend component: `/frontend/src/components/FileUpload.vue`
- Route definitions: `/backend/routes/routes.go`
- Models: `/backend/models/*.go`

## 🎉 Next Steps

1. Test the upload feature in your development environment
2. Integrate it into your forms (transactions, savings, FD, etc.)
3. Customize the accepted file types as needed
4. Add any additional validation rules
5. Consider adding cloud storage (AWS S3, Google Cloud Storage) for production

---

**Need Help?** Refer to the component props and events sections above, or check the example implementations.
