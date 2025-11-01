<template>
  <div class="transactions-view fade-in">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1 class="text-purple">Transactions</h1>
      <div class="d-flex gap-2">
        <button class="btn btn-outline-primary" @click="showTransferModal = true">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
            <path fill-rule="evenodd" d="M8 3a5 5 0 1 0 4.546 2.914.5.5 0 0 1 .908-.417A6 6 0 1 1 8 2z"/>
            <path d="M8 4.466V.534a.25.25 0 0 1 .41-.192l2.36 1.966c.12.1.12.284 0 .384L8.41 4.658A.25.25 0 0 1 8 4.466"/>
          </svg>
          Transfer Funds
        </button>
        <button class="btn btn-primary" @click="showAddModal = true">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
            <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4"/>
          </svg>
          Add Transaction
        </button>
      </div>
    </div>

    <!-- Summary Stats -->
    <div class="row g-3 mb-4">
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon green">📈</div>
          <div class="stat-value">{{ formatCurrency(totalIncome) }}</div>
          <div class="stat-label">Total Income</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon red">📉</div>
          <div class="stat-value">{{ formatCurrency(totalExpense) }}</div>
          <div class="stat-label">Total Expenses</div>
        </div>
      </div>
      <div class="col-12 col-md-4">
        <div class="stat-card">
          <div class="stat-icon blue">💰</div>
          <div class="stat-value">{{ formatCurrency(totalIncome - totalExpense) }}</div>
          <div class="stat-label">Net Cash Flow</div>
          <div :class="totalIncome - totalExpense >= 0 ? 'stat-change positive' : 'stat-change negative'">
            {{ totalIncome - totalExpense >= 0 ? '↑' : '↓' }}
          </div>
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="card mb-4">
      <div class="card-body">
        <div class="row g-3 mb-3">
          <div class="col-12 col-md-3">
            <label class="form-label">Type</label>
            <select class="form-select" v-model="filters.type" @change="applyFilters">
              <option value="">All</option>
              <option value="income">Income</option>
              <option value="expense">Expense</option>
              <option value="transfer">Transfer</option>
            </select>
          </div>
          <div class="col-12 col-md-3">
            <label class="form-label">Category</label>
            <select class="form-select" v-model="filters.category" @change="applyFilters">
              <option value="">All Categories</option>
              <option v-for="cat in allCategories" :key="cat.id" :value="cat.id">
                {{ cat.name }}
              </option>
            </select>
          </div>
          <div class="col-12 col-md-3">
            <label class="form-label">Account</label>
            <select class="form-select" v-model="filters.account" @change="applyFilters">
              <option value="">All Accounts</option>
              <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
                {{ acc.name }}
              </option>
            </select>
          </div>
          <div class="col-12 col-md-3">
            <label class="form-label">Search</label>
            <input
              type="text"
              class="form-control"
              v-model="filters.search"
              placeholder="Search transactions..."
              @input="applyFilters"
            />
          </div>
        </div>
        <div class="row g-3 align-items-end">
          <div class="col-12 col-md-3">
            <label class="form-label">Start Date</label>
            <input
              type="date"
              class="form-control"
              v-model="filters.startDate"
              @change="applyFilters"
            />
          </div>
          <div class="col-12 col-md-3">
            <label class="form-label">End Date</label>
            <input
              type="date"
              class="form-control"
              v-model="filters.endDate"
              @change="applyFilters"
            />
          </div>
          <div class="col-12 col-md-3">
            <label class="form-label">Quick Filter</label>
            <select class="form-select" v-model="quickDateFilter" @change="applyQuickFilter">
              <option value="">Custom Range</option>
              <option value="today">Today</option>
              <option value="yesterday">Yesterday</option>
              <option value="this_week">This Week</option>
              <option value="last_week">Last Week</option>
              <option value="this_month">This Month</option>
              <option value="last_month">Last Month</option>
              <option value="this_year">This Year</option>
            </select>
          </div>
          <div class="col-12 col-md-3">
            <div class="form-check">
              <input
                class="form-check-input"
                type="checkbox"
                id="groupByDay"
                v-model="groupByDay"
              />
              <label class="form-check-label" for="groupByDay">
                Group by Day
              </label>
            </div>
            <button class="btn btn-sm btn-outline-secondary mt-2" @click="clearFilters">
              Clear Filters
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Transactions Table -->
    <div class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <h5 class="mb-0">Transaction History</h5>
        <div class="d-flex align-items-center gap-2">
          <label class="mb-0 me-2">Items per page:</label>
          <select class="form-select form-select-sm" v-model.number="itemsPerPage" @change="changeItemsPerPage(itemsPerPage)" style="width: auto;">
            <option :value="20">20</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
            <option :value="500">500</option>
          </select>
        </div>
      </div>
      <div class="card-body p-0">
        <div v-if="filteredTransactions.length === 0" class="p-4 text-center text-muted">
          No transactions found
        </div>
        <div v-else-if="groupByDay" class="grouped-transactions">
          <div v-for="group in groupedTransactions" :key="group.date" class="day-group">
            <div class="day-header">
              <div class="day-date">
                <strong>{{ formatGroupDate(group.date) }}</strong>
                <span class="text-muted ms-2">({{ group.transactions.length }} transactions)</span>
              </div>
              <div class="day-summary">
                <span class="text-muted me-3">Income: <strong class="amount-income">{{ formatCurrency(group.totalIncome) }}</strong></span>
                <span class="text-muted">Expense: <strong class="amount-expense">{{ formatCurrency(group.totalExpense) }}</strong></span>
              </div>
            </div>
            <div class="table-responsive">
              <table class="table table-hover mb-0">
                <tbody>
                  <tr v-for="transaction in group.transactions" :key="transaction.id">
                    <td style="width: 35%">
                      <div>{{ transaction.description || '-' }}</div>
                      <small v-if="transaction.creditCardId" class="text-muted">
                        💳 {{ getCreditCardName(transaction.creditCardId) }}
                      </small>
                    </td>
                    <td style="width: 15%">
                      <span class="badge" :style="{ backgroundColor: getCategoryColor(transaction.categoryId) }">
                        {{ getCategoryName(transaction.categoryId) }}
                      </span>
                    </td>
                    <td style="width: 20%">
                      <span v-if="transaction.type === 'transfer'">
                        {{ transaction.accountName || getAccountName(transaction.accountId) }} → {{ transaction.toAccountName || getAccountName(transaction.toAccountId) }}
                      </span>
                      <span v-else>{{ transaction.accountName || getAccountName(transaction.accountId) }}</span>
                    </td>
                    <td style="width: 10%">
                      <span
                        class="badge transaction-type-badge"
                        :class="{
                          'badge-income': transaction.type === 'income',
                          'badge-expense': transaction.type === 'expense',
                          'badge-transfer': transaction.type === 'transfer'
                        }"
                      >
                        {{ transaction.type }}
                      </span>
                    </td>
                    <td class="text-end" style="width: 10%">
                      <span
                        class="fw-bold transaction-amount"
                        :class="{
                          'amount-income': transaction.type === 'income',
                          'amount-expense': transaction.type === 'expense',
                          'amount-transfer': transaction.type === 'transfer'
                        }"
                      >
                        {{ transaction.type === 'income' ? '+' : transaction.type === 'transfer' ? '↔' : '-' }}{{ formatCurrency(transaction.amount) }}
                      </span>
                    </td>
                    <td class="text-center" style="width: 5%">
                      <span v-if="transaction.attachments && transaction.attachments.length > 0" class="badge bg-info" style="cursor: pointer;" @click="viewAttachments(transaction)" :title="`${transaction.attachments.length} file(s)`">
                        📎 {{ transaction.attachments.length }}
                      </span>
                      <span v-else class="text-muted">-</span>
                    </td>
                    <td class="text-center" style="width: 5%">
                      <button
                        class="btn btn-sm btn-primary me-1"
                        @click="editTransaction(transaction)"
                      >
                        Edit
                      </button>
                      <button
                        class="btn btn-sm btn-danger"
                        @click="confirmDelete(transaction)"
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
        <div v-else class="table-responsive">
          <table class="table table-hover mb-0">
            <thead>
              <tr>
                <th>Date</th>
                <th>Description</th>
                <th>Category</th>
                <th>Account</th>
                <th>Type</th>
                <th class="text-end">Amount</th>
                <th class="text-center">Files</th>
                <th class="text-center">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="transaction in paginatedTransactions" :key="transaction.id">
                <td>{{ formatDate(transaction.date) }}</td>
                <td>
                  <div>{{ transaction.description || '-' }}</div>
                  <small v-if="transaction.creditCardId" class="text-muted">
                    💳 {{ getCreditCardName(transaction.creditCardId) }}
                  </small>
                </td>
                <td>
                  <span class="badge" :style="{ backgroundColor: getCategoryColor(transaction.categoryId) }">
                    {{ getCategoryName(transaction.categoryId) }}
                  </span>
                </td>
                <td>
                  <span v-if="transaction.type === 'transfer'">
                    {{ transaction.accountName || getAccountName(transaction.accountId) }} → {{ transaction.toAccountName || getAccountName(transaction.toAccountId) }}
                  </span>
                  <span v-else>{{ transaction.accountName || getAccountName(transaction.accountId) }}</span>
                </td>
                <td>
                  <span
                    class="badge transaction-type-badge"
                    :class="{
                      'badge-income': transaction.type === 'income',
                      'badge-expense': transaction.type === 'expense',
                      'badge-transfer': transaction.type === 'transfer'
                    }"
                  >
                    {{ transaction.type }}
                  </span>
                </td>
                <td class="text-end">
                  <span
                    class="fw-bold transaction-amount"
                    :class="{
                      'amount-income': transaction.type === 'income',
                      'amount-expense': transaction.type === 'expense',
                      'amount-transfer': transaction.type === 'transfer'
                    }"
                  >
                    {{ transaction.type === 'income' ? '+' : transaction.type === 'transfer' ? '↔' : '-' }}{{ formatCurrency(transaction.amount) }}
                  </span>
                </td>
                <td class="text-center">
                  <span v-if="transaction.attachments && transaction.attachments.length > 0" class="badge bg-info" style="cursor: pointer;" @click="viewAttachments(transaction)" :title="`${transaction.attachments.length} file(s)`">
                    📎 {{ transaction.attachments.length }}
                  </span>
                  <span v-else class="text-muted">-</span>
                </td>
                <td class="text-center">
                  <button
                    class="btn btn-sm btn-primary me-1"
                    @click="editTransaction(transaction)"
                  >
                    Edit
                  </button>
                  <button
                    class="btn btn-sm btn-danger"
                    @click="confirmDelete(transaction)"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <div class="card-footer" v-if="pagination.totalPages > 1">
        <div class="d-flex justify-content-between align-items-center">
          <div class="pagination-info">
            Showing {{ (pagination.currentPage - 1) * pagination.limit + 1 }} to {{ Math.min(pagination.currentPage * pagination.limit, pagination.totalCount) }} of {{ pagination.totalCount }} transactions
          </div>
          <nav aria-label="Transaction pagination">
            <ul class="pagination pagination-sm mb-0">
              <li class="page-item" :class="{ disabled: !pagination.hasPrev }">
                <a class="page-link" href="#" @click.prevent="changePage(currentPage - 1)">Previous</a>
              </li>
              <li
                v-for="page in visiblePages"
                :key="page"
                class="page-item"
                :class="{ active: page === currentPage }"
              >
                <a class="page-link" href="#" @click.prevent="changePage(page)">{{ page }}</a>
              </li>
              <li class="page-item" :class="{ disabled: !pagination.hasNext }">
                <a class="page-link" href="#" @click.prevent="changePage(currentPage + 1)">Next</a>
              </li>
            </ul>
          </nav>
        </div>
      </div>
    </div>

    <!-- Add/Edit Transaction Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showAddModal || showEditModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showAddModal || showEditModal"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ showEditModal ? 'Edit Transaction' : 'Add Transaction' }}</h5>
            <button type="button" class="btn-close" @click="closeModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveTransaction">
              <div class="mb-3">
                <label class="form-label">Type *</label>
                <select class="form-select" v-model="form.type" required>
                  <option value="income">Income</option>
                  <option value="expense">Expense</option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Amount *</label>
                <input
                  type="number"
                  step="0.01"
                  class="form-control"
                  v-model.number="form.amount"
                  required
                  placeholder="0.00"
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Category *</label>
                <select class="form-select" v-model="form.categoryId" required>
                  <option value="">Select category...</option>
                  <option
                    v-for="cat in filteredCategories"
                    :key="cat.id"
                    :value="cat.id"
                  >
                    {{ cat.icon }} {{ cat.name }}
                  </option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Account *</label>
                <select class="form-select" v-model="form.accountId" required>
                  <option value="">Select account...</option>
                  <option v-for="acc in accountsAndCreditCards" :key="acc.id" :value="acc.id">
                    {{ acc.name }}
                  </option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Date *</label>
                <input
                  type="date"
                  class="form-control"
                  v-model="form.date"
                  required
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Description</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="form.description"
                  placeholder="Optional description"
                />
              </div>

              <div class="mb-3">
                <FileUpload
                  v-model="transactionAttachments"
                  label="Attachments (Receipts/Invoices)"
                  :multiple="true"
                  :max-files="5"
                  :max-size="10485760"
                  accepted-types=".jpg,.jpeg,.png,.pdf,.doc,.docx"
                />
              </div>

              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="closeModal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary">
                  {{ showEditModal ? 'Update' : 'Create' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Transfer Funds Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showTransferModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showTransferModal"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Transfer Funds</h5>
            <button type="button" class="btn-close" @click="closeTransferModal"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveTransfer">
              <div class="mb-3">
                <label class="form-label">From Account *</label>
                <select class="form-select" v-model="transferForm.fromAccountId" required>
                  <option value="">Select account...</option>
                  <option v-for="acc in accounts" :key="acc.id" :value="acc.id">
                    {{ acc.name }} ({{ formatCurrency(acc.balance) }})
                  </option>
                </select>
              </div>

              <div class="mb-3">
                <label class="form-label">Transfer To *</label>
                <select class="form-select" v-model="transferForm.destinationType" required>
                  <option value="account">Another Account</option>
                  <option value="savings_goal">Savings Goal</option>
                  <option value="fixed_deposit">Fixed Deposit (Create New)</option>
                </select>
              </div>

              <div class="mb-3" v-if="transferForm.destinationType === 'account'">
                <label class="form-label">To Account *</label>
                <select class="form-select" v-model="transferForm.toAccountId" required>
                  <option value="">Select account...</option>
                  <option
                    v-for="acc in accounts"
                    :key="acc.id"
                    :value="acc.id"
                    :disabled="acc.id === transferForm.fromAccountId"
                  >
                    {{ acc.name }} ({{ formatCurrency(acc.balance) }})
                  </option>
                </select>
              </div>

              <div class="mb-3" v-if="transferForm.destinationType === 'savings_goal'">
                <label class="form-label">Savings Goal *</label>
                <select class="form-select" v-model="transferForm.savingsGoalId" required>
                  <option value="">Select savings goal...</option>
                  <option v-for="goal in savingsGoals" :key="goal.id" :value="goal.id">
                    {{ goal.name }} ({{ formatCurrency(goal.currentAmount) }} / {{ formatCurrency(goal.targetAmount) }})
                  </option>
                </select>
              </div>

              <div class="mb-3" v-if="transferForm.destinationType === 'fixed_deposit'">
                <label class="form-label">Fixed Deposit Name *</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="transferForm.fdName"
                  required
                  placeholder="e.g., 1-Year FD"
                />
              </div>

              <div class="mb-3" v-if="transferForm.destinationType === 'fixed_deposit'">
                <label class="form-label">Interest Rate (% per annum) *</label>
                <input
                  type="number"
                  step="0.01"
                  class="form-control"
                  v-model.number="transferForm.fdInterestRate"
                  required
                  placeholder="e.g., 5.5"
                />
              </div>

              <div class="mb-3" v-if="transferForm.destinationType === 'fixed_deposit'">
                <label class="form-label">Tenure (months) *</label>
                <input
                  type="number"
                  class="form-control"
                  v-model.number="transferForm.fdTenureMonths"
                  required
                  placeholder="e.g., 12"
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Amount *</label>
                <input
                  type="number"
                  step="0.01"
                  class="form-control"
                  v-model.number="transferForm.amount"
                  required
                  placeholder="0.00"
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Date *</label>
                <input
                  type="date"
                  class="form-control"
                  v-model="transferForm.date"
                  required
                />
              </div>

              <div class="mb-3">
                <label class="form-label">Description</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="transferForm.description"
                  placeholder="Optional description"
                />
              </div>

              <div class="alert alert-info">
                <small v-if="transferForm.destinationType === 'account'">
                  This will transfer {{ formatCurrency(transferForm.amount || 0) }} from {{ getAccountName(transferForm.fromAccountId) || 'source' }} to {{ getAccountName(transferForm.toAccountId) || 'destination' }}.
                </small>
                <small v-else-if="transferForm.destinationType === 'savings_goal'">
                  This will transfer {{ formatCurrency(transferForm.amount || 0) }} from {{ getAccountName(transferForm.fromAccountId) || 'source' }} to your savings goal.
                </small>
                <small v-else-if="transferForm.destinationType === 'fixed_deposit'">
                  This will create a new fixed deposit of {{ formatCurrency(transferForm.amount || 0) }} for {{ transferForm.fdTenureMonths || 0 }} months at {{ transferForm.fdInterestRate || 0 }}% interest.
                </small>
              </div>

              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-secondary" @click="closeTransferModal">
                  Cancel
                </button>
                <button type="submit" class="btn btn-primary">
                  Transfer
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- View Attachments Modal -->
    <div
      class="modal fade"
      :class="{ 'show d-block': showAttachmentsModal }"
      style="background-color: rgba(0,0,0,0.5);"
      v-if="showAttachmentsModal"
    >
      <div class="modal-dialog modal-lg modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Transaction Attachments</h5>
            <button type="button" class="btn-close" @click="showAttachmentsModal = false"></button>
          </div>
          <div class="modal-body">
            <div v-if="viewingTransaction && viewingTransaction.attachments && viewingTransaction.attachments.length > 0">
              <div class="row g-3">
                <div
                  v-for="(url, index) in viewingTransaction.attachments"
                  :key="index"
                  class="col-md-6 col-lg-4"
                >
                  <div class="attachment-card">
                    <div v-if="isImageUrl(url)" class="attachment-preview">
                      <img :src="url" :alt="`Attachment ${index + 1}`" class="img-fluid" @click="openAttachment(url)" style="cursor: pointer;" />
                    </div>
                    <div v-else class="attachment-file-icon" @click="openAttachment(url)" style="cursor: pointer;">
                      <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" fill="currentColor" viewBox="0 0 16 16">
                        <path d="M5 4a.5.5 0 0 0 0 1h6a.5.5 0 0 0 0-1zm-.5 2.5A.5.5 0 0 1 5 6h6a.5.5 0 0 1 0 1H5a.5.5 0 0 1-.5-.5M5 8a.5.5 0 0 0 0 1h6a.5.5 0 0 0 0-1zm0 2a.5.5 0 0 0 0 1h3a.5.5 0 0 0 0-1z"/>
                        <path d="M2 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2zm10-1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1"/>
                      </svg>
                    </div>
                    <div class="attachment-info">
                      <p class="mb-1 small text-truncate" :title="getFileNameFromUrl(url)">
                        {{ getFileNameFromUrl(url) }}
                      </p>
                      <a :href="url" target="_blank" class="btn btn-sm btn-outline-primary">
                        Download
                      </a>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="text-center text-muted py-4">
              No attachments available
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountsStore } from '@/stores/accounts'
import { useCreditCardsStore } from '@/stores/creditCards'
import { useSettingsStore } from '@/stores/settings'
import { useNotification } from '@/composables/useNotification'
import { FileUpload } from '@/components'

const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()
const creditCardsStore = useCreditCardsStore()
const settingsStore = useSettingsStore()
const { confirm, success, error } = useNotification()

const showAddModal = ref(false)
const showEditModal = ref(false)
const showTransferModal = ref(false)
const showAttachmentsModal = ref(false)
const editingTransaction = ref(null)
const viewingTransaction = ref(null)

const filters = ref({
  type: '',
  category: '',
  account: '',
  search: '',
  startDate: '',
  endDate: ''
})

const quickDateFilter = ref('')
const groupByDay = ref(false)
const currentPage = ref(1)
const itemsPerPage = ref(20)

const form = ref({
  type: 'expense',
  amount: 0,
  categoryId: '',
  accountId: '',
  date: new Date().toISOString().split('T')[0],
  description: '',
  tags: [],
  attachments: []
})

const transactionAttachments = ref([])

const transferForm = ref({
  fromAccountId: '',
  destinationType: 'account',
  toAccountId: '',
  savingsGoalId: '',
  fdName: '',
  fdInterestRate: 0,
  fdTenureMonths: 12,
  amount: 0,
  date: new Date().toISOString().split('T')[0],
  description: 'Transfer between accounts'
})

const transactions = computed(() => transactionsStore.allTransactions)
const accounts = computed(() => accountsStore.allAccounts)
const allCategories = computed(() => transactionsStore.categories)
const pagination = computed(() => transactionsStore.pagination)

// Combined accounts and credit cards for dropdowns
const accountsAndCreditCards = computed(() => {
  const regularAccounts = accounts.value.map(acc => ({
    id: acc.id,
    name: acc.name,
    type: 'account',
    balance: acc.balance
  }))

  const creditCardAccounts = creditCardsStore.allCreditCards.map(card => ({
    id: card.id,
    name: `💳 ${card.name}`,
    type: 'credit_card',
    balance: card.currentBalance
  }))

  return [...regularAccounts, ...creditCardAccounts]
})

const filteredCategories = computed(() => {
  return transactionsStore.categories.filter(c => c.type === form.value.type)
})

const totalIncome = computed(() => transactionsStore.totalIncome())
const totalExpense = computed(() => transactionsStore.totalExpense())

const filteredTransactions = computed(() => {
  let result = transactions.value

  // Client-side search filter only (other filters handled by backend)
  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    result = result.filter(t =>
      (t.description || '').toLowerCase().includes(search) ||
      getCategoryName(t.categoryId).toLowerCase().includes(search)
    )
  }

  return result
})

