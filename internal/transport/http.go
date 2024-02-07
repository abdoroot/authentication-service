package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/abdoroot/authentication-service/internal/auth"
	"github.com/abdoroot/authentication-service/internal/types"
	"github.com/abdoroot/authentication-service/middleware"
)

type httpTransport struct {
	mux        *http.ServeMux
	srv        *auth.Auth
	listenAddr string
}

func NewHttpTransport(srv *auth.Auth, httpListenAddr string) *httpTransport {
	return &httpTransport{
		mux:        http.NewServeMux(),
		srv:        srv,
		listenAddr: httpListenAddr,
	}
}

func (t *httpTransport) Strart() {
	fmt.Printf("HttpTransport running on port%v\n", t.listenAddr)
	t.mux.HandleFunc("/login", MakeHttpTransportHandler(t.handelPostLoginUser))
	t.mux.HandleFunc("/user/create", MakeHttpTransportHandler(t.handelPostCreateUser))
	t.mux.HandleFunc("/user/update", middleware.HttpLoginMiddleware(MakeHttpTransportHandler(t.handelPostUpdateUser), t.srv)) //Login middleware
	http.ListenAndServe(t.listenAddr, t.mux)
}

func (t *httpTransport) handelPostCreateUser(w http.ResponseWriter, r *http.Request) error {
	param := &types.CreateUserParam{}
	if err := json.NewDecoder(r.Body).Decode(param); err != nil {
		return err
	}
	user, err := param.CreateUserFromParam()
	if err != nil {
		return err
	}
	createdUser, err := t.srv.SignUp(context.Background(), user)
	if err != nil {
		return err
	}
	WriteJson(w, http.StatusOK, createdUser)
	return nil
}

func (t *httpTransport) handelPostUpdateUser(w http.ResponseWriter, r *http.Request) error {
	param := &types.UpdateUserParam{}
	if err := json.NewDecoder(r.Body).Decode(param); err != nil {
		return err
	}
	user := getUserFromRequestCtx(r.Context())
	updateReq, err := param.CreateUpdateRequest()
	updateReq.ID = user.ID
	if err != nil {
		return err
	}
	updatedUser, err := t.srv.Update(context.Background(), updateReq)
	if err != nil {
		return err
	}
	WriteJson(w, http.StatusOK, updatedUser)
	return nil
}

func (t *httpTransport) handelPostLoginUser(w http.ResponseWriter, r *http.Request) error {
	param := &types.LoginParam{}
	if err := json.NewDecoder(r.Body).Decode(param); err != nil {
		return err
	}
	user, err := t.srv.Login(context.Background(), param)
	if err != nil {
		return err
	}
	userIdString := strconv.Itoa(user.ID)
	token, err := auth.GenerateToken(userIdString, user.Email)
	if err != nil {
		return err
	}
	WriteJson(w, http.StatusOK, map[string]any{
		"user":  user,
		"token": token,
	})
	return nil
}

func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Add("Content-type", "Application/json")
	w.WriteHeader(status)
	json, _ := json.Marshal(data)
	w.Write(json)
}

func MakeHttpTransportHandler(next types.HttpApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})
		}
	}
}

func getUserFromRequestCtx(ctx context.Context) *types.User {
	return ctx.Value("user").(*types.User)
}
