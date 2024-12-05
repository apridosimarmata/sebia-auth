package user

type UserUsecase interface {
	GetUserBasicInformation() // for settings page
	GetUserBusinessEntity()
	GetUserAffiliateStatus()
}