const groupedTransactions = computed(() => {
  if (!groupByDay.value) return []

  const groups = {}

  filteredTransactions.value.forEach(transaction => {
    const dateKey = new Date(transaction.date).toISOString().split('T')[0]

    if (!groups[dateKey]) {
      groups[dateKey] = {
        date: dateKey,
        transactions: [],
        totalIncome: 0,
        totalExpense: 0
      }
    }

    groups[dateKey].transactions.push(transaction)

    if (transaction.type === 'income') {
      groups[dateKey].totalIncome += transaction.amount
    } else if (transaction.type === 'expense') {
      groups[dateKey].totalExpense += transaction.amount
    }
  })

  // Convert to array and sort by date descending
  return Object.values(groups).sort((a, b) => new Date(b.date) - new Date(a.date))
})

const paginatedTransactions = computed(() => {
  // Return transactions directly from API (already paginated)
  return filteredTransactions.value
})

const totalPages = computed(() => pagination.value.totalPages || 1)

const changePage = async (page) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  await loadTransactions()
}

const changeItemsPerPage = async (newLimit) => {
  itemsPerPage.value = newLimit
  currentPage.value = 1 // Reset to first page
  await loadTransactions()
}

const loadTransactions = async () => {
  try {
    const params = {
      page: currentPage.value,
      limit: itemsPerPage.value
    }

    // Add filters to params
    if (filters.value.type) params.type = filters.value.type
    if (filters.value.category) params.categoryId = filters.value.category
    if (filters.value.account) params.accountId = filters.value.account
    if (filters.value.startDate) params.startDate = filters.value.startDate
    if (filters.value.endDate) params.endDate = filters.value.endDate

    await transactionsStore.fetchTransactions(params.page, params.limit, params)
  } catch (err) {
    error('Error loading transactions')
  }
}

