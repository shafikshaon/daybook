<template>
  <div class="budgets-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Budgets</h1>
      <button class="btn btn-primary" @click="showAddModal = true">+ Add Budget</button>
    </div>

    <!-- Budget Summary -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon purple">💰</div>
          <div class="stat-value">{{ formatCurrency(budgetsStore.totalBudgeted) }}</div>
          <div class="stat-label">Total Budgeted</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon red">📉</div>
          <div class="stat-value">{{ formatCurrency(budgetsStore.totalSpent) }}</div>
          <div class="stat-label">Total Spent</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon green">💵</div>
          <div class="stat-value">{{ formatCurrency(budgetsStore.totalBudgeted - budgetsStore.totalSpent) }}</div>
          <div class="stat-label">Remaining</div>
        </div>
      </div>
    </div>

    <!-- Budget Alerts -->
    <div v-if="budgetAlerts.length > 0" class="alert alert-warning mb-4">
      <h6 class="alert-heading">⚠️ Budget Alerts</h6>
      <ul class="mb-0">
        <li v-for="alert in budgetAlerts" :key="alert.id">
          {{ alert.message }}: {{ formatCurrency(alert.amount) }} / {{ formatCurrency(alert.budget) }}
        </li>
      </ul>
    </div>

    <!-- Debug Info -->
    <div v-if="budgets.length === 0" class="alert alert-info mb-4">
      <p class="mb-0">No budgets found. Click "Add Budget" to create your first budget.</p>
    </div>

    <!-- Budgets List -->
    <div class="row g-3">
      <div v-for="item in budgets" :key="(item.budget || item).id" class="col-12 col-md-6 col-lg-4">
        <div class="card">
          <div class="card-body">
            <h5 class="card-title">{{ getCategoryName((item.budget || item).categoryId) }}</h5>
            <p class="text-muted mb-3">{{ (item.budget || item).period }}</p>

            <div class="mb-3">
              <div class="d-flex justify-content-between mb-1">
                <span>{{ formatCurrency(item.totalSpent || 0) }}</span>
                <span>{{ formatCurrency((item.budget || item).amount) }}</span>
              </div>
              <div class="progress" style="height: 12px;">
                <div
                  class="progress-bar"
                  :class="getProgressClass(item)"
                  :style="{ width: Math.min(item.percentageUsed || 0, 100) + '%' }"
                ></div>
              </div>
              <small class="text-muted">
                {{ Math.round(item.percentageUsed || 0) }}% used
              </small>
            </div>

            <div class="mb-2">
              <small class="text-muted">
                Remaining: {{ formatCurrency(item.remaining || 0) }}
              </small>
            </div>

            <div v-if="item.alertTriggered || item.isOverBudget" class="mb-2">
              <span v-if="item.isOverBudget" class="badge bg-danger">
                ⛔ Over Budget
              </span>
              <span v-else-if="item.alertTriggered" class="badge bg-warning text-dark">
                ⚠️ Alert
              </span>
            </div>

            <div class="d-flex justify-content-between">
              <button class="btn btn-sm btn-outline-primary" @click="editBudget(item.budget || item)">Edit</button>
              <button class="btn btn-sm btn-danger" @click="deleteBudget((item.budget || item).id)">Delete</button>
            </div>
          </div>
        </div>
      </div>
    </div>
