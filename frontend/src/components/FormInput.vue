<template>
  <div class="mb-3">
    <label v-if="label" class="form-label">
      {{ label }}
      <span v-if="required" class="text-danger">*</span>
    </label>

    <!-- Text/Number/Date inputs -->
    <input
      v-if="['text', 'number', 'email', 'password', 'date', 'time', 'datetime-local'].includes(type)"
      :type="type"
      class="form-control"
      :class="{ 'is-invalid': error }"
      :value="modelValue"
      @input="handleInput"
      @blur="handleBlur"
      :placeholder="placeholder"
      :required="required"
      :disabled="disabled"
      :readonly="readonly"
      :min="min"
      :max="max"
      :step="step"
    />

    <!-- Textarea -->
    <textarea
      v-else-if="type === 'textarea'"
      class="form-control"
      :class="{ 'is-invalid': error }"
      :value="modelValue"
      @input="handleInput"
      @blur="handleBlur"
      :placeholder="placeholder"
      :required="required"
      :disabled="disabled"
      :readonly="readonly"
      :rows="rows"
    ></textarea>

    <!-- Select -->
    <select
      v-else-if="type === 'select'"
      class="form-select"
      :class="{ 'is-invalid': error }"
      :value="modelValue"
      @change="handleInput"
      @blur="handleBlur"
      :required="required"
      :disabled="disabled"
    >
      <option value="" v-if="placeholder">{{ placeholder }}</option>
      <option
        v-for="option in options"
        :key="option.value"
        :value="option.value"
      >
        {{ option.label }}
      </option>
    </select>

    <!-- Checkbox -->
    <div v-else-if="type === 'checkbox'" class="form-check">
      <input
        type="checkbox"
        class="form-check-input"
        :class="{ 'is-invalid': error }"
        :checked="modelValue"
        @change="handleCheckbox"
        :required="required"
        :disabled="disabled"
      />
      <label class="form-check-label">{{ checkboxLabel }}</label>
    </div>

    <!-- Color picker -->
    <input
      v-else-if="type === 'color'"
      type="color"
      class="form-control form-control-color"
      :class="{ 'is-invalid': error }"
      :value="modelValue"
      @input="handleInput"
      :disabled="disabled"
    />

    <div v-if="error" class="invalid-feedback d-block">
      {{ error }}
    </div>
    <small v-else-if="helpText" class="form-text text-muted">
      {{ helpText }}
    </small>
  </div>
</template>

<script setup>
defineProps({
  modelValue: {
    type: [String, Number, Boolean],
    default: ''
  },
  type: {
    type: String,
    default: 'text'
  },
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  required: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  readonly: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  },
  helpText: {
    type: String,
    default: ''
  },
  options: {
    type: Array,
    default: () => []
  },
  rows: {
    type: Number,
    default: 3
  },
  checkboxLabel: {
    type: String,
    default: ''
  },
  min: {
    type: [String, Number],
    default: undefined
  },
  max: {
    type: [String, Number],
    default: undefined
  },
  step: {
    type: [String, Number],
    default: undefined
  }
})

const emit = defineEmits(['update:modelValue', 'blur'])

const handleInput = (event) => {
  let value = event.target.value
  if (event.target.type === 'number') {
    value = value === '' ? null : Number(value)
  }
  emit('update:modelValue', value)
}

const handleCheckbox = (event) => {
  emit('update:modelValue', event.target.checked)
}

const handleBlur = () => {
  emit('blur')
}
</script>
