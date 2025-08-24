package handler

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gensan0223/snulog/internal/auth"
	"github.com/gensan0223/snulog/internal/repository"
	pb "github.com/gensan0223/snulog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WebHandler struct {
	grpcAddr    string
	authService *auth.AuthService
	userRepo    repository.UserRepository
}

func NewWebHandler(grpcAddr string, db *sql.DB) *WebHandler {
	return &WebHandler{
		grpcAddr:    grpcAddr,
		authService: auth.NewAuthService(),
		userRepo:    repository.NewPostgresUserRepository(db),
	}
}

func (h *WebHandler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	session, authenticated := h.authService.GetSessionFromRequest(r)
	if !authenticated {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Username string
	}{
		Username: session.Username,
	}

	tmpl.Execute(w, data)
}

func (h *WebHandler) ServeLogin(w http.ResponseWriter, r *http.Request) {
	// 既にログインしている場合はリダイレクト
	if _, authenticated := h.authService.GetSessionFromRequest(r); authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/login.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func (h *WebHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<div class="error-message">ユーザー名とパスワードを入力してください</div>`)
		return
	}

	user, err := h.userRepo.GetUserByUsername(username)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<div class="error-message">ユーザー名またはパスワードが間違っています</div>`)
		return
	}

	if !h.authService.CheckPassword(password, user.PasswordHash) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<div class="error-message">ユーザー名またはパスワードが間違っています</div>`)
		return
	}

	token, err := h.authService.CreateSession(username)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<div class="error-message">ログインに失敗しました</div>`)
		return
	}

	h.authService.SetSessionCookie(w, token)
	w.Header().Set("HX-Redirect", "/")
}

func (h *WebHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		h.authService.DeleteSession(cookie.Value)
	}

	h.authService.ClearSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *WebHandler) AddLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, authenticated := h.authService.GetSessionFromRequest(r)
	if !authenticated {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<div class="error-message">ログインが必要です</div>`)
		return
	}

	userName := session.Username // セッションからユーザー名を取得
	status := r.FormValue("status")
	feeling := r.FormValue("feeling")

	if status == "" || feeling == "" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<div class="error-message">すべてのフィールドを入力してください</div>`)
		return
	}

	conn, err := grpc.NewClient(h.grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<div class="error-message">サーバー接続エラー: %v</div>`, err)
		return
	}
	defer conn.Close()

	client := pb.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	entry := &pb.LogEntry{
		UserName:  userName,
		Status:    status,
		Feeling:   feeling,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	resp, err := client.AddLogs(ctx, entry)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<div class="error-message">ログの追加に失敗しました: %v</div>`, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<div class="success-message">✅ ログが正常に追加されました: %s</div>`, resp.Message)
}

func (h *WebHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	_, authenticated := h.authService.GetSessionFromRequest(r)
	if !authenticated {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<div class="error-message">ログインが必要です</div>`)
		return
	}

	conn, err := grpc.NewClient(h.grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<div class="error-message">サーバー接続エラー: %v</div>`, err)
		return
	}
	defer conn.Close()

	client := pb.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := client.FetchLogs(ctx, &pb.FetchRequest{
		TeamId: "default",
	})
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<div class="error-message">ログの取得に失敗しました: %v</div>`, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if len(resp.Logs) == 0 {
		fmt.Fprint(w, `<p>まだログがありません</p>`)
		return
	}

	for _, log := range resp.Logs {
		fmt.Fprintf(w, `
			<div class="log-entry">
				<strong>👤 %s</strong> - 📝 %s - 😀 %s
				<div class="log-meta">🕒 %s</div>
			</div>
		`, log.UserName, log.Status, log.Feeling, log.Timestamp)
	}
}
