    <!-- Budgets List -->
    <div class="row g-3">
      <div v-for="budgetProgress in budgets" :key="getBudgetData(budgetProgress).id" class="col-12 col-md-6 col-lg-4">
        <div class="card">
          <div class="card-body">
            <h5 class="card-title">{{ getCategoryName(getBudgetData(budgetProgress).categoryId) }}</h5>
            <p class="text-muted mb-3">{{ getBudgetData(budgetProgress).period }}</p>

            <div class="mb-3">
              <div class="d-flex justify-content-between mb-1">
                <span>{{ formatCurrency(budgetProgress.totalSpent || 0) }}</span>
                <span>{{ formatCurrency(getBudgetData(budgetProgress).amount) }}</span>
              </div>
              <div class="progress" style="height: 12px;">
                <div
                  class="progress-bar"
                  :class="getProgressClass(budgetProgress)"
                  :style="{ width: Math.min(budgetProgress.percentageUsed || 0, 100) + '%' }"
                ></div>
              </div>
              <small class="text-muted">
                {{ Math.round(budgetProgress.percentageUsed || 0) }}% used
              </small>
            </div>

            <div class="mb-2">
              <small class="text-muted">
                Remaining: {{ formatCurrency(budgetProgress.remaining || 0) }}
              </small>
            </div>

            <div v-if="budgetProgress.alertTriggered || budgetProgress.isOverBudget" class="mb-2">
              <span v-if="budgetProgress.isOverBudget" class="badge bg-danger">
                ⛔ Over Budget
              </span>
              <span v-else-if="budgetProgress.alertTriggered" class="badge bg-warning text-dark">
                ⚠️ Alert
              </span>
            </div>

            <div class="d-flex justify-content-between">
              <button class="btn btn-sm btn-outline-primary" @click="editBudget(getBudgetData(budgetProgress))">Edit</button>
              <button class="btn btn-sm btn-danger" @click="deleteBudget(getBudgetData(budgetProgress).id)">Delete</button>
            </div>
          </div>
        </div>
      </div>
    </div>
