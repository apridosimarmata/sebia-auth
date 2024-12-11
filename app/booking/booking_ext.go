package booking

import (
	"fmt"
	"mini-wallet/domain/business"
	"mini-wallet/domain/inquiry"
	"mini-wallet/domain/services"
	"strings"
)

func buildGuestBookingConfirmationMessage(inquiryEntity inquiry.InquiryEntity, confimationCode string, serviceEntity services.ServiceEntity, host business.BusinessEntity) string {
	return fmt.Sprintf("Hi %s, Berikut adalah kode konfirmasimu [%s] untuk pemesanan di %s!\nBila membutuhkan informasi, bisa menghubungi %s di +%s,\n\nSelamat liburan! 😊", inquiryEntity.FullName, confimationCode, serviceEntity.Title, host.Name, host.PhoneNumber)
}

func buildHostBookingConfirmationMessage(inquiryEntity inquiry.InquiryEntity, host business.BusinessEntity, serviceEntity services.ServiceEntity) string {
	return fmt.Sprintf("Hi %s\n\n%s telah melakukan pemesanan %s di tanggal berikut [%s].\nNomor kontakmu sudah dibagikan kepada tamu, selalu siap barangkali tamu menghubungi kamu ya! 😊", host.Name, inquiryEntity.FullName, serviceEntity.Title, strings.Join(inquiryEntity.SelectedDates, ", "))
}
