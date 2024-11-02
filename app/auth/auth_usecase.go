package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"mini-wallet/domain"
	"mini-wallet/domain/auth"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/user"
	"mini-wallet/infrastructure"
	"mini-wallet/utils"
	"net/http"
	"time"
)

type authUsecase struct {
	userRepository user.UserRepository
}

func NewAuthUsecase(repositories domain.Repositories) auth.AuthUsecase {
	return &authUsecase{
		userRepository: repositories.UserRepository,
	}
}

func (usecase *authUsecase) CheckIdentifier(ctx context.Context, req auth.CheckIndentifierDTO) (res response.Response[string]) {
	user, err := usecase.userRepository.GetUserByIdentifier(ctx, req.Identifier)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if user == nil {
		res.NotFound("Pengguna tidak ditemukan", nil)
		return
	}

	res.Success("Pengguna ditemukan")
	return
}

func (usecase *authUsecase) VerifyResetPasswordToken(ctx context.Context, req auth.VerifyResetPasswordTokenDTO) (res response.Response[string]) {
	now, _ := utils.GetJktTime()
	passwordReset, err := usecase.userRepository.GetUserPasswordResetEntity(ctx, req.PasswordResetToken, now.Unix())
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if passwordReset == nil {
		res.BadRequest("Link sudah digunakan atau kedaluwarsa", nil)
		return
	}

	res.Success("Link dapat digunakan")
	return
}

func (usecase *authUsecase) ResetUserPassword(ctx context.Context, req auth.PasswordResetSubmissionDTO) (res response.Response[string]) {
	now, _ := utils.GetJktTime()
	passwordReset, err := usecase.userRepository.GetUserPasswordResetEntity(ctx, req.PasswordResetToken, now.Unix())
	if err != nil {
		res.InternalServerError(err.Error())
		log.Fatal(err)
		return
	}

	if passwordReset == nil {
		res.BadRequest("Link sudah digunakan atau kedaluwarsa", nil)
		return
	}

	user, err := usecase.userRepository.GetUserByEmail(ctx, passwordReset.Email)
	if err != nil {
		res.InternalServerError(err.Error())
		log.Fatal(err)
		return
	}

	if user == nil {
		res.BadRequest("Pengguna tidak ditemukan", nil)
		return
	}

	user.ChangePassword(req.Password)

	err = usecase.userRepository.UpsertUser(ctx, *user)
	if err != nil {
		log.Fatal(err)
		res.InternalServerError(err.Error())
		return
	}

	_ = usecase.userRepository.DeleteUserPasswordResetEntity(ctx, user.Email)

	res.SuccessWithMessage("Kata sandi diubah, silakan masuk")
	return
}

func (usecase *authUsecase) VerifyEmail(ctx context.Context, req auth.VerifyEmailDTO) (res response.Response[string]) {
	now, err := utils.GetJktTime()
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	temporaryUser, err := usecase.userRepository.GetTemporaryUserByVerificationToken(ctx, req.Token, now.Unix())
	if err != nil || temporaryUser == nil {
		res.BadRequest("Link sudah digunakan atau kedaluwarsa", nil)
		return
	}

	userEntity, err := temporaryUser.ToUserEntity()
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	err = usecase.userRepository.InsertUser(ctx, *userEntity)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	err = usecase.userRepository.DeleteTemporaryUser(ctx, userEntity.Email)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	res.SuccessWithMessage("Verifikasi email berhasil")
	return
}

func (usecase *authUsecase) SendPasswordResetLink(ctx context.Context, req auth.PasswordResetDTO) (res response.Response[string]) {
	existingUser, err := usecase.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if existingUser == nil {
		res.BadRequest("Pengguna tidak ditemukan", nil)
		return
	}

	now, _ := utils.GetJktTime()
	passwordResetToken, _ := GenerateRandomString(32)
	userPasswordResetEntity := user.UserPasswordResetEntity{
		Email:              existingUser.Email,
		ExpiredAt:          now.Add(time.Minute * 15).Unix(),
		PasswordResetToken: passwordResetToken,
	}

	err = usecase.userRepository.InsertUserPasswordResetEntity(ctx, userPasswordResetEntity)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	go infrastructure.SendPasswordResetLink(existingUser.Email, existingUser.Name, passwordResetToken)

	res.SuccessWithMessage("Instruksi atur ulang kata sandi terkirim")
	return res
}

func GenerateRandomString(n int) (string, error) {
	// Create a slice of random bytes
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Encode the bytes to a URL-safe base64 string
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}

