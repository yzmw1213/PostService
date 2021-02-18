# PostService
Goで構築するマイクロサービスの投稿サービス

## 使用技術
- Go 1.15
- Docker,docker-compose
- dockerize v0.6.1
- protoc 3.11.0
- gRPC v1.31.0
- AWS(IAM,VPC,ECS,ECR,RDS,Route53,ELB,S3,CloudWatch)
- terraform
- CircleCI

## 構成図
![PortfolioArchitecture](https://user-images.githubusercontent.com/36359899/108287470-2ef5e280-71ce-11eb-9301-a2c3c8ed5d01.png)

## 機能一覧
- 投稿
  - 新規登録、編集、削除、全件取得、検索
  - タグ登録
  - 投稿タグ付け
  - 投稿お気に入り機能
  - 投稿コメント
  - go-playground/validatorを用いたバリデーション
- サービス間通信
  - Envoyプロキシを介した他サービスとの通信

## アピールポイント
1. マイクロサービスアーキテクチャを採用している
2. gRPCでサービス間通信を行っている
3. テストコードを書いている
4. interfaceを書いてメソッドの実装チェックを行っている
5. linterを使っている
6. issueとプルリクエストを活用している

## 関連レポジトリ
- [フロントエンド](https://github.com/yzmw1213/Front)
- [Envoyプロキシ](https://github.com/yzmw1213/Proxy)
- [ユーザーサービス](https://github.com/yzmw1213/UserService)
