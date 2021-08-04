# ca-game-api

## 概要

TechTrain の MISSION である[オンライン版　CA Tech Dojo サーバサイド (Go)編](https://techbowl.co.jp/techtrain/missions/12)に取り組んだものです. 
ミッション達成後にレスポンスボディの JSON のフィールドやデータベースに変更を加えています. 
オリジナルの機能としてユーザランキング機能, キャラクター合成機能を追加しました.

## 実行手順

1. `.env`ファイルを作成する

環境変数を指定するための`.env`ファイルをプロジェクトルートディレクトリに作成してください.
以下の`<>`で囲われた部分は適切な値に置き換えてください.

```test:.env
PORT=8080
DSN=postgres://<username>:<password>@ca-game-api-db:5432/<dbname>?sslmode=disable

POSTGRES_PASSWORD=<password>
POSTGRES_USER=<username>
POSTGRES_DB=<dbname>
```

2. APIを起動する

以下のコマンドを実行し, 起動します.

```bash
$ docker compose up -d
```

## 動作例

動作例は全てローカルで, cURLでリクエストを行い, 確認した動作を載せています.

### ユーザを作成する

ユーザ名を指定しユーザを作成できます.
レスポンスとして認証トークンが返されます.

```bash
$ curl -X POST "http://localhost:8080/user/create" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"name\": \"minguu\"}"
{
  "token": "2JMZe9atOCkE8q0YH5s-Wr"
}
```

### ユーザを確認する

認証トークンでユーザ名を確認できます.
レスポンスとしてユーザ名が返されます.

```bash
$ curl -X GET "http://localhost:8080/user/get" -H  "accept: application/json" -H  "x-token: 2JMZe9atOCkE8q0YH5s-Wr"
{
  "name": "minguu2"
}
```

### ユーザ名を変更する

認証トークンでユーザ名を変更できます.

```bash
$ curl -X PUT "http://localhost:8080/user/update" -H  "accept: application/json" -H  "x-token: 2JMZe9atOCkE8q0YH5s-Wr" -H  "Content-Type: application/json" -d "{  \"name\": \"new_minguu\"}"
$ curl -X GET "http://localhost:8080/user/get" -H  "accept: application/json" -H  "x-token: 2JMZe9atOCkE8q0YH5s-Wr" 
{
  "name": "new_minguu"
}
```

### ガチャを回す

認証トークンで回数を指定してガチャを回し, キャラクターを取得できます.
レスポンスとしてガチャの結果が返されます.

```bash
$ curl -X POST "http://localhost:8080/gacha/draw" -H  "accept: application/json" -H  "x-token: 2JMZe9atOCkE8q0YH5s-Wr" -H  "Content-Type: application/json" -d "{  \"times\": 3}"
{
  "results": [
    {
      "characterID": 30000006,
      "name": "normal_character6"
    },
    {
      "characterID": 30000001,
      "name": "normal_character1"
    },
    {
      "characterID": 30000002,
      "name": "normal_character2"
    }
  ]
}
```

### ユーザ所持キャラクターを一覧取得する

認証トークンで所持しているキャラクターを一覧で取得できます.
レスポンスとして所有しているキャラクター一覧が返されます.

```bash
$ curl -X GET "http://localhost:8080/character/list" -H  "accept: application/json" -H  "x-token: 2JMZe9atOCkE8q0YH5s-Wr"
{
  "characters": [
    {
      "userCharacterID": 1,
      "characterID": 30000006,
      "name": "normal_character6",
      "level": 2,
      "experience": 700,
      "power": 410
    },
    {
      "userCharacterID": 2,
      "characterID": 30000001,
      "name": "normal_character1",
      "level": 1,
      "experience": 100,
      "power": 1
    },
    {
      "userCharacterID": 3,
      "characterID": 30000002,
      "name": "normal_character2",
      "level": 2,
      "experience": 400,
      "power": 400
    }
  ]
}
```

### ユーザのランキングを取得する

認証トークンでユーザのランキングを取得できます.
レスポンスとしてユーザのランキングが返されます.
ユーザのランキングはユーザの所有しているキャラクターのレベルとキャラクター固有のパワーの合計値で決定し, 高い方から3人のユーザが返されます.

```bash
$ curl -X 'GET' 'http://localhost:8080/user/ranking' -H 'accept: application/json' -H 'x-token: 2JMZe9atOCkE8q0YH5s-Wr'                         
{
  "users": [
    {
      "name": "example",
      "sumPower": 300000
    },
    {
      "name": "minguu",
      "sumPower": 223600
    },
    {
      "name": "example2",
      "sumPower": 204000
    }
  ]
}
```

### キャラクターを合成する

認証トークン, ベースとなるキャラクターのID, 合成されるキャラクターのIDでキャラクターを合成できます. 
レスポンスとして合成後のキャラクターが返されます.
キャラクターを合成することでキャラクターのレベルを上げられます.

```bash
$ curl -X PUT "http://localhost:8080/character/compose" -H  "accept: application/json" -H  "x-token: 2JMZe9atOCkE8q0YH5s-Wr" -H  "Content-Type: application/json" -d "{  \"baseUserCharacterID\": 1,  \"materialUserCharacterID\": 2}"
{
  "userCharacterID": 1,
  "characterID": 30000006,
  "name": "normal_character6",
  "level": 2,
  "experience": 700,
  "power": 410
}
```