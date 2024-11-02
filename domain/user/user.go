package user

import (
	"context"
	"mini-wallet/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserEntity struct {
	UID string `json:"uid"`

	Name           string  `bson:"name"`
	PhoneNumber    *string `bson:"phone_number"`
	HashedPassword *string `bson:"hashed_password"`
	Email          string  `bson:"email"`
	Gender         *string `bson:"gender"`

	CreatedAt             string  `bson:"created_at"`
	UpdatedAt             string  `bson:"updated_at"`
	EmailVerifiedAt       string  `bson:"email_verified_at"`
	PhoneNumberVerifiedAt *string `bson:"phone_number_verified_at"`
	PasswordSalt          *string `bson:"password_salt"`
}

func (p *UserEntity) ChangePassword(newPassword string) error {
	if p.PasswordSalt == nil {
		salt, err := utils.GenerateSalt(16)
		if err != nil {
			return err
		}
		p.PasswordSalt = &salt
	}

	saltedPassword := newPassword + *p.PasswordSalt
	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stringHashedPassword := string(hash)
	p.HashedPassword = &stringHashedPassword

	return nil
}

type UserPasswordResetEntity struct {
	UID string `json:"uid"`

	Email              string `bson:"email"`
	ExpiredAt          int64  `bson:"expired_at"`
	PasswordResetToken string `bson:"password_reset_token"`
}

type TemporaryUserEntity struct {
	UID string `json:"uid"`

	Name           string  `bson:"name"`
	PhoneNumber    *string `bson:"phone_number"`
	HashedPassword *string `bson:"hashed_password"`
	Email          string  `bson:"email"`

	CreatedAt         string  `bson:"created_at"`
	ExpiredAt         int     `bson:"expired_at"`
	PasswordSalt      *string `bson:"password_salt"`
	VerificationToken string  `bson:"verification_token"`
}

func (p *TemporaryUserEntity) ToUserEntity() (*UserEntity, error) {
	now, err := utils.GetJktTime()
	if err != nil {
		return nil, err
	}

	return &UserEntity{
		Name:           p.Name,
		Email:          p.Email,
		PhoneNumber:    p.PhoneNumber,
		HashedPassword: p.HashedPassword,
		Gender:         nil,

		CreatedAt:       p.CreatedAt,
		UpdatedAt:       now.Format(time.RFC3339),
		EmailVerifiedAt: now.Format(time.RFC3339),
		PasswordSalt:    p.PasswordSalt,
	}, nil
}

type UserDto struct {
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Email       string  `json:"email"`
	Gender      *string `json:"gender,omitempty"`
}

type UserRepository interface {
	InsertUser(ctx context.Context, user UserEntity) (err error)
	UpsertUser(ctx context.Context, user UserEntity) (err error)
	InsertTemporaryUser(ctx context.Context, user TemporaryUserEntity) (err error)
	DeleteTemporaryUser(ctx context.Context, email string) (err error)
	GetTemporaryUserByVerificationToken(ctx context.Context, token string, now int64) (user *TemporaryUserEntity, err error)
	GetTemporaryUserByIdentifier(ctx context.Context, identifier string, now int64) (user *TemporaryUserEntity, err error)

	GetUserByEmail(ctx context.Context, email string) (user *UserEntity, err error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (user *UserEntity, err error)
	GetUserByIdentifier(ctx context.Context, identifier string) (user *UserEntity, err error)

	InsertUserPasswordResetEntity(ctx context.Context, entity UserPasswordResetEntity) (err error)
	DeleteUserPasswordResetEntity(ctx context.Context, email string) (err error)
	GetUserPasswordResetEntity(ctx context.Context, token string, now int64) (res *UserPasswordResetEntity, err error)
}

func (p *UserEntity) VerifyPassword(providedPassword string) error {
	// Combine the provided password with the stored salt
	saltedPassword := providedPassword + *p.PasswordSalt

	// Compare the hashed version of the provided password with the stored hash
	return bcrypt.CompareHashAndPassword([]byte(*p.HashedPassword), []byte(saltedPassword))
}
