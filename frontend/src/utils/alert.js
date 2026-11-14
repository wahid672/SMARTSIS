import { alertController, toastController, IonicSafeString } from '@ionic/vue'

/**
 * Modern Native Ionic Confirm Dialog
 * Menggunakan alertController bawaan Ionic Framework:
 * 100% stabil, tidak bentrok dengan DOM/Z-Index, mendukung tampilan mobile & desktop.
 * 
 * @param {Object} options
 * @param {string} options.title - Judul modal
 * @param {string} options.text - Pesan/konten teks (mendukung tag HTML sederhana)
 * @param {string} [options.confirmButtonText='Ya, Lanjutkan']
 * @param {string} [options.cancelButtonText='Batal']
 * @param {boolean} [options.isDanger=false] - Jika bernilai true, tombol berwarna merah bahaya
 * @returns {Promise<boolean>}
 */
export const showConfirm = async ({
  title = 'Konfirmasi Tindakan',
  text = 'Apakah Anda yakin ingin melanjutkan?',
  confirmButtonText = 'Ya, Lanjutkan',
  cancelButtonText = 'Batal',
  isDanger = false
}) => {
  const alert = await alertController.create({
    header: title,
    message: new IonicSafeString(text.replace(/\n/g, '<br/>')),
    cssClass: 'smartsis-ionic-alert ' + (isDanger ? 'alert-is-danger' : 'alert-is-primary'),
    backdropDismiss: false,
    buttons: [
      {
        text: cancelButtonText,
        role: 'cancel',
        cssClass: 'alert-btn-cancel'
      },
      {
        text: confirmButtonText,
        role: 'confirm',
        cssClass: isDanger ? 'alert-btn-danger' : 'alert-btn-confirm'
      }
    ]
  })

  await alert.present()
  const { role } = await alert.onDidDismiss()
  return role === 'confirm'
}

/**
 * Modern Native Ionic Alert / Information Dialog
 */
export const showAlert = async ({
  title = 'Informasi',
  text = '',
  confirmButtonText = 'Tutup'
}) => {
  const alert = await alertController.create({
    header: title,
    message: new IonicSafeString(text.replace(/\n/g, '<br/>')),
    cssClass: 'smartsis-ionic-alert alert-is-info',
    backdropDismiss: true,
    buttons: [
      {
        text: confirmButtonText,
        role: 'confirm',
        cssClass: 'alert-btn-confirm'
      }
    ]
  })

  await alert.present()
  await alert.onDidDismiss()
}

/**
 * Modern Native Ionic Toast Notifier (Pojok Atas Layar)
 */
export const showToast = async (title = 'Berhasil', icon = 'success') => {
  const toast = await toastController.create({
    message: title,
    duration: 3200,
    position: 'top',
    color: icon === 'error' ? 'danger' : 'success',
    cssClass: 'smartsis-ionic-toast'
  })
  await toast.present()
}

export default {
  showConfirm,
  showAlert,
  showToast
}
