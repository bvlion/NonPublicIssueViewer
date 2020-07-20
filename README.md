# 非公開 Issue Viewer

個人的に非公開の Issue をパスだけで公開したくなったので作りました。

## 目的

* GitHub で付けていた日々の食事の記録を非エンジニアでもそれとなく見れるサイトを用意する
* Go の勉強

## GAE リリース

GAE にリリースするための準備です。

### リリース用 JSON の権限

以下の権限を付与する  
<img src="https://user-images.githubusercontent.com/24517539/83963454-45c70400-a8e1-11ea-89d2-dee5e5320418.png">

### JSON

取得した json は base 64 エンコードして Secrets.SECRET_YAML_BASE64 に設定する

## GitHub

GitHub Actions を使ってデプロイを実施しています。

### Secrets

変数 | 説明
--- | ---
BG_URL | ログイン画面トップに設定する画像の URL
NOT_FOUND_URL | 404 ページで使う画像
GCP_PROJECT_ID | [setup-gcloud](https://github.com/GoogleCloudPlatform/github-actions/tree/master/setup-gcloud)を参照
GCP_SA_KEY | [setup-gcloud](https://github.com/GoogleCloudPlatform/github-actions/tree/master/setup-gcloud)を参照
SECRET_YAML_BASE64 | 上記で作成した base 64 エンコード された JSON
SLACK_CHANNEL_ID | ビルドを投稿する Slack チャンネル
SLACK_ICON_URL | Slack に使用するアイコン
SLACK_WEBHOOK_URL | Slack の Webhook URL

### 元データ

author の備忘録用です。

https://docs.google.com/spreadsheets/d/1dVHTc7CVa2ULY8xHKrgITgtFHJG1AEq5uEBKmaZQ7NA/edit

## コーディング

以下の内容で `source/secret.yaml` を準備します。

```
github:
  token: 
  user: 
  project: 

login_pass: 
session_key: 

footer_links:
  作成者ブログ(鰯の管詰): https://www.ambitious-i.net
 ```

js と css は gulp を使用して圧縮と配置を実施しています。

## ライセンス

MITです。  
`LICENSE.txt` をご確認ください。