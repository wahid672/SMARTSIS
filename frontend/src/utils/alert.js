import Swal from 'sweetalert2'
import 'sweetalert2/dist/sweetalert2.min.css'

// Custom instance dengan styling modern ala SMARTSIS (Tailwind-like classes)
// `heightAuto: false` sangat krusial di SPA (Ionic / Vue / Tailwind) agar tidak memecah layout html/body
const customSwal = Swal.mixin({
  heightAuto: false,
  buttonsStyling: false,
  customClass: {
    popup: 'rounded-3xl p-6 shadow-2xl border border-slate-100 font-sans',
    title: 'text-lg font-bold text-slate-800 tracking-tight',
    htmlContainer: 'text-xs text-slate-600 mt-2 leading-relaxed',
    confirmButton: 'px-5 py-2.5 rounded-xl text-xs font-bold text-white shadow-md mx-1.5 transition-all outline-none',
    cancelButton: 'px-5 py-2.5 rounded-xl text-xs font-bold text-slate-600 bg-slate-100 hover:bg-slate-200 mx-1.5 transition-all outline-none',
    actions: 'mt-6 gap-2'
  }
})

/**
 * Modern Confirm Dialog
 * @param {Object} options
 * @param {string} options.title
 * @param {string} options.text
 * @param {string} [options.icon='warning'] - 'warning' | 'info' | 'question' | 'error' | 'success'
 * @param {string} [options.confirmButtonText='Ya, Lanjutkan']
 * @param {string} [options.cancelButtonText='Batal']
 * @param {boolean} [options.isDanger=false] - Jika aksi hapus / berbahaya (warna tombol merah)
 * @returns {Promise<boolean>}
 */
export const showConfirm = async ({
  title = 'Konfirmasi Tindakan',
  text = 'Apakah Anda yakin ingin melanjutkan?',
  icon = 'question',
  confirmButtonText = 'Ya, Lanjutkan',
  cancelButtonText = 'Batal',
  isDanger = false
}) => {
  const confirmClass = isDanger
    ? 'px-5 py-2.5 rounded-xl text-xs font-bold text-white bg-rose-600 hover:bg-rose-700 shadow-md shadow-rose-500/20 mx-1.5 transition-all outline-none'
    : 'px-5 py-2.5 rounded-xl text-xs font-bold text-white bg-emerald-600 hover:bg-emerald-700 shadow-md shadow-emerald-500/20 mx-1.5 transition-all outline-none'

  const result = await customSwal.fire({
    title,
    html: text.replace(/\n/g, '<br/>'),
    icon,
    showCancelButton: true,
    confirmButtonText,
    cancelButtonText,
    reverseButtons: true,
    focusCancel: isDanger,
    heightAuto: false,
    customClass: {
      popup: 'rounded-3xl p-6 shadow-2xl border border-slate-100',
      title: 'text-base sm:text-lg font-bold text-slate-800',
      htmlContainer: 'text-xs sm:text-sm text-slate-600 mt-2 leading-relaxed',
      confirmButton: confirmClass,
      cancelButton: 'px-5 py-2.5 rounded-xl text-xs font-bold text-slate-600 bg-slate-100 hover:bg-slate-200 mx-1.5 transition-all outline-none',
      actions: 'mt-5'
    }
  })

  return result.isConfirmed
}

/**
 * Modern Alert / Information Dialog
 */
export const showAlert = async ({
  title = 'Informasi',
  text = '',
  icon = 'info',
  confirmButtonText = 'Tutup'
}) => {
  return customSwal.fire({
    title,
    html: text.replace(/\n/g, '<br/>'),
    icon,
    confirmButtonText,
    heightAuto: false,
    customClass: {
      popup: 'rounded-3xl p-6 shadow-2xl border border-slate-100',
      title: 'text-base sm:text-lg font-bold text-slate-800',
      htmlContainer: 'text-xs sm:text-sm text-slate-600 mt-2 leading-relaxed',
      confirmButton: 'px-6 py-2.5 rounded-xl text-xs font-bold text-white bg-blue-600 hover:bg-blue-700 shadow-md shadow-blue-500/20 transition-all outline-none',
      actions: 'mt-5'
    }
  })
}

/**
 * Modern Toast Notifier (Pojok Kanan Atas)
 */
export const showToast = (title = 'Berhasil', icon = 'success') => {
  const Toast = Swal.mixin({
    toast: true,
    position: 'top-end',
    showConfirmButton: false,
    timer: 3000,
    timerProgressBar: true,
    heightAuto: false,
    didOpen: (toast) => {
      toast.onmouseenter = Swal.stopTimer
      toast.onmouseleave = Swal.resumeTimer
    }
  })
  Toast.fire({
    icon,
    title
  })
}

export default Swal
