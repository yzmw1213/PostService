# PostService
Goで構築するマイクロサービスの投稿機能サービス

# 概要
- 簡易的な記事投稿のCRUD機能

## 使用技術
- Go 1.12.17
- Docker docker-compose
- dockerize v0.6.1
- protoc 3.11.0
- gRPC v1.31.0
- AWS(VPC,ECS,ECR,RDS,ELB)
- terraform
- CircleCI

## 構成図
![AWS_stracture](https://user-images.githubusercontent.com/36359899/89097162-79bd3200-d417-11ea-83e5-8c998c824a0f.png)

## 機能一覧
- 投稿
  - 新規登録、編集、削除、全件取得
  - go-playground/validatorを用いたバリデーション
- サービス間通信
  - Envoyプロキシを介した他サービスとの通信
