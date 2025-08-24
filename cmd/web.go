package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gensan0223/snulog/internal/handler"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		grpcAddr, _ := cmd.Flags().GetString("grpc-addr")

		// データベース接続
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			dsn = "postgres://postgres:password@localhost:5432/snulog?sslmode=disable"
		}

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		webHandler := handler.NewWebHandler(grpcAddr, db)

		// Static files
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

		// Routes
		http.HandleFunc("/", webHandler.ServeIndex)
		http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				webHandler.ServeLogin(w, r)
			} else if r.Method == http.MethodPost {
				webHandler.HandleLogin(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})
		http.HandleFunc("/logout", webHandler.HandleLogout)
		http.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				webHandler.AddLog(w, r)
			} else if r.Method == http.MethodGet {
				webHandler.GetLogs(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})

		fmt.Printf("🌐 Web server starting on http://localhost:%s\n", port)
		fmt.Printf("📡 Connecting to gRPC server at %s\n", grpcAddr)
		fmt.Printf("🔐 Login page: http://localhost:%s/login\n", port)

		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal("Failed to start web server:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().StringP("port", "p", "8080", "Port to run the web server on")
	webCmd.Flags().StringP("grpc-addr", "g", "localhost:50051", "gRPC server address")
}