const applyFilters = () => {
  currentPage.value = 1 // Reset to first page when filters change
  loadTransactions()
}

const applyQuickFilter = () => {
  const today = new Date()
  const value = quickDateFilter.value

  if (!value) {
    return
  }

  switch (value) {
    case 'today':
      filters.value.startDate = today.toISOString().split('T')[0]
      filters.value.endDate = today.toISOString().split('T')[0]
      break
    case 'yesterday':
      const yesterday = new Date(today)
      yesterday.setDate(yesterday.getDate() - 1)
      filters.value.startDate = yesterday.toISOString().split('T')[0]
      filters.value.endDate = yesterday.toISOString().split('T')[0]
      break
    case 'this_week':
      const startOfWeek = new Date(today)
      startOfWeek.setDate(today.getDate() - today.getDay())
      filters.value.startDate = startOfWeek.toISOString().split('T')[0]
      filters.value.endDate = today.toISOString().split('T')[0]
      break
    case 'last_week':
      const lastWeekStart = new Date(today)
      lastWeekStart.setDate(today.getDate() - today.getDay() - 7)
      const lastWeekEnd = new Date(lastWeekStart)
      lastWeekEnd.setDate(lastWeekStart.getDate() + 6)
      filters.value.startDate = lastWeekStart.toISOString().split('T')[0]
      filters.value.endDate = lastWeekEnd.toISOString().split('T')[0]
      break
    case 'this_month':
      const startOfMonth = new Date(today.getFullYear(), today.getMonth(), 1)
      filters.value.startDate = startOfMonth.toISOString().split('T')[0]
      filters.value.endDate = today.toISOString().split('T')[0]
      break
    case 'last_month':
      const lastMonthStart = new Date(today.getFullYear(), today.getMonth() - 1, 1)
      const lastMonthEnd = new Date(today.getFullYear(), today.getMonth(), 0)
      filters.value.startDate = lastMonthStart.toISOString().split('T')[0]
      filters.value.endDate = lastMonthEnd.toISOString().split('T')[0]
      break
    case 'this_year':
      const startOfYear = new Date(today.getFullYear(), 0, 1)
      filters.value.startDate = startOfYear.toISOString().split('T')[0]
      filters.value.endDate = today.toISOString().split('T')[0]
      break
    default:
      // Custom range - do nothing
      break
  }

  applyFilters()
}

