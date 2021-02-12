# ca-game-api

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
$ docker exec -i mysql-container sh -c 'exec mysql -u <上のユーザ名> -D ca_game_api_db -p"<上のユーザのパスワード>"' < init.sql
```