func (usecase *authUsecase) AuthenticateRegularUser(ctx context.Context, req auth.AuthenticationDTO) (res response.Response[auth.AuthenticationResponse]) {
	now, err := utils.GetJktTime()
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	temporaryUser, err := usecase.userRepository.GetTemporaryUserByIdentifier(ctx, req.Identifier, now.Unix())
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if temporaryUser != nil {
		res.BadRequest("Verifikasi email terlebih dahulu", nil)
		return
	}

	existingUser, err := usecase.userRepository.GetUserByEmail(ctx, req.Identifier)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if existingUser == nil {
		existingUser, err = usecase.userRepository.GetUserByPhoneNumber(ctx, req.Identifier)
		if err != nil {
			res.InternalServerError(err.Error())
			return
		}

		if existingUser == nil {
			res.BadRequest("Pengguna tidak ditemukan", nil)
			return
		}
	}

	if existingUser.HashedPassword == nil {
		res.BadRequest("Silakan masuk menggunakan Google", nil)
		return
	}

	if existingUser.VerifyPassword(req.Password) != nil {
		res.BadRequest("Kata sandi salah", nil)
		return
	}

	accessToken, _ := auth.GenerateJWT(*existingUser, "ACCESS")
	refreshToken, _ := auth.GenerateJWT(*existingUser, "REFRESH")
	res.SuccessWithCookie("success", auth.AuthenticationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, []*http.Cookie{
		{
			Name:     "access_token",
			Value:    accessToken,
			Domain:   "tobacamping.id",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			Expires:  now.Add(time.Hour * 24 * 31),
		},
		{
			Name:     "refresh_token",
			Value:    accessToken,
			Domain:   "tobacamping.id",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			Expires:  now.Add(time.Hour * 24 * 31),
		},
	})

	return res
}

func (usecase *authUsecase) AuthenticateByGoogle(ctx context.Context, req auth.GoogleRegisterDTO) (res response.Response[auth.AuthenticationResponse]) {
	userEntity, err := req.ToUserEntity()
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	existingUser, err := usecase.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	accessToken, _ := auth.GenerateJWT(*userEntity, "ACCESS")
	refreshToken, _ := auth.GenerateJWT(*userEntity, "REFRESH")
	now, _ := utils.GetJktTime()
	if existingUser != nil {
		res.SuccessWithCookie("success", auth.AuthenticationResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, []*http.Cookie{
			{
				Name:     "access_token",
				Value:    accessToken,
				Domain:   "tobacamping.id",
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				Expires:  now.Add(time.Hour * 24 * 31),
			},
			{
				Name:     "refresh_token",
				Value:    accessToken,
				Domain:   "tobacamping.id",
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				Expires:  now.Add(time.Hour * 24 * 31),
			},
		})
		return
	}

	res.BadRequest("Akun belum terdaftar, lanjutkan pendaftaran", nil)
	return res
}

func (usecase *authUsecase) RegisterUser(ctx context.Context, req auth.UserRegistrationDTO) (res response.Response[interface{}]) {
	userEntity, err := req.ToTemporaryUserEntity()
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	existingUser, err := usecase.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if existingUser != nil {
		res.BadRequest("Email sudah digunakan", nil)
		return
	}

	existingUser, err = usecase.userRepository.GetUserByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if existingUser != nil {
		res.BadRequest("Nomor handphone sudah digunakan", nil)
		return
	}

	// insert to temporary user, expiring in 15 min
	err = usecase.userRepository.InsertTemporaryUser(ctx, *userEntity)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	go infrastructure.SendEmailVerificationLink(userEntity.Email, userEntity.Name, userEntity.VerificationToken)

	res.SuccessWithMessage("Link verifikasi dikirimkan ke email Anda")
	return res
}

func (usecase *authUsecase) RegisterByGoogle(ctx context.Context, req auth.GoogleRegisterDTO) (res response.Response[auth.AuthenticationResponse]) {
	userEntity, err := req.ToUserEntity()
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	existingUser, err := usecase.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	accessToken, _ := auth.GenerateJWT(*userEntity, "ACCESS")
	refreshToken, _ := auth.GenerateJWT(*userEntity, "REFRESH")
	now, _ := utils.GetJktTime()
	if existingUser != nil {
		res.SuccessWithCookie("success", auth.AuthenticationResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, []*http.Cookie{
			{
				Name:     "access_token",
				Value:    accessToken,
				Domain:   "tobacamping.id",
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				Expires:  now.Add(time.Hour * 24 * 31),
			},
			{
				Name:     "refresh_token",
				Value:    accessToken,
				Domain:   "tobacamping.id",
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				Expires:  now.Add(time.Hour * 24 * 31),
			},
		})
		return
	}

	err = usecase.userRepository.InsertUser(ctx, *userEntity)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	accessToken, _ = auth.GenerateJWT(*userEntity, "ACCESS")
	refreshToken, _ = auth.GenerateJWT(*userEntity, "REFRESH")

	res.SuccessWithCookie("success", auth.AuthenticationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, []*http.Cookie{
		{
			Name:     "access_token",
			Value:    accessToken,
			Domain:   "tobacamping.id",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			Expires:  now.Add(time.Hour * 24 * 31),
		},
		{
			Name:     "refresh_token",
			Value:    accessToken,
			Domain:   "tobacamping.id",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			Expires:  now.Add(time.Hour * 24 * 31),
		},
	})

	return res
}
