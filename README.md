# dev-practice-todo

研修「開発実践編」で GitHub Flow を練習するための題材リポジトリです。

最小構成の TODO アプリ（タスクの追加・一覧・完了トグル・削除）を題材にしています。
**アプリそのものの作り込みは主題ではありません。** ブランチを切る・PR を出す・レビューを受ける・
バグを直すといった一連の開発フローを、実際に動くアプリの上で練習することが目的です。

## 構成

```
.
├── api/                 Go (net/http + chi) の API サーバ
│   ├── main.go          エントリポイント（DB 接続・マイグレーション・起動）
│   ├── handlers.go      ルーティングと各エンドポイントのハンドラ
│   ├── store.go         DB アクセス（タスクの CRUD）
│   ├── handlers_test.go ハンドラのテスト（DB 不要）
│   ├── schema.sql       tasks テーブル定義（冪等に適用）
│   └── Dockerfile
├── web/                 Next.js (App Router) の 1 ページ UI
│   ├── app/page.tsx     タスク一覧・追加・トグル・削除
│   └── Dockerfile
├── e2e/                 Playwright の E2E テスト
│   └── todo.spec.ts     ベースライン（追加→一覧→削除）
├── docs/
│   └── known-issues.md  既知の不具合メモ
├── docker-compose.yml   db + api + web をまとめて起動
└── .github/workflows/ci.yml
```

## セットアップ

Docker さえあれば動きます。

```bash
docker compose up --build
```

- UI: http://localhost:3000
- API ヘルスチェック: http://localhost:8080/healthz

ブラウザで http://localhost:3000 を開き、タスクの追加・削除ができることを確認してください。

停止と後片付け（DB のボリュームも削除）:

```bash
docker compose down -v
```

## API

| メソッド | パス               | 内容                              |
| -------- | ------------------ | --------------------------------- |
| GET      | `/tasks`           | 全タスクを作成日時の昇順で返す     |
| POST     | `/tasks`           | `{ "title": "..." }` を作成（201） |
| PATCH    | `/tasks/{id}`      | `{ "done": true }` で完了状態を更新 |
| DELETE   | `/tasks/{id}`      | タスクを削除（204）                |
| GET      | `/healthz`         | ヘルスチェック（200）              |

## テスト

### API のユニットテスト

```bash
cd api
go test ./...
```

### E2E（Playwright）

`docker compose up` でアプリを起動した状態で、別ターミナルから実行します。

```bash
cd e2e
pnpm install
pnpm exec playwright install --with-deps chromium
pnpm test
```

## 既知の不具合

このリポジトリには既知の不具合が **1 件** あります（[docs/known-issues.md](docs/known-issues.md)）。
完了トグルにチェックを付けてもリロードすると外れます。これは研修の課題で実際に直してもらう題材です。
