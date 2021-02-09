# ca-game-api

## 実行手順

1. `.env`ファイルを作成する

環境変数を指定するための`.env`ファイルをプロジェクトルートディレクトリに作成してください。
以下の<>で囲われた部分を適切な値に変更し、作成してください。

```env
DRIVER=mysql
DATASOURCE=<ユーザ名>:<ユーザのパスワード>@(mysql-container:3306)/ca_game_api_db
MYSQL_ROOT_PASSWORD=<上とは異なる任意のパスワード>
```

2. データベースを初期化し、APIを起動する

```bash
$ docker-compose up -d db

# ルートユーザでDBに接続する
# パスワードを求めらるのでMYSQL_ROOT_PASSWORD=で指定した値を打つ
$ docker-compose exec db mysql -u root -p

# ユーザとデータベースを作成し、ユーザに権限を与える
mysql> CREATE USER '<上のユーザ名>'@'%' IDENTIFIED BY '<上のユーザのパスワード>';
mysql> CREATE DATABASE ca_game_api_db;
mysql> GRANT ALL ON ca_game_api_db.* TO '<上のユーザ名>'@'%';
mysql> exit

# テーブルを作成する
$ docker exec -i mysql-container sh -c 'exec mysql -u <上のユーザ名> -D ca_game_api_db -p"<上のユーザのパスワード>"' < init.sql

# APIのコンテナを起動する。（-dオプションは任意でつける）
$ docker-compose up api
```