## 運用Tips / 画面・体験メモ

### 運用Tips

- TTL/GC: oauth_statesは10〜15分で削除。sessionsはexpires_atで定期削除（revoked=trueも考慮）。
- クッキー: `HttpOnly + Secure + SameSite=Lax(or None)`必須。
- 端末管理: sessionsにdevice情報を追加で「他端末ログアウト」も容易。

---

### 画面/体験メモ

- 初回ログイン後ダイアログ: 「あなたのメールドメインと一致するテナントが見つかりました: Kogakuin University。参加しますか？ [参加] [スキップ]」+「参加コードをお持ちの場合はこちら」入力欄。
- ヘッダーのテナント切替: 所属複数時はセレクタで `active_membership` を切替 → APIは常に `active_membership` をコンテキストに動作。
- スキップ時: `active_membership_id=NULL` でもアプリを使えるが、テナント限定機能は非表示/エンプティステート。

