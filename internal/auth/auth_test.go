package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthService_HashPassword(t *testing.T) {
	auth := NewAuthService()

	password := "testpassword"
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" {
		t.Error("Expected non-empty hash")
	}

	if hash == password {
		t.Error("Hash should not equal original password")
	}
}

func TestAuthService_CheckPassword(t *testing.T) {
	auth := NewAuthService()

	password := "testpassword"
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// 正しいパスワード
	if !auth.CheckPassword(password, hash) {
		t.Error("Expected password to match hash")
	}

	// 間違ったパスワード
	if auth.CheckPassword("wrongpassword", hash) {
		t.Error("Expected wrong password to not match hash")
	}
}

func TestAuthService_SessionManagement(t *testing.T) {
	auth := NewAuthService()
	username := "testuser"

	// セッション作成
	token, err := auth.CreateSession(username)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}

	// セッション取得
	session, exists := auth.GetSession(token)
	if !exists {
		t.Error("Expected session to exist")
	}

	if session.Username != username {
		t.Errorf("Expected username %s, got %s", username, session.Username)
	}

	// セッション削除
	auth.DeleteSession(token)
	_, exists = auth.GetSession(token)
	if exists {
		t.Error("Expected session to be deleted")
	}
}

func TestAuthService_SessionExpiry(t *testing.T) {
	auth := NewAuthService()
	username := "testuser"

	token, err := auth.CreateSession(username)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// セッションの作成時間を過去に設定（期限切れをシミュレート）
	auth.mutex.Lock()
	auth.sessions[token].CreatedAt = time.Now().Add(-25 * time.Hour)
	auth.mutex.Unlock()

	// 期限切れセッションは取得できない
	_, exists := auth.GetSession(token)
	if exists {
		t.Error("Expected expired session to not exist")
	}
}

func TestAuthService_SessionCookie(t *testing.T) {
	auth := NewAuthService()

	// レスポンスレコーダー
	w := httptest.NewRecorder()
	token := "test-token"

	// クッキー設定
	auth.SetSessionCookie(w, token)

	// クッキーが設定されているかチェック
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Expected 1 cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != "session_token" {
		t.Errorf("Expected cookie name 'session_token', got %s", cookie.Name)
	}

	if cookie.Value != token {
		t.Errorf("Expected cookie value %s, got %s", token, cookie.Value)
	}

	if !cookie.HttpOnly {
		t.Error("Expected cookie to be HttpOnly")
	}
}

func TestAuthService_GetSessionFromRequest(t *testing.T) {
	auth := NewAuthService()
	username := "testuser"

	// セッション作成
	token, err := auth.CreateSession(username)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// リクエスト作成
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: token,
	})

	// リクエストからセッション取得
	session, exists := auth.GetSessionFromRequest(req)
	if !exists {
		t.Error("Expected session to exist")
	}

	if session.Username != username {
		t.Errorf("Expected username %s, got %s", username, session.Username)
	}

	// クッキーがないリクエスト
	reqNoCookie := httptest.NewRequest("GET", "/", nil)
	_, exists = auth.GetSessionFromRequest(reqNoCookie)
	if exists {
		t.Error("Expected no session for request without cookie")
	}
}
