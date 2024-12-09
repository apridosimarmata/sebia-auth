package auth

import (
	"context"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/user"
	"mini-wallet/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	RegisterByGoogle(ctx context.Context, req GoogleRegisterDTO) (res response.Response[AuthenticationResponse])
	AuthenticateByGoogle(ctx context.Context, req GoogleRegisterDTO) (res response.Response[AuthenticationResponse])
	AuthenticateRegularUser(ctx context.Context, req AuthenticationDTO) (res response.Response[AuthenticationResponse])
	RegisterUser(ctx context.Context, req UserRegistrationDTO) (res response.Response[interface{}])
	SendPasswordResetLink(ctx context.Context, req PasswordResetDTO) (res response.Response[string])
	ResetUserPassword(ctx context.Context, req PasswordResetSubmissionDTO) (res response.Response[string])
	VerifyResetPasswordToken(ctx context.Context, req VerifyResetPasswordTokenDTO) (res response.Response[string])
	CheckIdentifier(ctx context.Context, req CheckIndentifierDTO) (res response.Response[string])
	VerifyPhoneNumber(ctx context.Context, req VerifyEmailDTO) (res response.Response[AuthenticationResponse])
	RefreshAccess(ctx context.Context) (res response.Response[AuthenticationResponse])
	RegisterUserFromInquiry(ctx context.Context, req AuthFromInquiryDTO) (res response.Response[string])
	AuthenticateFromInquiry(ctx context.Context, req AuthFromInquiryDTO) (res response.Response[AuthenticationResponse])
}

type AuthFromInquiryDTO struct {
	InquiryID string `json:"inquiry_id"`
	Password  string `json:"password"`
}

func (p *AuthFromInquiryDTO) Validate() error {
	err := utils.ValidateRequired(p.InquiryID)
	if err != nil {
		return err
	}

	err = utils.ValidateRequired(p.Password)
	if err != nil {
		return err
	}

	err = utils.ValidatePassword(p.Password)
	if err != nil {
		return err
	}

	return nil
}

type UserIDContext struct {
}

type TokenStatus struct {
}

type AuthRepository interface {
	AddToken(ctx context.Context, token string, walletId string) (err error)
	GetTokenWalletId(ctx context.Context, token string) (walletId string, err error)
}

type Token struct {
	Token string `json:"token"`
}

type UserRegistrationDTO struct {
	Name        string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (p *UserRegistrationDTO) Validate() (err error) {
	err = utils.ValidateRequired(p.Email)
	if err != nil {
		return err
	}

	err = utils.ValidateRequired(p.Name)
	if err != nil {
		return err
	}

	err = utils.ValidateRequired(p.Password)
	if err != nil {
		return err
	}

	err = utils.ValidateRequired(p.PhoneNumber)
	if err != nil {
		return err
	}

	err = utils.ValidateEmail(p.Email)
	if err != nil {
		return err
	}

	err = utils.ValidatePassword(p.Password)
	if err != nil {
		return err
	}

	newPhoneNumber, err := utils.ValidatePhoneNumber(p.PhoneNumber)
	if err != nil {
		return err
	}
	p.PhoneNumber = *newPhoneNumber

	err = utils.ValidateFullName(p.Name)
	if err != nil {
		return err
	}

	p.PhoneNumber = utils.ConvertPhoneNumber(p.PhoneNumber)

	return nil
}

func (p *AuthenticationDTO) Validate() (err error) {
	err = utils.ValidateRequired(p.Identifier)
	if err != nil {
		return err
	}

	err = utils.ValidateRequired(p.Password)
	if err != nil {
		return err
	}

	return nil
}

type CheckIndentifierDTO struct {
	Identifier string `json:"identifier"`
}

func (p *CheckIndentifierDTO) Validate() (err error) {
	err = utils.ValidateRequired(p.Identifier)
	if err != nil {
		return err
	}

	return nil
}

type PasswordResetDTO struct {
	Email string `json:"email"`
}

type VerifyResetPasswordTokenDTO struct {
	PasswordResetToken string `json:"password_reset_token"`
}

func (p *VerifyResetPasswordTokenDTO) Validate() (err error) {
	err = utils.ValidateRequired(p.PasswordResetToken)
	if err != nil {
		return err
	}

	return nil
}

type PasswordResetSubmissionDTO struct {
	Password           string `json:"password"`
	PasswordResetToken string `json:"password_reset_token"`
}

func (p *PasswordResetSubmissionDTO) Validate() (err error) {
	err = utils.ValidateRequired(p.PasswordResetToken)
	if err != nil {
		return err
	}

	err = utils.ValidateRequired(p.Password)
	if err != nil {
		return err
	}

	err = utils.ValidatePassword(p.Password)
	if err != nil {
		return err
	}

	return nil
}

type VerifyEmailDTO struct {
	Token string `json:"token"`
}

func (p *VerifyEmailDTO) Validate() (err error) {
	err = utils.ValidateRequired(p.Token)
	if err != nil {
		return err
	}

	return nil
}

func (p *PasswordResetDTO) Validate() (err error) {
	err = utils.ValidateRequired(p.Email)
	if err != nil {
		return err
	}

	err = utils.ValidateEmail(p.Email)
	if err != nil {
		return err
	}

	return nil
}

type AuthenticationDTO struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type GoogleRegisterDTO struct {
	Name            string          `json:"displayName"`
	StsTokenManager StsTokenManager `json:"stsTokenManager"`
	Email           string          `json:"email"`
}

type StsTokenManager struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthenticationResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AcessTokenClaims struct {
	jwt.RegisteredClaims
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}

func (p *GoogleRegisterDTO) ToUserEntity() (res *user.UserEntity, err error) {
	now, err := utils.GetJktTime()
	if err != nil {
		return nil, err
	}

	nowString := now.Format(time.RFC3339)
	return &user.UserEntity{
		UID:             utils.GenerateUniqueId(),
		Name:            p.Name,
		Email:           p.Email,
		EmailVerifiedAt: nowString,
		CreatedAt:       now.Format(time.RFC3339),
		UpdatedAt:       now.Format(time.RFC3339),
	}, nil
}

func (p *UserRegistrationDTO) ToTemporaryUserEntity() (res *user.TemporaryUserEntity, err error) {
	now, err := utils.GetJktTime()
	if err != nil {
		return nil, err
	}

	salt, err := utils.GenerateSalt(16)
	if err != nil {
		return nil, err
	}

	saltedPassword := p.Password + salt
	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	stringHashedPassword := string(hash)
	stringSalt := string(salt)

	VerificationToken, err := utils.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	return &user.TemporaryUserEntity{
		Name:              p.Name,
		Email:             p.Email,
		PhoneNumber:       &p.PhoneNumber,
		HashedPassword:    &stringHashedPassword,
		CreatedAt:         now.Format(time.RFC3339),
		ExpiredAt:         int(now.Add(time.Minute * 15).Unix()),
		PasswordSalt:      &stringSalt,
		VerificationToken: VerificationToken,
	}, nil
}
