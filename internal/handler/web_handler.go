package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	pb "github.com/gensan0223/snulog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WebHandler struct {
	grpcAddr string
}

func NewWebHandler(grpcAddr string) *WebHandler {
	return &WebHandler{
		grpcAddr: grpcAddr,
	}
}

func (h *WebHandler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func (h *WebHandler) AddLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userName := r.FormValue("user_name")
	status := r.FormValue("status")
	feeling := r.FormValue("feeling")

	if userName == "" || status == "" || feeling == "" {
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
