# ca-game-api

## 概要

TechTrainのMISSIONである[オンライン版　CA Tech Dojo サーバサイド (Go)編](https://techbowl.co.jp/techtrain/missions/12)に取り組んだものです。
ランキングや合成機能は自分のオリジナルの機能として追加で作成しました。

## 実行手順

1. `.env`ファイルを作成する

環境変数を指定するための`.env`ファイルをプロジェクトルートディレクトリに作成してください。
以下の`<>`で囲われた部分を適切な値に変更し、作成してください。

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
$ docker exec -i ca-game-api-db sh -c 'exec psql -U minguu -d ca_game_api_db -w' < init.sql
```

## 動作例

動作例は全てローカルで、cURLでリクエストを行い、確認した動作を載せています。

### ユーザを作成する

ユーザ名を指定し、ユーザを作成できます。
`name`の値は適切な値に変更してください。
レスポンスとして、JSON形式でユーザ固有の認証トークンが返ります。
ユーザ名は重複できず、既に存在する場合は、`400 Bad Request`が返ります。

```bash
$ curl -i -X POST "http://localhost:8000/user/create" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"name\": \"minguu\"}"
HTTP/1.1 200 OK
Date: Sat, 06 Mar 2021 10:33:06 GMT
Content-Length: 40
Content-Type: text/plain; charset=utf-8

{
  "token": "ceKeMPeYr0eF3K5e4Lfjfe"
}
```

### ユーザを確認する

認証トークンで、ユーザ名を確認します。
`x-token`の値は適切な値に変更してください。
レスポンスとして、JSON形式でユーザ名が返ります。
トークンが適切でない場合は、`401 Unauthorized`を返します。

```bash
$ curl -i -X GET "http://localhost:8000/user/get" -H  "accept: application/json" -H  "x-token: ceKeMPeYr0eF3K5e4Lfjfe"
HTTP/1.1 200 OK
Date: Sat, 06 Mar 2021 10:34:15 GMT
Content-Length: 23
Content-Type: text/plain; charset=utf-8

{
  "name": "minguu"
}
```

### ユーザ名を変更する

認証トークンで、ユーザ名を変更します。
`x-token`, `name`の値は適切な値に変更してください。
現在と同一の名前を指定しても、そのまま`200 OK`で返します。
指定した名前が既に存在する場合は、`400 Bad Request`を返します。

```bash
$ curl -i -X PUT "http://localhost:8000/user/update" -H  "accept: application/json" -H  "x-token: ceKeMPeYr0eF3K5e4Lfjfe" -H  "Content-Type: application/json" -d "{  \"name\": \"newMinguu\"}"
HTTP/1.1 200 OK
Date: Sat, 06 Mar 2021 10:38:46 GMT
Content-Length: 0
```

### ガチャを回す

認証トークンと回数を指定し、ガチャを回し、キャラクターを取得します。
`x-token`、`times`の値は適切な値に変更してください。
レスポンスとして、JSON形式でガチャの結果が返ります。
回数が正の整数でない場合は`400 Bad Request`を返します。

```bash
$ curl -i -X POST "http://localhost:8000/gacha/draw" -H  "accept: application/json" -H  "x-token: ceKeMPeYr0eF3K5e4Lfjfe" -H  "Content-Type: application/json" -d "{  \"times\": 3}"
HTTP/1.1 200 OK
Date: Sat, 06 Mar 2021 10:40:30 GMT
Content-Length: 259
Content-Type: text/plain; charset=utf-8

{
  "results": [
    {
      "characterID": "40000010",
      "name": "rare_character10"
    },
    {
      "characterID": "30000004",
      "name": "normal_character4"
    },
    {
      "characterID": "40000006",
      "name": "rare_character6"
    }
  ]
}
```

### ユーザ所持キャラクターを一覧取得する

認証トークンで、ユーザが所持しているキャラクターを一覧で取得します。
`x-token`の値は適切な値に変更してください。
レスポンスとして、JSON形式で所有しているキャラクター一覧が返ります。
キャラクタを所有していないユーザの場合は、 ステータス行は`200 OK`で、 レスポンスボディは`{"characters": null}`で返します。

```bash
$ curl -i -X GET "http://localhost:8000/character/list" -H  "accept: application/json" -H  "x-token: ceKeMPeYr0eF3K5e4Lfjfe"
HTTP/1.1 200 OK
Date: Sat, 06 Mar 2021 10:42:47 GMT
Content-Length: 407
Content-Type: text/plain; charset=utf-8

{
  "characters": [
    {
      "userCharacterID": "1",
      "characterID": "40000010",
      "name": "rare_character10",
      "level": 3
    },
    {
      "userCharacterID": "2",
      "characterID": "30000004",
      "name": "normal_character4",
      "level": 2
    },
    {
      "userCharacterID": "3",
      "characterID": "40000006",
      "name": "rare_character6",
      "level": 10
    }
  ]
}
```

### ユーザのランキングを取得する

ユーザのランキングを取得します。
レスポンスとして、JSON形式でユーザのランキングが返ります。
ユーザのランキングは、合計値の高い方から3人のユーザを返します。
ユーザのランキングは、ユーザの所有しているキャラクターのレベルとキャラクター固有のパワーの合計値で決定します。

```bash
$ curl -i -X GET "http://localhost:8000/ranking/user" -H  "accept: application/json"                                                                            
HTTP/1.1 200 OK
Date: Sat, 06 Mar 2021 10:48:59 GMT
Content-Length: 274
Content-Type: text/plain; charset=utf-8

{
  "userRankings": [
    {
      "userID": "3",
      "name": "user2",
      "sumPower": "5420"
    },
    {
      "userID": "2",
      "name": "user1",
      "sumPower": "5380"
    },
    {
      "userID": "1",
      "name": "minguu",
      "sumPower": "4740"
    }
  ]
}
```

### キャラクターを合成する

認証トークン、ベースとなるキャラクターのID、合成されるキャラクターのIDで、キャラクターを合成します。
キャラクターを合成することでキャラクターのレベルを上げることができます。
`x-token`、`baseUserCharacterID`、`materialUserCharacterID`の値は適切な値に変更してください。
レスポンスとして、JSON形式で合成後のキャラクターが返ります。
所有していないキャラクターを指定した場合は、`400 Bad Request`を返します。

```bash
$ curl -i -X PUT "http://localhost:8080/character/compose" -H  "accept: application/json" -H  "x-token: -nRs7IX1H2dRiPttorkAL5" -H  "Content-Type: application/json" -d "{  \"baseUserCharacterID\": 4,  \"materialUserCharacterID\": 5}"
HTTP/1.1 200 OK
Date: Sat, 06 Mar 2021 10:53:34 GMT
Content-Length: 103
Content-Type: text/plain; charset=utf-8

{
  "userCharacterID": "4",
  "characterID": "30000009",
  "name": "normal_character9",
  "level": 4
}
```


curl -i -X POST "http://localhost:8000/user/create" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{  \"name\": \"minguu\"}"

curl -i -X GET "http://localhost:8000/user/get" -H  "accept: application/json" -H  "x-token: PYi05uOulU4Rshqz5YZQ-c"

curl -i -X PUT "http://localhost:8000/user/update" -H  "accept: application/json" -H  "x-token: PYi05uOulU4Rshqz5YZQ-c" -H  "Content-Type: application/json" -d "{  \"name\": \"newMinguu\"}"