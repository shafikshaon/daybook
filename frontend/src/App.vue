<template>
  <div id="app" :class="{ 'dark-mode': isDarkMode }">
    <!-- Top Navbar -->
    <nav class="top-navbar" v-if="authStore.isAuthenticated && !isAuthPage">
      <div class="top-navbar-content">
        <!-- Mobile Menu Toggle (visible only on mobile) -->
        <button class="mobile-menu-toggle" @click="toggleMobileMenu" v-if="!isAuthPage">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 16 16">
            <path fill-rule="evenodd" d="M2.5 12a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5m0-4a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5m0-4a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 0 1H3a.5.5 0 0 1-.5-.5"/>
          </svg>
        </button>

        <!-- Desktop Sidebar Toggle (hidden on mobile) -->
        <button class="sidebar-toggle-btn" @click="toggleSidebar" v-if="!isAuthPage">
          ☰
        </button>

        <div class="navbar-brand">
          💰 Expense Tracker
        </div>

        <div class="navbar-right" v-if="authStore.isAuthenticated">
          <!-- Profile Dropdown -->
          <div class="dropdown">
            <button
              class="btn btn-sm profile-btn dropdown-toggle"
              type="button"
              @click="showProfileDropdown = !showProfileDropdown"
            >
              <span class="profile-avatar">{{ userInitials }}</span>
              <span class="profile-name">{{ userName }}</span>
            </button>
            <ul
              class="dropdown-menu dropdown-menu-end"
              :class="{ show: showProfileDropdown }"
            >
              <li><router-link class="dropdown-item" to="/settings">⚙️ Settings</router-link></li>
              <li><button class="dropdown-item" @click="toggleDarkMode">🌙 {{ isDarkMode ? 'Light' : 'Dark' }} Mode</button></li>
              <li><hr class="dropdown-divider"></li>
              <li><button class="dropdown-item text-danger" @click="handleSignOut">🚪 Sign Out</button></li>
            </ul>
          </div>
        </div>
      </div>
    </nav>

    <!-- Mobile Drawer Menu -->
    <div class="mobile-drawer-overlay" :class="{ show: showMobileMenu }" @click="closeMobileMenu" v-if="!isAuthPage"></div>
    <nav class="mobile-drawer" :class="{ show: showMobileMenu }" v-if="!isAuthPage">
      <div class="mobile-drawer-header">
        <div class="drawer-user-info">
          <span class="user-avatar">{{ userInitials }}</span>
          <div class="user-details">
            <span class="user-name">{{ userName }}</span>
            <span class="user-email">{{ userEmail }}</span>
          </div>
        </div>
        <button class="drawer-close-btn" @click="closeMobileMenu">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 16 16">
            <path d="M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8z"/>
          </svg>
        </button>
      </div>

      <div class="mobile-drawer-content">
        <ul class="drawer-nav">
          <li class="drawer-nav-item">
            <router-link to="/" class="drawer-nav-link" @click="closeMobileMenu">
              <svg class="drawer-nav-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
                <path d="M8.707 1.5a1 1 0 0 0-1.414 0L.646 8.146a.5.5 0 0 0 .708.708L2 8.207V13.5A1.5 1.5 0 0 0 3.5 15h9a1.5 1.5 0 0 0 1.5-1.5V8.207l.646.647a.5.5 0 0 0 .708-.708L13 5.793V2.5a.5.5 0 0 0-.5-.5h-1a.5.5 0 0 0-.5.5v1.293zM13 7.207V13.5a.5.5 0 0 1-.5.5h-9a.5.5 0 0 1-.5-.5V7.207l5-5z"/>
              </svg>
              <span>Dashboard</span>
            </router-link>
          </li>

          <li class="drawer-section-title">Accounts & Transactions</li>
          <li class="drawer-nav-item">
            <router-link to="/accounts" class="drawer-nav-link" @click="closeMobileMenu">
              <svg class="drawer-nav-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
                <path d="M4 10.781c.148 1.667 1.513 2.85 3.591 3.003V15h1.043v-1.216c2.27-.179 3.678-1.438 3.678-3.3 0-1.59-.947-2.51-2.956-3.028l-.722-.187V3.467c1.122.11 1.879.714 2.07 1.616h1.47c-.166-1.6-1.54-2.748-3.54-2.875V1H7.591v1.233c-1.939.23-3.27 1.472-3.27 3.156 0 1.454.966 2.483 2.661 2.917l.61.162v4.031c-1.149-.17-1.94-.8-2.131-1.718zm3.391-3.836c-1.043-.263-1.6-.825-1.6-1.616 0-.944.704-1.641 1.8-1.828v3.495l-.2-.05zm1.591 1.872c1.287.323 1.852.859 1.852 1.769 0 1.097-.826 1.828-2.2 1.939V8.73z"/>
              </svg>
              <span>Accounts</span>
            </router-link>
          </li>
          <li class="drawer-nav-item">
            <router-link to="/transactions" class="drawer-nav-link" @click="closeMobileMenu">
              <svg class="drawer-nav-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
                <path d="M0 4a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v5H0zm11.5 1a.5.5 0 0 0-.5.5v1a.5.5 0 0 0 .5.5h2a.5.5 0 0 0 .5-.5v-1a.5.5 0 0 0-.5-.5zM0 11v-1h16v1a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2"/>
              </svg>
              <span>Transactions</span>
            </router-link>
          </li>
          <li class="drawer-nav-item">
            <router-link to="/credit-cards" class="drawer-nav-link" @click="closeMobileMenu">
              <svg class="drawer-nav-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
                <path d="M11 5.5a.5.5 0 0 1 .5-.5h2a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5z"/>
                <path d="M2 2a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V4a2 2 0 0 0-2-2zm13 2v5H1V4a1 1 0 0 1 1-1h12a1 1 0 0 1 1 1m-1 9H2a1 1 0 0 1-1-1v-1h14v1a1 1 0 0 1-1 1"/>
              </svg>
              <span>Credit Cards</span>
            </router-link>
          </li>

          <li class="drawer-section-title">Planning & Budgets</li>
          <li class="drawer-nav-item">
            <router-link to="/budgets" class="drawer-nav-link" @click="closeMobileMenu">
              <svg class="drawer-nav-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
                <path d="M1 3a1 1 0 0 1 1-1h12a1 1 0 0 1 1 1zm7 8a2 2 0 1 0 0-4 2 2 0 0 0 0 4"/>
                <path d="M0 5a1 1 0 0 1 1-1h14a1 1 0 0 1 1 1v8a1 1 0 0 1-1 1H1a1 1 0 0 1-1-1zm3 0a2 2 0 0 1-2 2v4a2 2 0 0 1 2 2h10a2 2 0 0 1 2-2V7a2 2 0 0 1-2-2z"/>
              </svg>
              <span>Budgets</span>
            </router-link>
          </li>
          <li class="drawer-nav-item">
            <router-link to="/bills" class="drawer-nav-link" @click="closeMobileMenu">
              <svg class="drawer-nav-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
                <path d="M14 4.5V14a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V2a2 2 0 0 1 2-2h5.5zm-3 0A1.5 1.5 0 0 1 9.5 3V1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V4.5z"/>
              </svg>
              <span>Bills & Reminders</span>
            </router-link>
          </li>
          <li class="drawer-nav-item">
            <router-link to="/savings-goals" class="drawer-nav-link" @click="closeMobileMenu">
              <svg class="drawer-nav-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
                <path d="M9.669.864 8 0 6.331.864l-1.858.282-.842 1.68-1.337 1.32L2.6 6l-.306 1.854 1.337 1.32.842 1.68 1.858.282L8 12l1.669-.864 1.858-.282.842-1.68 1.337-1.32L13.4 6l.306-1.854-1.337-1.32-.842-1.68zm1.196 1.193.684 1.365 1.086 1.072L12.387 6l.248 1.506-1.086 1.072-.684 1.365-1.51.229L8 10.874l-1.355-.702-1.51-.229-.684-1.365-1.086-1.072L3.614 6l-.25-1.506 1.087-1.072.684-1.365 1.51-.229L8 1.126l1.356.702z"/>
                <path d="M4 11.794V16l4-1 4 1v-4.206l-2.018.306L8 13.126 6.018 12.1z"/>
              </svg>
              <span>Savings Goals</span>
            </router-link>
          </li>

          <li class="drawer-section-title">Investments & Assets</li>
          <li class="drawer-nav-item">
            <router-link to="/investments" class="drawer-nav-link" @click="closeMobileMenu">
              <svg class="drawer-nav-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
                <path fill-rule="evenodd" d="M0 0h1v15h15v1H0zm14.817 3.113a.5.5 0 0 1 .07.704l-4.5 5.5a.5.5 0 0 1-.74.037L7.06 6.767l-3.656 5.027a.5.5 0 0 1-.808-.588l4-5.5a.5.5 0 0 1 .758-.06l2.609 2.61 4.15-5.073a.5.5 0 0 1 .704-.07"/>
              </svg>
              <span>Investments</span>
            </router-link>
          </li>

          <li class="drawer-section-title">Reports & Analytics</li>
          <li class="drawer-nav-item">
            <router-link to="/reports" class="drawer-nav-link" @click="closeMobileMenu">
              <svg class="drawer-nav-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
                <path d="M4 11H2v3h2zm5-4H7v7h2zm5-5v12h-2V2zm-2-1a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h2a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1zM6 7a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1v7a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1zm-5 4a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1v3a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1z"/>
              </svg>
              <span>Reports</span>
            </router-link>
          </li>
          <li class="drawer-section-title">Settings</li>
          <li class="drawer-nav-item">
            <router-link to="/account-types" class="drawer-nav-link" @click="closeMobileMenu">
              <svg class="drawer-nav-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
                <path d="M1 2.5A1.5 1.5 0 0 1 2.5 1h3A1.5 1.5 0 0 1 7 2.5v3A1.5 1.5 0 0 1 5.5 7h-3A1.5 1.5 0 0 1 1 5.5zM2.5 2a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5zm6.5.5A1.5 1.5 0 0 1 10.5 1h3A1.5 1.5 0 0 1 15 2.5v3A1.5 1.5 0 0 1 13.5 7h-3A1.5 1.5 0 0 1 9 5.5zm1.5-.5a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5zM1 10.5A1.5 1.5 0 0 1 2.5 9h3A1.5 1.5 0 0 1 7 10.5v3A1.5 1.5 0 0 1 5.5 15h-3A1.5 1.5 0 0 1 1 13.5zm1.5-.5a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5zm6.5.5A1.5 1.5 0 0 1 10.5 9h3a1.5 1.5 0 0 1 1.5 1.5v3a1.5 1.5 0 0 1-1.5 1.5h-3A1.5 1.5 0 0 1 9 13.5zm1.5-.5a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5z"/>
              </svg>
              <span>Account Types</span>
            </router-link>
          </li>
          <li class="drawer-nav-item">
            <router-link to="/settings" class="drawer-nav-link" @click="closeMobileMenu">
              <svg class="drawer-nav-icon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
                <path d="M9.405 1.05c-.413-1.4-2.397-1.4-2.81 0l-.1.34a1.464 1.464 0 0 1-2.105.872l-.31-.17c-1.283-.698-2.686.705-1.987 1.987l.169.311c.446.82.023 1.841-.872 2.105l-.34.1c-1.4.413-1.4 2.397 0 2.81l.34.1a1.464 1.464 0 0 1 .872 2.105l-.17.31c-.698 1.283.705 2.686 1.987 1.987l.311-.169a1.464 1.464 0 0 1 2.105.872l.1.34c.413 1.4 2.397 1.4 2.81 0l.1-.34a1.464 1.464 0 0 1 2.105-.872l.31.17c1.283.698 2.686-.705 1.987-1.987l-.169-.311a1.464 1.464 0 0 1 .872-2.105l.34-.1c1.4-.413 1.4-2.397 0-2.81l-.34-.1a1.464 1.464 0 0 1-.872-2.105l.17-.31c.698-1.283-.705-2.686-1.987-1.987l-.311.169a1.464 1.464 0 0 1-2.105-.872zM8 10.93a2.929 2.929 0 1 1 0-5.86 2.929 2.929 0 0 1 0 5.858z"/>
              </svg>
              <span>Settings</span>
            </router-link>
          </li>
        </ul>
      </div>
    </nav>

    <!-- Sidebar Navigation (hidden for auth pages) -->
    <nav class="sidebar" :class="{ 'collapsed': !showSidebar }" v-if="!isAuthPage">
      <div class="sidebar-header">
        <div class="sidebar-user-info" v-if="showSidebar">
          <span class="user-initials">{{ userInitials }}</span>
          <span class="user-name">{{ userName }}</span>
        </div>
        <span v-else class="sidebar-user-collapsed">{{ userInitials }}</span>
      </div>

      <ul class="nav flex-column">
        <!-- Dashboard -->
        <li class="nav-item">
          <router-link to="/" class="nav-link" exact-active-class="active" :title="!showSidebar ? 'Dashboard' : ''">
            <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path d="M8.707 1.5a1 1 0 0 0-1.414 0L.646 8.146a.5.5 0 0 0 .708.708L2 8.207V13.5A1.5 1.5 0 0 0 3.5 15h9a1.5 1.5 0 0 0 1.5-1.5V8.207l.646.647a.5.5 0 0 0 .708-.708L13 5.793V2.5a.5.5 0 0 0-.5-.5h-1a.5.5 0 0 0-.5.5v1.293zM13 7.207V13.5a.5.5 0 0 1-.5.5h-9a.5.5 0 0 1-.5-.5V7.207l5-5z"/>
            </svg>
            <span v-if="showSidebar">Dashboard</span>
          </router-link>
        </li>

        <!-- Accounts & Transactions -->
        <li v-if="showSidebar" class="nav-section-title">Accounts & Transactions</li>
        <li v-if="!showSidebar" class="nav-divider"></li>

        <li class="nav-item">
          <router-link to="/accounts" class="nav-link" active-class="active" :title="!showSidebar ? 'Accounts' : ''">
            <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path d="M4 10.781c.148 1.667 1.513 2.85 3.591 3.003V15h1.043v-1.216c2.27-.179 3.678-1.438 3.678-3.3 0-1.59-.947-2.51-2.956-3.028l-.722-.187V3.467c1.122.11 1.879.714 2.07 1.616h1.47c-.166-1.6-1.54-2.748-3.54-2.875V1H7.591v1.233c-1.939.23-3.27 1.472-3.27 3.156 0 1.454.966 2.483 2.661 2.917l.61.162v4.031c-1.149-.17-1.94-.8-2.131-1.718zm3.391-3.836c-1.043-.263-1.6-.825-1.6-1.616 0-.944.704-1.641 1.8-1.828v3.495l-.2-.05zm1.591 1.872c1.287.323 1.852.859 1.852 1.769 0 1.097-.826 1.828-2.2 1.939V8.73z"/>
            </svg>
            <span v-if="showSidebar">Accounts</span>
          </router-link>
        </li>
        <li class="nav-item">
          <router-link to="/transactions" class="nav-link" active-class="active" :title="!showSidebar ? 'Transactions' : ''">
            <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path d="M0 4a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v5H0zm11.5 1a.5.5 0 0 0-.5.5v1a.5.5 0 0 0 .5.5h2a.5.5 0 0 0 .5-.5v-1a.5.5 0 0 0-.5-.5zM0 11v-1h16v1a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2"/>
            </svg>
            <span v-if="showSidebar">Transactions</span>
          </router-link>
        </li>
        <li class="nav-item">
          <router-link to="/credit-cards" class="nav-link" active-class="active" :title="!showSidebar ? 'Credit Cards' : ''">
            <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path d="M11 5.5a.5.5 0 0 1 .5-.5h2a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5z"/>
              <path d="M2 2a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V4a2 2 0 0 0-2-2zm13 2v5H1V4a1 1 0 0 1 1-1h12a1 1 0 0 1 1 1m-1 9H2a1 1 0 0 1-1-1v-1h14v1a1 1 0 0 1-1 1"/>
            </svg>
            <span v-if="showSidebar">Credit Cards</span>
          </router-link>
        </li>

        <!-- Planning & Budgets -->
        <li v-if="showSidebar" class="nav-section-title">Planning & Budgets</li>
        <li v-if="!showSidebar" class="nav-divider"></li>

        <li class="nav-item">
          <router-link to="/budgets" class="nav-link" active-class="active" :title="!showSidebar ? 'Budgets' : ''">
            <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path d="M1 3a1 1 0 0 1 1-1h12a1 1 0 0 1 1 1zm7 8a2 2 0 1 0 0-4 2 2 0 0 0 0 4"/>
              <path d="M0 5a1 1 0 0 1 1-1h14a1 1 0 0 1 1 1v8a1 1 0 0 1-1 1H1a1 1 0 0 1-1-1zm3 0a2 2 0 0 1-2 2v4a2 2 0 0 1 2 2h10a2 2 0 0 1 2-2V7a2 2 0 0 1-2-2z"/>
            </svg>
            <span v-if="showSidebar">Budgets</span>
          </router-link>
        </li>
        <li class="nav-item">
          <router-link to="/bills" class="nav-link" active-class="active" :title="!showSidebar ? 'Bills & Reminders' : ''">
            <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path d="M14 4.5V14a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V2a2 2 0 0 1 2-2h5.5zm-3 0A1.5 1.5 0 0 1 9.5 3V1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V4.5z"/>
            </svg>
            <span v-if="showSidebar">Bills & Reminders</span>
          </router-link>
        </li>
        <li class="nav-item">
          <router-link to="/savings-goals" class="nav-link" active-class="active" :title="!showSidebar ? 'Savings Goals' : ''">
            <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path d="M9.669.864 8 0 6.331.864l-1.858.282-.842 1.68-1.337 1.32L2.6 6l-.306 1.854 1.337 1.32.842 1.68 1.858.282L8 12l1.669-.864 1.858-.282.842-1.68 1.337-1.32L13.4 6l.306-1.854-1.337-1.32-.842-1.68zm1.196 1.193.684 1.365 1.086 1.072L12.387 6l.248 1.506-1.086 1.072-.684 1.365-1.51.229L8 10.874l-1.355-.702-1.51-.229-.684-1.365-1.086-1.072L3.614 6l-.25-1.506 1.087-1.072.684-1.365 1.51-.229L8 1.126l1.356.702z"/>
              <path d="M4 11.794V16l4-1 4 1v-4.206l-2.018.306L8 13.126 6.018 12.1z"/>
            </svg>
            <span v-if="showSidebar">Savings Goals</span>
          </router-link>
        </li>

        <!-- Investments & Assets -->
        <li v-if="showSidebar" class="nav-section-title">Investments & Assets</li>
        <li v-if="!showSidebar" class="nav-divider"></li>

        <li class="nav-item">
          <router-link to="/investments" class="nav-link" active-class="active" :title="!showSidebar ? 'Investments' : ''">
            <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path fill-rule="evenodd" d="M0 0h1v15h15v1H0zm14.817 3.113a.5.5 0 0 1 .07.704l-4.5 5.5a.5.5 0 0 1-.74.037L7.06 6.767l-3.656 5.027a.5.5 0 0 1-.808-.588l4-5.5a.5.5 0 0 1 .758-.06l2.609 2.61 4.15-5.073a.5.5 0 0 1 .704-.07"/>
            </svg>
            <span v-if="showSidebar">Investments</span>
          </router-link>
        </li>

        <!-- Reports & Analytics -->
        <li v-if="showSidebar" class="nav-section-title">Reports & Analytics</li>
        <li v-if="!showSidebar" class="nav-divider"></li>

        <li class="nav-item">
          <router-link to="/reports" class="nav-link" active-class="active" :title="!showSidebar ? 'Reports' : ''">
            <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path d="M4 11H2v3h2zm5-4H7v7h2zm5-5v12h-2V2zm-2-1a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h2a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1zM6 7a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1v7a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1zm-5 4a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1v3a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1z"/>
            </svg>
            <span v-if="showSidebar">Reports</span>
          </router-link>
        </li>

        <!-- Settings -->
        <li v-if="showSidebar" class="nav-section-title">Settings</li>
        <li v-if="!showSidebar" class="nav-divider"></li>

        <li class="nav-item">
          <router-link to="/account-types" class="nav-link" active-class="active" :title="!showSidebar ? 'Account Types' : ''">
            <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path d="M1 2.5A1.5 1.5 0 0 1 2.5 1h3A1.5 1.5 0 0 1 7 2.5v3A1.5 1.5 0 0 1 5.5 7h-3A1.5 1.5 0 0 1 1 5.5zM2.5 2a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5zm6.5.5A1.5 1.5 0 0 1 10.5 1h3A1.5 1.5 0 0 1 15 2.5v3A1.5 1.5 0 0 1 13.5 7h-3A1.5 1.5 0 0 1 9 5.5zm1.5-.5a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5zM1 10.5A1.5 1.5 0 0 1 2.5 9h3A1.5 1.5 0 0 1 7 10.5v3A1.5 1.5 0 0 1 5.5 15h-3A1.5 1.5 0 0 1 1 13.5zm1.5-.5a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5zm6.5.5A1.5 1.5 0 0 1 10.5 9h3a1.5 1.5 0 0 1 1.5 1.5v3a1.5 1.5 0 0 1-1.5 1.5h-3A1.5 1.5 0 0 1 9 13.5zm1.5-.5a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5z"/>
            </svg>
            <span v-if="showSidebar">Account Types</span>
          </router-link>
        </li>
        <li class="nav-item">
          <router-link to="/settings" class="nav-link" active-class="active" :title="!showSidebar ? 'Settings' : ''">
            <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
              <path d="M9.405 1.05c-.413-1.4-2.397-1.4-2.81 0l-.1.34a1.464 1.464 0 0 1-2.105.872l-.31-.17c-1.283-.698-2.686.705-1.987 1.987l.169.311c.446.82.023 1.841-.872 2.105l-.34.1c-1.4.413-1.4 2.397 0 2.81l.34.1a1.464 1.464 0 0 1 .872 2.105l-.17.31c-.698 1.283.705 2.686 1.987 1.987l.311-.169a1.464 1.464 0 0 1 2.105.872l.1.34c.413 1.4 2.397 1.4 2.81 0l.1-.34a1.464 1.464 0 0 1 2.105-.872l.31.17c1.283.698 2.686-.705 1.987-1.987l-.169-.311a1.464 1.464 0 0 1 .872-2.105l.34-.1c1.4-.413 1.4-2.397 0-2.81l-.34-.1a1.464 1.464 0 0 1-.872-2.105l.17-.31c.698-1.283-.705-2.686-1.987-1.987l-.311.169a1.464 1.464 0 0 1-2.105-.872zM8 10.93a2.929 2.929 0 1 1 0-5.86 2.929 2.929 0 0 1 0 5.858z"/>
            </svg>
            <span v-if="showSidebar">Settings</span>
          </router-link>
        </li>
      </ul>
    </nav>

    <!-- Main Content -->
    <main class="main-content" :class="{ 'sidebar-collapsed': !showSidebar || isAuthPage, 'auth-page': isAuthPage }">
      <router-view />
    </main>

    <!-- Mobile Bottom Navigation (iOS/Android Style) -->
    <nav class="mobile-bottom-nav" v-if="!isAuthPage">
      <router-link to="/" class="nav-item" exact-active-class="active">
        <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 16 16">
          <path d="M8.707 1.5a1 1 0 0 0-1.414 0L.646 8.146a.5.5 0 0 0 .708.708L2 8.207V13.5A1.5 1.5 0 0 0 3.5 15h9a1.5 1.5 0 0 0 1.5-1.5V8.207l.646.647a.5.5 0 0 0 .708-.708L13 5.793V2.5a.5.5 0 0 0-.5-.5h-1a.5.5 0 0 0-.5.5v1.293zM13 7.207V13.5a.5.5 0 0 1-.5.5h-9a.5.5 0 0 1-.5-.5V7.207l5-5z"/>
        </svg>
        <span class="nav-label">Home</span>
      </router-link>

      <router-link to="/transactions" class="nav-item" active-class="active">
        <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 16 16">
          <path d="M0 4a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v5H0zm11.5 1a.5.5 0 0 0-.5.5v1a.5.5 0 0 0 .5.5h2a.5.5 0 0 0 .5-.5v-1a.5.5 0 0 0-.5-.5zM0 11v-1h16v1a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2"/>
        </svg>
        <span class="nav-label">Payments</span>
      </router-link>

      <router-link to="/budgets" class="nav-item" active-class="active">
        <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 16 16">
          <path d="M1 3a1 1 0 0 1 1-1h12a1 1 0 0 1 1 1zm7 8a2 2 0 1 0 0-4 2 2 0 0 0 0 4"/>
          <path d="M0 5a1 1 0 0 1 1-1h14a1 1 0 0 1 1 1v8a1 1 0 0 1-1 1H1a1 1 0 0 1-1-1zm3 0a2 2 0 0 1-2 2v4a2 2 0 0 1 2 2h10a2 2 0 0 1 2-2V7a2 2 0 0 1-2-2z"/>
        </svg>
        <span class="nav-label">Budgets</span>
      </router-link>

      <router-link to="/reports" class="nav-item" active-class="active">
        <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 16 16">
          <path d="M4 11H2v3h2zm5-4H7v7h2zm5-5v12h-2V2zm-2-1a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h2a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1zM6 7a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1v7a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1zm-5 4a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1v3a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1z"/>
        </svg>
        <span class="nav-label">Reports</span>
      </router-link>

      <router-link to="/settings" class="nav-item" active-class="active">
        <svg class="nav-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 16 16">
          <path d="M9.405 1.05c-.413-1.4-2.397-1.4-2.81 0l-.1.34a1.464 1.464 0 0 1-2.105.872l-.31-.17c-1.283-.698-2.686.705-1.987 1.987l.169.311c.446.82.023 1.841-.872 2.105l-.34.1c-1.4.413-1.4 2.397 0 2.81l.34.1a1.464 1.464 0 0 1 .872 2.105l-.17.31c-.698 1.283.705 2.686 1.987 1.987l.311-.169a1.464 1.464 0 0 1 2.105.872l.1.34c.413 1.4 2.397 1.4 2.81 0l.1-.34a1.464 1.464 0 0 1 2.105-.872l.31.17c1.283.698 2.686-.705 1.987-1.987l-.169-.311a1.464 1.464 0 0 1 .872-2.105l.34-.1c1.4-.413 1.4-2.397 0-2.81l-.34-.1a1.464 1.464 0 0 1-.872-2.105l.17-.31c.698-1.283-.705-2.686-1.987-1.987l-.311.169a1.464 1.464 0 0 1-2.105-.872zM8 10.93a2.929 2.929 0 1 1 0-5.86 2.929 2.929 0 0 1 0 5.858z"/>
        </svg>
        <span class="nav-label">Settings</span>
      </router-link>
    </nav>

    <!-- Global Alert Toast -->
    <AlertToast ref="alertToastRef" />

    <!-- Global Confirm Modal -->
    <ConfirmModal
      :show="confirmState.show"
      :title="confirmState.title"
      :message="confirmState.message"
      :confirmText="confirmState.confirmText"
      :cancelText="confirmState.cancelText"
      :variant="confirmState.variant"
      @confirm="confirmState.onConfirm"
      @cancel="confirmState.onCancel"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { useAuthStore } from '@/stores/auth'