const clearFilters = () => {
  filters.value = {
    type: '',
    category: '',
    account: '',
    search: '',
    startDate: '',
    endDate: ''
  }
  quickDateFilter.value = ''
  applyFilters()
}

const visiblePages = computed(() => {
  const pages = []
  const total = totalPages.value
  const current = currentPage.value

  // Always show first page
  pages.push(1)

  // Show pages around current page
  for (let i = Math.max(2, current - 1); i <= Math.min(total - 1, current + 1); i++) {
    if (!pages.includes(i)) {
      pages.push(i)
    }
  }

  // Always show last page if there are multiple pages
  if (total > 1 && !pages.includes(total)) {
    pages.push(total)
  }

  return pages.sort((a, b) => a - b)
})

const formatCurrency = (amount) => {
  return settingsStore.formatCurrency(amount)
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

const formatGroupDate = (dateString) => {
  const date = new Date(dateString)
  const today = new Date()
  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)

  const isToday = date.toDateString() === today.toDateString()
  const isYesterday = date.toDateString() === yesterday.toDateString()

  if (isToday) {
    return 'Today, ' + date.toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })
  } else if (isYesterday) {
    return 'Yesterday, ' + date.toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })
  } else {
    return date.toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric', year: 'numeric' })
  }
}

const getCategoryName = (categoryId) => {
  const category = transactionsStore.getCategoryById(categoryId)
  return category ? category.name : categoryId
}

