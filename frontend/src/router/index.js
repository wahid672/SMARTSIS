import { createRouter, createWebHistory } from '@ionic/vue-router'
import { useAuthStore } from '@/stores/auth'

import LoginView from '@/views/LoginView.vue'
import TabsLayout from '@/views/TabsLayout.vue'
import DashboardView from '@/views/DashboardView.vue'
import ProfilView from '@/views/ProfilView.vue'
import PlaceholderView from '@/views/PlaceholderView.vue'
import PengaturanView from '@/views/PengaturanView.vue'
import RolesView from '@/views/RolesView.vue'
import PegawaiView from '@/views/PegawaiView.vue'
import SiswaView from '@/views/SiswaView.vue'
import KelasView from '@/views/KelasView.vue'
import KeuanganView from '@/views/KeuanganView.vue'
import TahunAjaranView from '@/views/TahunAjaranView.vue'
import PresensiView from '@/views/PresensiView.vue'
import NilaiView from '@/views/NilaiView.vue'
import MapelView from '@/views/MapelView.vue'
import TagihanView from '@/views/TagihanView.vue'

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
    path: '/app',
    component: TabsLayout,
    children: [
      {
        path: '',
        redirect: '/tabs/dashboard'
      },
      {
        path: '/tabs/dashboard',
        name: 'Dashboard',
        component: DashboardView,
        meta: { title: 'Beranda' }
      },
      {
        path: '/tabs/presensi',
        name: 'Presensi',
        component: PresensiView,
        meta: { title: 'Presensi & Absensi' }
      },
      {
        path: '/tabs/nilai',
        name: 'Nilai',
        component: NilaiView,
        meta: { title: 'Nilai & Rapor' }
      },
      {
        path: '/tabs/spp',
        redirect: '/keuangan/kas'
      },
      {
        path: '/keuangan/kas',
        name: 'BukuKas',
        component: KeuanganView,
        meta: { title: 'Buku Kas & Pos' }
      },
      {
        path: '/keuangan/tagihan',
        name: 'TagihanSiswa',
        component: TagihanView,
        meta: { title: 'Tagihan Siswa' }
      },
      {
        path: '/tabs/profil',
        name: 'Profil',
        component: ProfilView,
        meta: { title: 'Profil Akun' }
      },
      {
        path: '/admin/siswa',
        name: 'AdminSiswa',
        component: SiswaView,
        meta: { title: 'Manajemen Data Siswa' }
      },
      {
        path: '/admin/kelas',
        name: 'AdminKelas',
        component: KelasView,
        meta: { title: 'Master Data Kelas' }
      },
      {
        path: '/admin/mapel',
        name: 'AdminMapel',
        component: MapelView,
        meta: { title: 'Master Mata Pelajaran' }
      },
      {
        path: '/admin/tahun-ajaran',
        name: 'AdminTahunAjaran',
        component: TahunAjaranView,
        meta: { title: 'Master Tahun Ajaran' }
      },
      {
        path: '/admin/pegawai',
        name: 'AdminPegawai',
        component: PegawaiView,
        meta: { title: 'Manajemen Pegawai & Guru' }
      },
      {
        path: '/admin/roles',
        name: 'AdminRoles',
        component: RolesView,
        meta: { title: 'Pengaturan Hak Akses & Menu' }
      },
      {
        path: '/admin/pengaturan',
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
