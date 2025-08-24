#!/bin/bash

# .envファイルから環境変数を読み込み
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

echo "🚀 Snulogサーバーを起動中..."
echo "📡 gRPCサーバー: localhost:50051"
echo "🌐 Webサーバー: http://localhost:8080"
echo ""

# バックグラウンドでgRPCサーバーを起動
echo "Starting gRPC server..."
go run server/server.go &
GRPC_PID=$!

# 少し待ってからWebサーバーを起動
sleep 2
echo "Starting web server..."
go run main.go web &
WEB_PID=$!

# 終了時にプロセスをクリーンアップ
cleanup() {
    echo ""
    echo "🛑 サーバーを停止中..."
    kill $GRPC_PID 2>/dev/null
    kill $WEB_PID 2>/dev/null
    exit 0
}

trap cleanup SIGINT SIGTERM

echo "✅ サーバーが起動しました"
echo "Ctrl+C で停止します"

# プロセスが終了するまで待機
wait