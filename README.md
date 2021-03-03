# ca-game-api

## 概要

このリポジトリはTechTrainのMISSION[オンライン版　CA Tech Dojo サーバサイド (Go)編](https://techbowl.co.jp/techtrain/missions/12)に取り組んだものです。
ランキングや合成機能は自分のオリジナルの機能として作成しました。

## 動作

サンプルはcURLでリクエストを行い、全てローカルで確認できた動作を載せています。

### ユーザを作成する

ユーザ名を指定し、ユーザを作成できます。
`<username>`は適切な文字列に置き換えてください。
レスポンスとして、JSON形式でユーザ固有の認証トークンが返ります。
ユーザ名は重複せず、既に存在する場合は、`400 Bad Request`が返ります。

```bash
$ curl -X POST "http://localhost:8080/user/create" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"name\": \"<username>\"}"
{"token":"IOCECQAPBnkX9apWwvdcSJ"}
```

### ユーザを確認する

認証トークンで、ユーザ名を確認します。
`<token>`は適切な文字列に置き換えてください。
レスポンスとして、JSON形式でユーザ名が返ります。
一致するトークンが存在しない場合は、`401 Unauthorized`を返します。

```bash
$ curl -X GET "http://localhost:8080/user/get" -H  "accept: application/json" -H  "x-token: <token>"
{"name":"username"}
```

### ユーザ名を変更する

認証トークンで、ユーザ名を変更します。
`<token>`、`<username>`は適切な文字列に置き換えてください。
現在と同一の名前を送信しても、そのまま`200 OK`で返します。
指定した名前が既に存在する場合は、`400 Bad Request`を返します。

```bash
$ curl -X PUT "http://localhost:8080/user/update" -H  "accept: application/json" -H  "x-token: <token>" -H  "Content-Type: application/json" -d "{  \"name\": \"<username>\"}"
```

### ガチャを回す

認証トークンと回数を指定し、ガチャを回し、キャラクターを取得します。
`<token>`、`<times>`は適切な値に置き換えてください。
回数が正の整数でない場合は`400 Bad Request`を返します。

```bash
$ curl -X POST "http://localhost:8080/gacha/draw" -H  "accept: application/json" -H  "x-token: <token>" -H  "Content-Type: application/json" -d "{  \"times\": <times>}"
{
  "results": [
  
  ]
}
```

### ユーザ所持キャラクターを一覧取得する

認証トークンで、ユーザが所持しているキャラクターを一覧で取得します。

```bash
$ curl -X GET "http://localhost:8080/character/list" -H  "accept: application/json" -H  "x-token: <token>"
```

## 実行手順

1. `.env`ファイルを作成する

環境変数を指定するための`.env`ファイルをプロジェクトルートディレクトリに作成してください。
以下の<>で囲われた部分を適切な値に変更し、作成してください。

```env
DRIVER=mysql
DATASOURCE=<ユーザ名>:<ユーザのパスワード>@(mysql-container:3306)/ca_game_api_db
MYSQL_ROOT_PASSWORD=<上とは異なる任意のパスワード>
MYSQL_DATABASE=ca_game_api_db
MYSQL_USER=<上と等しいユーザ名>
MYSQL_PASSWORD=<上と等しいユーザのパスワード>
```

2. APIを起動する

```bash
$ docker-compose up -d

# （初回のみ、2回目以降は行わない）テーブルとキャラクターを作成する
$ docker exec -i mysql-container sh -c 'exec mysql -u <上と等しいユーザ名> -D ca_game_api_db -p"<上と等しいユーザのパスワード>"' < init.sql
```