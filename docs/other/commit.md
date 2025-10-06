以下の規約に従うとコミットメッセージが見やすくなる

https://www.conventionalcommits.org/ja/v1.0.0/

## 1. 全体の目的と背景

**Conventional Commits** は、Git のコミットメッセージに一貫したフォーマットを定める軽量の規約で、人間にもツールにも意味が明示される履歴を作ることを目的としています。

これにより：

- 変更内容（機能追加・バグ修正・破壊的変更など）を一目で分かるようにする
- 自動で CHANGELOG を生成したり、Semantic Versioning（SemVer）に基づいてバージョンを上げたりするツールと連携しやすくする

という利点が生まれます。[Conventional Commits](https://www.conventionalcommits.org/ja/v1.0.0/?utm_source=chatgpt.com)

（分かりやすい日本語での解説もあり、初心者にも「なぜやるか／どう書くか」が整理されている。[Zenn](https://zenn.dev/wakamsha/articles/about-conventional-commits?utm_source=chatgpt.com)）

---

## 2. コミットメッセージの構造（必須／任意のパーツ）

基本形：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

日本語訳だと：

```
<型>[任意のスコープ]: <タイトル>

[任意の本文]

[任意のフッター]
```

[Conventional Commits](https://www.conventionalcommits.org/ja/v1.0.0/?utm_source=chatgpt.com)

各要素の意味：

- **type（型）**：コミットのカテゴリを示す短いキーワード。例：`feat`、`fix`、`docs` など。後述の標準的な型が推奨される。
- **scope（スコープ）**（任意）：変更の対象領域を括弧で囲んで示す。例：`feat(auth): ...` や `fix(ui): ...` のように。
- **description（タイトル）**：そのコミットの要約（短く、命令形が多い）。
- **body（本文）**（任意）：なぜその変更をしたか、背景や詳細、トレードオフ、補足情報。
- **footer（フッター）**（任意）：関連する issue 番号や破壊的変更の明示など。特別なフォーマット（例：`BREAKING CHANGE: ...`）もここに入れる。

---

## 3. 標準的な type（型）の一覧と用途

主に推奨されている型（カスタマイズも可能だが、チームで合意しておくとよい）：

| type | 意味・用途例 |
| --- | --- |
| `feat` | 新機能の追加。これがあると SemVer ではマイナーバージョン上昇の候補になる。 |
| `fix` | バグ修正。パッチバージョンアップ。 |
| `docs` | ドキュメントのみの変更（コードの挙動には影響しない）。 |
| `style` | フォーマットや空白、セミコロンなどのコード整形（機能に影響しない）。 |
| `refactor` | バグ修正でも機能追加でもないコードの変更（構造改善など）。 |
| `perf` | パフォーマンス向上のための変更（機能は変わらない）。 |
| `test` | テスト追加・修正（実装を変更しない）。 |
| `build` | ビルド関連（例：依存関係、ビルドスクリプトの変更）。 |
| `ci` | CI 設定やスクリプトの変更（テストやビルド自体の構成）。 |
| `chore` | その他の雑務的な変更（例：ライブラリ更新だがビルドには関係しない）。 |
| `revert` | 以前のコミットの取り消し。 |

[Conventional Commits](https://www.conventionalcommits.org/ja/v1.0.0/?utm_source=chatgpt.com)[enlume.com](https://www.enlume.com/blogs/mastering-commit-messages-a-guide-to-conventional-commits-and-best-practices/?utm_source=chatgpt.com)

---

## 4. 破壊的変更（Breaking Changes）の扱い

破壊的変更を含む場合、以下のどちらか／両方で明示する：

1. **ヘッダーに `!` を追加**
    
    例：`feat!: drop support for Node 10`
    
2. **フッターに明示的に書く**
    
    ```
    feat(parser): change API
    
    BREAKING CHANGE: the `parse()` function now returns a promise instead of value.
    ```
    

破壊的変更があると SemVer のメジャーバージョンを上げるトリガーになる。[Conventional Commits](https://www.conventionalcommits.org/ja/v1.0.0/?utm_source=chatgpt.com)

---

## 5. 具体例

### 単純な新機能

```
feat(auth): add OAuth2 login flow

ユーザーが Google / GitHub を使ってログインできるようにした。
セッション管理は JWT を用い、期限は 7 日に設定。
```

### バグ修正

```
fix(api): handle nil pointer in user fetch

user が存在しないときに nil dereference していたのを防ぐためのチェックを追加。
```

### 破壊的変更

```
feat!: remove v1 endpoints

BREAKING CHANGE: all API endpoints under /v1 are removed; clients must migrate to /v2.
```

### ドキュメントのみ

```
docs: update README with setup instructions
```

---

## 6. 付加的な運用ポイント

- **複数の変更を含むなら分割を検討する**：一つのコミットに `feat` と `fix` の両方が本質的にある場合は、可能なら別コミットに分ける。[Conventional Commits](https://www.conventionalcommits.org/ja/v1.0.0-beta.4/?utm_source=chatgpt.com)
- **Issue 参照**：フッターで `Closes #123` や `Refs #456` の形で関連 issue をつなげるのが一般的（これは GitHub などと連携して自動クローズすることもある）。
- **自動化との連携**：たとえば `release-please` や `semantic-release` などのツールはこれらの type/破壊的変更を解析して自動で CHANGELOG を生成したり、バージョンタグを付けたりする。[Zenn](https://zenn.dev/kuritify/articles/conventional-commits-definitive-guide?utm_source=chatgpt.com)[enlume.com](https://www.enlume.com/blogs/mastering-commit-messages-a-guide-to-conventional-commits-and-best-practices/?utm_source=chatgpt.com)

---

## 7. 導入のベストプラクティス

1. **チームで型一覧とルールを決める**（必要なら `type` を拡張）。[Conventional Commits](https://www.conventionalcommits.org/ja/v1.0.0-beta.4/?utm_source=chatgpt.com)
2. **コミットメッセージのlint を導入する**（たとえば `commitlint` + `husky` など）。[Zenn](https://zenn.dev/wakamsha/articles/about-conventional-commits?utm_source=chatgpt.com)
3. **CI やリリースパイプラインに組み込む**（SemVer に従ったタグ付け・CHANGELOG 自動生成）。[Zenn](https://zenn.dev/kuritify/articles/conventional-commits-definitive-guide?utm_source=chatgpt.com)
4. **PR／レビュー時にも一貫性をチェック**（例：タイトルを Conventional Commits 準拠にすることでレビューワーが意図を理解しやすくなる）。[enlume.com](https://www.enlume.com/blogs/mastering-commit-messages-a-guide-to-conventional-commits-and-best-practices/?utm_source=chatgpt.com)

---

## 8. まとめ：何が得られるか

- 履歴が**機械でも解釈可能**な構造になる
- 自動リリース（バージョニング／CHANGELOG）ができる
- 開発者間の**意図共有が明確化**される
- 後から履歴を追うときの**可読性と信頼性が上がる**

公式の語り口はシンプルだが、周辺ツールと組み合わせることで運用がぐっと楽になる。[Conventional Commits](https://www.conventionalcommits.org/ja/v1.0.0/?utm_source=chatgpt.com)[Zenn](https://zenn.dev/wakamsha/articles/about-conventional-commits?utm_source=chatgpt.com)[Zenn](https://zenn.dev/kuritify/articles/conventional-commits-definitive-guide?utm_source=chatgpt.com)

---

必要なら、あなたのプロジェクトに適用するための `.commitlintrc` 設定例、GitHub Actions での自動リリースワークフロー雛形、既存コミット履歴の一括変換スクリプト（例：`git rebase -i` を使った書き換え方）なども出せます。どうしますか？

ほかにも以下のようなものがある(CLI上からコミットを行うとき有効)

- https://github.com/streamich/git-cz
- https://github.com/commitizen/cz-cli
