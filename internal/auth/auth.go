package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	Username  string
	CreatedAt time.Time
}

type AuthService struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
}

func NewAuthService() *AuthService {
	return &AuthService{
		sessions: make(map[string]*Session),
	}
}

func (a *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (a *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *AuthService) GenerateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (a *AuthService) CreateSession(username string) (string, error) {
	token, err := a.GenerateSessionToken()
	if err != nil {
		return "", err
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.sessions[token] = &Session{
		Username:  username,
		CreatedAt: time.Now(),
	}

	return token, nil
}

func (a *AuthService) GetSession(token string) (*Session, bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	session, exists := a.sessions[token]
	if !exists {
		return nil, false
	}

	// セッションの有効期限チェック（24時間）
	if time.Since(session.CreatedAt) > 24*time.Hour {
		delete(a.sessions, token)
		return nil, false
	}

	return session, true
}

func (a *AuthService) DeleteSession(token string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	delete(a.sessions, token)
}

func (a *AuthService) GetSessionFromRequest(r *http.Request) (*Session, bool) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, false
	}
	return a.GetSession(cookie.Value)
}

func (a *AuthService) SetSessionCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   24 * 60 * 60, // 24時間
		HttpOnly: true,
		Secure:   false, // 開発環境ではfalse
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

func (a *AuthService) ClearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
