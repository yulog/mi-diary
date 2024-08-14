# mi-diary

[![][mit-badge]][mit]
![GitHub go.mod Go version][go-version-badge]
[![GitHub Tag][tag-badge]][tag-url]
[![GitHub Release][release-badge]][release-url]

Misskey のリアクションの履歴を記録して振り返るツール。  
(ローカルで動く Misskey 版 favolog 的なものを目指す)

## ToDo

- ユーザー
  - [x] ユーザーごとのカウント
  - [x] ユーザーごとのノート一覧
- アーカイブ
  - [x] 日付ごとのカウント
  - [x] 日付ごとのノート一覧
    - [x] UTCになっているので直す？
- ノート
  - [x] ノートの本文保存、表示
  - [x] ノートにユーザー名表示
    - [x] アイコン表示
      - [ ] ローカルに保存する
  - [x] ノート一覧表示
  - [x] ノート一覧にページング
- リアクション
  - [x] リアクション画像
    - [ ] ローカルに保存する
- 添付ファイル
  - 画像
    - [x] ノートの添付画像表示
    - [x] 画像の形式の制限
    - [ ] ノートの添付画像保存
    - [x] 添付画像一覧
      - [x] ある画像を含むノート一覧
      - [x] 似た色で絞り込み
- その他
  - [x] 見た目
  - [x] リアクション履歴を全件取得する
  - [x] マルチプロファイル
  - [ ] 検索
    - [x] ノート本文(like)
  - [x] MiAuth
  - [ ] Wails

## License

MIT

## Author

yulog

[mit]:            http://opensource.org/licenses/MIT
[mit-badge]:      https://img.shields.io/badge/License-MIT-yellow.svg
[go-version-badge]:https://img.shields.io/github/go-mod/go-version/yulog/mi-diary
[tag-badge]:https://img.shields.io/github/v/tag/yulog/mi-diary
[release-badge]:https://img.shields.io/github/v/release/yulog/mi-diary
[tag-url]:https://github.com/yulog/mi-diary/tags
[release-url]:https://github.com/yulog/mi-diary/releases