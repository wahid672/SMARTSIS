import { createRouter, createWebHistory } from '@ionic/vue-router'
import { useAuthStore } from '@/stores/auth'

import LoginView from '@/views/LoginView.vue'
import TabsLayout from '@/views/TabsLayout.vue'
import DashboardView from '@/views/DashboardView.vue'
import ProfilView from '@/views/ProfilView.vue'
import PlaceholderView from '@/views/PlaceholderView.vue'
import PengaturanView from '@/views/PengaturanView.vue'
import RolesView from '@/views/RolesView.vue'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: { public: true }
  },
  {
    path: '/tabs/',
    component: TabsLayout,
    children: [
      {
        path: '',
        redirect: '/tabs/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: DashboardView,
        meta: { title: 'Beranda' }
      },
      {
        path: 'presensi',
        name: 'Presensi',
        component: PlaceholderView,
        meta: { title: 'Presensi & Absensi' }
      },
      {
        path: 'nilai',
        name: 'Nilai',
        component: PlaceholderView,
        meta: { title: 'Nilai & Rapor' }
      },
      {
        path: 'spp',
        name: 'SPP',
        component: PlaceholderView,
        meta: { title: 'Keuangan & SPP' }
      },
      {
        path: 'profil',
        name: 'Profil',
        component: ProfilView,
        meta: { title: 'Profil Akun' }
      }
    ]
  },
  {
    path: '/admin/',
    component: TabsLayout,
    children: [
      {
        path: 'siswa',
        name: 'AdminSiswa',
        component: PlaceholderView,
        meta: { title: 'Manajemen Data Siswa' }
      },
      {
        path: 'pegawai',
        name: 'AdminPegawai',
        component: PlaceholderView,
        meta: { title: 'Manajemen Pegawai & Guru' }
      },
      {
        path: 'roles',
        name: 'AdminRoles',
        component: RolesView,
        meta: { title: 'Pengaturan Hak Akses & Menu' }
      },
      {
        path: 'pengaturan',
        name: 'AdminPengaturan',
        component: PengaturanView,
        meta: { title: 'Identitas Lembaga' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Route Guard untuk proteksi autentikasi
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  if (!to.meta.public && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    next('/tabs/dashboard')
  } else {
    next()
  }
})

export default router