import { useRouter, useRoute } from 'vue-router'
import { useNotification } from '@/composables/useNotification'
import { AlertToast, ConfirmModal } from '@/components'

const settingsStore = useSettingsStore()
const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

// Notification system
const { registerAlertInstance, confirmState } = useNotification()
const alertToastRef = ref(null)

const showSidebar = ref(true)
const showProfileDropdown = ref(false)
const showMobileMenu = ref(false)

const isDarkMode = computed(() => settingsStore.isDarkMode)

const isAuthPage = computed(() => {
  return route.meta.hideLayout || route.path === '/login' || route.path === '/signup'
})

const userName = computed(() => {
  return authStore.user?.fullName || authStore.user?.username || 'User'
})

const userEmail = computed(() => {
  return authStore.user?.email || ''
})

const userInitials = computed(() => {
  return userName.value
    .split(' ')
    .map(n => n[0])
    .join('')
    .toUpperCase()
})

const toggleSidebar = () => {
  showSidebar.value = !showSidebar.value
}

const toggleMobileMenu = () => {
  showMobileMenu.value = !showMobileMenu.value
  // Prevent body scroll when menu is open
  if (showMobileMenu.value) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
}

const closeMobileMenu = () => {
  showMobileMenu.value = false
  document.body.style.overflow = ''
}