const getCategoryColor = (categoryId) => {
  const category = transactionsStore.getCategoryById(categoryId)
  return category ? category.color : '#6b7280'
}

const getAccountName = (accountId) => {
  const account = accountsStore.getAccountById(accountId)
  return account ? account.name : accountId
}

const getCreditCardName = (cardId) => {
  const card = creditCardsStore.getCreditCardById(cardId)
  return card ? card.name : 'Credit Card'
}

const editTransaction = (transaction) => {
  editingTransaction.value = transaction
  form.value = {
    ...transaction,
    date: new Date(transaction.date).toISOString().split('T')[0],
    // If transaction has creditCardId, use it as accountId for the dropdown
    accountId: transaction.creditCardId || transaction.accountId
  }
  // Load existing attachments
  transactionAttachments.value = (transaction.attachments || []).map(url => ({
    fileUrl: url,
    originalName: url.split('/').pop(),
    fileName: url.split('/').pop()
  }))
  showEditModal.value = true
}

const confirmDelete = async (transaction) => {
  const confirmed = await confirm({
    title: 'Delete Transaction',
    message: `Are you sure you want to delete this transaction? This action cannot be undone.`,
    confirmText: 'Delete',
    cancelText: 'Cancel',
    variant: 'danger'
  })

  if (confirmed) {
    try {
      await transactionsStore.deleteTransaction(transaction.id)
      success('Transaction deleted successfully')
      await loadTransactions()
    } catch (err) {
      error(err.response?.data?.message || err.message || 'Error deleting transaction')
    }
  }
}

