# ca-game-api

## 概要

TechTrain の MISSION である[オンライン版　CA Tech Dojo サーバサイド (Go)編](https://techbowl.co.jp/techtrain/missions/12)に取り組んだものです. 
ミッション達成後にレスポンスボディの JSON のフィールドやデータベースに変更を加えています. 
オリジナルの機能としてユーザランキング機能, キャラクター合成機能を追加しました.

## セットアップ

1. `.env`ファイルを作成する

環境変数を指定するための`.env`ファイルをプロジェクトルートディレクトリに作成してください.
以下の`<>`で囲われた部分は適切な値に置き換えてください.

```test:.env
PORT=8080
DSN=postgres://<username>:<password>@ca-game-api-db-dev:5432/<dbname>?sslmode=disable

POSTGRES_PASSWORD=<password>
POSTGRES_USER=<username>
POSTGRES_DB=<dbname>
```

2. APIを起動する

以下のコマンドを実行し, 起動します.

```bash
$ docker compose up -d
```

## ドキュメント

実行できるエンドポイントなどは[こちら](https://minguu42.github.io/ca-game-api/)で確認できます.

## 開発環境

- OS: Mac OS
- プログラミング言語：Go
- 主なライブラリ・フレームワーク：net/http, database/sql
- データベース：PostgreSQL
- デプロイ：Heroku
- フォーマッタ・リンタ：gofmt, goimports, govet, staticcheck
- テスト：testing, net/http/httptest
- タスクランナー：GNU Make

### ローカル実行

```bash
make dev
# 終了時
make down
```

### ドキュメント

```bash
make docs
```

注意点：open コマンドでブラウザで所定の URL を開くようにしていますが, 環境によってうまく動きません.

### 自動整形

```bash
make fmt
```

注意点：goimports をインストールしている必要があります.

### 静的解析

```bash
make lint
```

注意点：staticcheck をインストールしている必要があります.

### テスト

```bash
make test
```

注意点：psql をインストールしている必要があります.
