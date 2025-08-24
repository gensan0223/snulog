#!/bin/bash

echo "🚀 Snulog セットアップスクリプト"

# .envファイルから環境変数を読み込み
if [ -f .env ]; then
    echo "📄 .envファイルを読み込み中..."
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "❌ .envファイルが見つかりません"
    exit 1
fi

# PostgreSQLコンテナが起動しているかチェック
echo "🔍 PostgreSQLコンテナの状態をチェック中..."
if ! docker ps | grep -q snulog-db; then
    echo "🐳 PostgreSQLコンテナを起動中..."
    docker-compose up -d db
    
    # ヘルスチェック待機
    echo "⏳ データベースの準備を待機中..."
    timeout=60
    while [ $timeout -gt 0 ]; do
        if docker exec snulog-db pg_isready -U postgres > /dev/null 2>&1; then
            echo "✅ データベースが準備完了"
            break
        fi
        sleep 1
        timeout=$((timeout-1))
    done
    
    if [ $timeout -eq 0 ]; then
        echo "❌ データベースの起動がタイムアウトしました"
        exit 1
    fi
else
    echo "✅ PostgreSQLコンテナは既に起動しています"
fi

# マイグレーション実行
echo "🔄 データベースマイグレーションを実行中..."
docker-compose run --rm migrate

echo "🎉 セットアップ完了！"
echo ""
echo "次のコマンドでサーバーを起動できます："
echo "  gRPCサーバー: go run server/server.go"
echo "  Webサーバー:  go run main.go web"
echo "  デバッグ:     go run main.go debug"