const saveTransaction = async () => {
  try {
    // Check if selected account is a credit card
    const isCreditCard = creditCardsStore.allCreditCards.some(card => card.id === form.value.accountId)

    const transactionData = {
      ...form.value,
      date: new Date(form.value.date).toISOString(),
      attachments: transactionAttachments.value.map(f => f.fileUrl)
    }

    // Set the correct ID field based on whether it's a credit card or account
    if (isCreditCard) {
      transactionData.creditCardId = form.value.accountId
      transactionData.accountId = null
    } else {
      transactionData.accountId = form.value.accountId
      transactionData.creditCardId = null
    }

    if (showEditModal.value) {
      await transactionsStore.updateTransaction(editingTransaction.value.id, transactionData)
      success('Transaction updated successfully')
    } else {
      await transactionsStore.createTransaction(transactionData)
      success('Transaction created successfully')
    }

    closeModal()
    await loadTransactions()
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error saving transaction')
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  editingTransaction.value = null
  form.value = {
    type: 'expense',
    amount: 0,
    categoryId: '',
    accountId: '',
    date: new Date().toISOString().split('T')[0],
    description: '',
    tags: [],
    attachments: []
  }
  transactionAttachments.value = []
}

const saveTransfer = async () => {
  try {
    const dateISO = new Date(transferForm.value.date).toISOString()

    if (transferForm.value.destinationType === 'account') {
      if (transferForm.value.fromAccountId === transferForm.value.toAccountId) {
        error('Source and destination accounts must be different')
        return
      }

      await transactionsStore.transferFunds(
        transferForm.value.fromAccountId,
        transferForm.value.toAccountId,
        transferForm.value.amount,
        transferForm.value.description,
        dateISO
      )
      success('Funds transferred successfully')
    } else if (transferForm.value.destinationType === 'savings_goal') {
      if (!transferForm.value.savingsGoalId) {
        error('Please select a savings goal')
        return
      }

      await transactionsStore.transferToSavingsGoal(
        transferForm.value.fromAccountId,
        transferForm.value.savingsGoalId,
        transferForm.value.amount,
        transferForm.value.description,
        dateISO
      )

      success('Funds transferred to savings goal successfully')
    } else if (transferForm.value.destinationType === 'fixed_deposit') {
      if (!transferForm.value.fdName || !transferForm.value.fdInterestRate || !transferForm.value.fdTenureMonths) {
        error('Please fill in all fixed deposit details')
        return
      }

      // Calculate maturity date
      const maturityDate = new Date(transferForm.value.date)
      maturityDate.setMonth(maturityDate.getMonth() + transferForm.value.fdTenureMonths)

      // Deduct from source account
      await accountsStore.updateBalance(transferForm.value.fromAccountId, transferForm.value.amount, 'subtract')

      // Create transaction record for fixed deposit
      await transactionsStore.createTransaction({
        type: 'expense',
        amount: transferForm.value.amount,
        categoryId: 'other_expense',
        accountId: transferForm.value.fromAccountId,
        date: dateISO,
        description: `Fixed Deposit: ${transferForm.value.fdName}`,
        tags: ['fixed_deposit']
      })

      success('Fixed deposit transaction created successfully')
    }

    closeTransferModal()
    // Refresh accounts to show updated balances
    await accountsStore.fetchAccounts()
    await loadTransactions()
  } catch (err) {
    error(err.response?.data?.message || err.message || 'Error transferring funds')
  }
}

const closeTransferModal = () => {
  showTransferModal.value = false
  transferForm.value = {
    fromAccountId: '',
    destinationType: 'account',
    toAccountId: '',
    savingsGoalId: '',
    fdName: '',
    fdInterestRate: 0,
    fdTenureMonths: 12,
    amount: 0,
    date: new Date().toISOString().split('T')[0],
    description: 'Transfer between accounts'
  }
}

// File attachment helper functions
const viewAttachments = (transaction) => {
  viewingTransaction.value = transaction
  showAttachmentsModal.value = true
}

const isImageUrl = (url) => {
  return /\.(jpg|jpeg|png|gif|webp)$/i.test(url)
}

const getFileNameFromUrl = (url) => {
  return url.split('/').pop()
}

const openAttachment = (url) => {
  window.open(url, '_blank')
}

// Watch for manual date changes to reset quick filter
watch([() => filters.value.startDate, () => filters.value.endDate], () => {
  // Reset quick filter to custom when dates are manually changed
  // Only if a quick filter was previously selected
  if (quickDateFilter.value !== '') {
    quickDateFilter.value = ''
  }
})

onMounted(async () => {
  await Promise.all([
    transactionsStore.fetchTransactions(currentPage.value, itemsPerPage.value),
    accountsStore.fetchAccounts(),
    creditCardsStore.fetchCreditCards()
  ])
})
</script>

<style scoped>
.attachment-card {
  border: 1px solid #e3e8ee;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.2s ease;
}

.attachment-card:hover {
  border-color: #635bff;
  box-shadow: 0 2px 8px rgba(99, 91, 255, 0.2);
}

.attachment-preview {
  width: 100%;
  height: 200px;
  overflow: hidden;
  background-color: #f6f9fc;
  display: flex;
  align-items: center;
  justify-content: center;
}

.attachment-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.attachment-file-icon {
  width: 100%;
  height: 200px;
  background-color: #f6f9fc;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #8898aa;
}

.attachment-info {
  padding: 12px;
  background-color: #ffffff;
}

/* Pagination styles */
.pagination-info {
  font-size: 0.875rem;
  color: #4b5563;
  font-weight: 500;
}

.pagination .page-link {
  color: #4b5563;
  border-color: #d1d5db;
  font-weight: 500;
}

.pagination .page-link:hover {
  color: #635bff;
  background-color: #f3f4f6;
  border-color: #635bff;
}

.pagination .page-item.active .page-link {
  background-color: #635bff;
  border-color: #635bff;
  color: #ffffff;
  font-weight: 600;
}

.pagination .page-item.disabled .page-link {
  color: #9ca3af;
  background-color: #f9fafb;
  border-color: #e5e7eb;
}

/* Dark mode styles */
.dark-mode .attachment-card {
  border-color: #374151;
  background-color: #1f2937;
}

.dark-mode .attachment-preview,
.dark-mode .attachment-file-icon {
  background-color: #111827;
}

.dark-mode .attachment-info {
  background-color: #1f2937;
}

.dark-mode .attachment-info p {
  color: #e5e7eb;
}

.dark-mode .pagination-info {
  color: #d1d5db;
}

.dark-mode .pagination .page-link {
  color: #d1d5db;
  background-color: #374151;
  border-color: #4b5563;
}

.dark-mode .pagination .page-link:hover {
  color: #818cf8;
  background-color: #4b5563;
  border-color: #818cf8;
}

.dark-mode .pagination .page-item.active .page-link {
  background-color: #635bff;
  border-color: #635bff;
  color: #ffffff;
}

.dark-mode .pagination .page-item.disabled .page-link {
  color: #6b7280;
  background-color: #1f2937;
  border-color: #374151;
}

/* Transaction type badges */
.transaction-type-badge {
  font-size: 0.75rem;
  padding: 0.35rem 0.65rem;
  font-weight: 600;
  text-transform: capitalize;
}

.badge-income {
  background-color: #dbeafe;
  color: #1e40af;
}

.badge-expense {
  background-color: #fee2e2;
  color: #991b1b;
}

.badge-transfer {
  background-color: #e0e7ff;
  color: #3730a3;
}

/* Transaction amounts */
.transaction-amount {
  font-size: 0.95rem;
}

.amount-income {
  color: #1e40af;
}

.amount-expense {
  color: #991b1b;
}

.amount-transfer {
  color: #3730a3;
}

/* Dark mode transaction colors */
.dark-mode .badge-income {
  background-color: #1e3a5f;
  color: #93c5fd;
}

.dark-mode .badge-expense {
  background-color: #5f1e1e;
  color: #fca5a5;
}

.dark-mode .badge-transfer {
  background-color: #312e5f;
  color: #a5b4fc;
}

.dark-mode .amount-income {
  color: #93c5fd;
}

.dark-mode .amount-expense {
  color: #fca5a5;
}

.dark-mode .amount-transfer {
  color: #a5b4fc;
}

/* Grouped transactions styles */
.grouped-transactions {
  padding: 0;
}

.day-group {
  border-bottom: 1px solid #e5e7eb;
}

.day-group:last-child {
  border-bottom: none;
}

.day-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  background-color: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
}

.day-date {
  font-size: 1rem;
}

.day-summary {
  font-size: 0.875rem;
}

.day-group .table {
  margin-bottom: 0;
}

.day-group .table tbody tr:last-child td {
  border-bottom: none;
}

/* Dark mode grouped transactions */
.dark-mode .day-group {
  border-bottom-color: #374151;
}

.dark-mode .day-header {
  background-color: #1f2937;
  border-bottom-color: #374151;
}

.dark-mode .day-date strong,
.dark-mode .day-summary {
  color: #e5e7eb;
}
</style>
