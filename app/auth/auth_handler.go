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
)

type authHandler struct {
	authUsecase         _auth.AuthUsecase
	firebaseAdminClient *auth.Client
}

func SetAuthHandler(router *chi.Mux, usecases domain.Usecases) {
	app, err := firebaseAdmin.NewApp(context.Background(), nil)
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
	}

	router.Route("/auth/", func(r chi.Router) {
		r.Post("/register/google", authHandler.RegisterWithGoogle)
		r.Post("/login/google", authHandler.AuthenticateByGoogle)
		r.Post("/check-identifier", authHandler.CheckIndentifier)
		r.Post("/login", authHandler.AuthenticateRegularUser)
		r.Post("/password-reset", authHandler.SendPasswordResetLink)
		r.Post("/password-reset-submission", authHandler.SubmitPasswordReset)
		r.Post("/verify-reset-password-token", authHandler.VerifyResetPasswordToken)

		r.Post("/verify-email", authHandler.VerifyEmail)

		r.Post("/register", authHandler.RegisterUser)

		r.Get("/logout", authHandler.Logout)
		r.Get("/", authHandler.VerifyAccessToken)
		r.Get("/refresh", authHandler.Logout)

	})
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

func (handler *authHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
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

	res := handler.authUsecase.VerifyEmail(context.Background(), req)
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

	cookie, err := r.Cookie("access_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "No cookie found", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
		return
	}

	token := cookie.Value

	_, errCode := _auth.ValidateToken(token)
	switch errCode {
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
			Name:     "access_token",
			Value:    "",
			Domain:   "tobacamping.id",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			Expires:  now.Add(time.Minute * -30),
		},
		{
			Name:     "refresh_token",
			Value:    "",
			Domain:   "tobacamping.id",
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
