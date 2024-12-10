package auth

import (
	"context"
	"log"
	"mini-wallet/domain"
	_auth "mini-wallet/domain/auth"
	"mini-wallet/domain/common/response"
	"mini-wallet/utils"
	"net/http"
	"time"

	firebaseAdmin "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"google.golang.org/api/option"
)

type authHandler struct {
	authUsecase         _auth.AuthUsecase
	firebaseAdminClient *auth.Client
	config              *utils.AppConfig
}

func SetAuthHandler(
	router *chi.Mux,
	usecases domain.Usecases,
	middleware _auth.AuthMiddleware,
	config *utils.AppConfig,
) {
	app, err := firebaseAdmin.NewApp(context.Background(), nil, option.WithCredentialsFile(
		config.GoogleCredentialsPath,
	))
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
		return
	}

	authHandler := authHandler{
		authUsecase:         usecases.AuthUsecase,
		firebaseAdminClient: client,
		config:              config,
	}

	router.Route("/auth/verify", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/", authHandler.VerifyAccessToken)
	})

	router.Route("/auth/", func(r chi.Router) {
		// google handler
		r.Post("/register/google", authHandler.RegisterWithGoogle)
		r.Post("/login/google", authHandler.AuthenticateByGoogle)

		// from inquiry
		r.Post("/register/inquiry", authHandler.RegisterUserFromInquiry)
		r.Post("/login/inquiry", authHandler.AuthenticateFromInquiry)

		// pre authenticated
		r.Post("/check-identifier", authHandler.CheckIndentifier)
		r.Post("/login", authHandler.AuthenticateRegularUser)
		r.Post("/register", authHandler.RegisterUser)

		// authenticated
		r.Get("/logout", authHandler.Logout)

		// verifications
		r.Post("/verify-reset-password-token", authHandler.VerifyResetPasswordToken)
		r.Post("/verify-phone-number", authHandler.VerifyPhoneNumber)

		// passwords
		r.Post("/password-reset", authHandler.SendPasswordResetLink)
		r.Post("/password-reset-submission", authHandler.SubmitPasswordReset)

	})

	router.Route("/auth/refresh", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/", authHandler.RefreshAccess)
	})
}

func (handler *authHandler) AuthenticateFromInquiry(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := _auth.AuthFromInquiryDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	err := req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	res := handler.authUsecase.AuthenticateFromInquiry(context.Background(), req)
	res.Writer = w
	res.WriteResponse()
}

func (handler *authHandler) RegisterUserFromInquiry(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := _auth.AuthFromInquiryDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	err := req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	res := handler.authUsecase.RegisterUserFromInquiry(context.Background(), req)
	res.Writer = w
	res.WriteResponse()
}

func (handler *authHandler) RefreshAccess(w http.ResponseWriter, r *http.Request) {
	res := handler.authUsecase.RefreshAccess(r.Context())

	res.Writer = w
	res.WriteResponse()
}

func (handler *authHandler) CheckIndentifier(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := _auth.CheckIndentifierDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	err := req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	res := handler.authUsecase.CheckIdentifier(context.Background(), req)
	res.Writer = w
	res.WriteResponse()
}

func (handler *authHandler) VerifyResetPasswordToken(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := _auth.VerifyResetPasswordTokenDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	err := req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	res := handler.authUsecase.VerifyResetPasswordToken(context.Background(), req)
	res.Writer = w
	res.WriteResponse()
}

func (handler *authHandler) SubmitPasswordReset(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := _auth.PasswordResetSubmissionDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	err := req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	res := handler.authUsecase.ResetUserPassword(context.Background(), req)
	res.Writer = w
	res.WriteResponse()
}

func (handler *authHandler) VerifyPhoneNumber(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := _auth.VerifyEmailDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	err := req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	res := handler.authUsecase.VerifyPhoneNumber(context.Background(), req)
	res.Writer = w
	res.WriteResponse()
}

func (handler *authHandler) SendPasswordResetLink(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := _auth.PasswordResetDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	err := req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	res := handler.authUsecase.SendPasswordResetLink(context.Background(), req)
	res.Writer = w
	res.WriteResponse()
}

func (handler *authHandler) AuthenticateRegularUser(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := _auth.AuthenticationDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	err := req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	res := handler.authUsecase.AuthenticateRegularUser(context.Background(), req)
	res.Writer = w
	res.WriteResponse()
}

func (handler *authHandler) AuthenticateByGoogle(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := _auth.GoogleRegisterDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	err := handler.verifyToken(req.StsTokenManager.AccessToken)
	if err != nil {
		resp.Unauthorized(err.Error())
		return
	}

	res := handler.authUsecase.AuthenticateByGoogle(context.Background(), req)
	res.Writer = w
	res.WriteResponse()
}

func (handler *authHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := _auth.UserRegistrationDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	err := req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	res := handler.authUsecase.RegisterUser(context.Background(), req)
	res.Writer = w
	res.WriteResponse()
}

func (handler *authHandler) VerifyAccessToken(w http.ResponseWriter, r *http.Request) {
	res := response.Response[interface{}]{
		Writer: w,
	}

	tokenStatus := r.Context().Value(_auth.TokenStatus{}).(int)
	switch tokenStatus {
	case 0:
		res.Success(nil)
	case _auth.ERROR_INVALID_TOKEN:
		res.Unauthorized("InvalidToken")
	case _auth.ERROR_EXPIRED_TOKEN:
		res.Unauthorized("ExpiredToken")
	}

	res.WriteResponse()
}

func (handler *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	res := response.Response[_auth.AuthenticationResponse]{
		Writer: w,
	}
	now, _ := utils.GetJktTime()

	res.SuccessWithCookie("success", _auth.AuthenticationResponse{
		AccessToken:  "",
		RefreshToken: "",
	}, []*http.Cookie{
		{
			Name:     handler.config.AccessTokenKey,
			Value:    "",
			Domain:   ".sebia.id",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			Expires:  now.Add(time.Minute * -30),
		},
		{
			Name:     handler.config.RefreshTokenKey,
			Value:    "",
			Domain:   ".sebia.id",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			Expires:  now.Add(time.Minute * -30),
		},
	})

	res.WriteResponse()
}

func (handler *authHandler) RegisterWithGoogle(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := _auth.GoogleRegisterDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	err := handler.verifyToken(req.StsTokenManager.AccessToken)
	if err != nil {
		resp.Unauthorized(err.Error())
		return
	}

	res := handler.authUsecase.RegisterByGoogle(context.Background(), req)
	res.Writer = w
	res.WriteResponse()
}

func (handler *authHandler) verifyToken(token string) (err error) {
	_, err = handler.firebaseAdminClient.VerifyIDToken(context.Background(), token)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
		return err
	}

	return nil
}