const toggleDarkMode = () => {
  settingsStore.toggleDarkMode()
  showProfileDropdown.value = false
}

const handleSignOut = async () => {
  showProfileDropdown.value = false
  const { confirm } = useNotification()
  const confirmed = await confirm({
    title: 'Sign Out',
    message: 'Are you sure you want to sign out?',
    confirmText: 'Sign Out',
    cancelText: 'Cancel',
    variant: 'danger'
  })
  if (confirmed) {
    await authStore.logout()
    router.push('/login')
  }
}

// Close dropdown when clicking outside
const handleClickOutside = (event) => {
  if (!event.target.closest('.dropdown')) {
    showProfileDropdown.value = false
  }
}

onMounted(async () => {
  // Register alert instance for global notifications
  if (alertToastRef.value) {
    registerAlertInstance(alertToastRef.value)
  }

  // Initialize authentication
  await authStore.initializeAuth()

  // Load settings
  settingsStore.loadSettings()

  // Add click outside listener
  document.addEventListener('click', handleClickOutside)

  // Check if on mobile, collapse sidebar
  if (window.innerWidth <= 768) {
    showSidebar.value = false
  }

  // Responsive handler
  const handleResize = () => {
    if (window.innerWidth <= 768) {
      showSidebar.value = false
    }
  }
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style>
* {
  box-sizing: border-box;
}

#app {
  min-height: 100vh;
  overflow-x: hidden;
  width: 100%;
}

body {
  overflow-x: hidden;
  width: 100%;
  margin: 0;
  padding: 0;
}

.main-content.auth-page {
  margin-left: 0 !important;
  margin-top: 0 !important;
  padding: 0 !important;
  width: 100% !important;
  max-width: 100vw !important;
  overflow-x: hidden !important;
}
</style>
