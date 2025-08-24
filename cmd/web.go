package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gensan0223/snulog/internal/handler"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		grpcAddr, _ := cmd.Flags().GetString("grpc-addr")

		webHandler := handler.NewWebHandler(grpcAddr)

		// Static files
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

		// Routes
		http.HandleFunc("/", webHandler.ServeIndex)
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
