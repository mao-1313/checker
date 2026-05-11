# 課題：並行URLチェッカー

複数のURLに同時にHTTPリクエストを送り、各URLのステータスを返すツール。

```
入力: []string{"https://example.com", "https://golang.org", ...}
出力: 各URLのステータスコードまたはエラー
```

## 学習ステップ

### Step 1 - 直列で実装する

goroutine なしでまず動くものを作る。

```go
for _, url := range urls {
    // 1件ずつリクエスト
}
```

### Step 2 - goroutine で並列化する

`go` キーワードで関数を別goroutineで起動する。

```go
for _, url := range urls {
    go func(u string) {
        // 並行リクエスト
    }(url)
}
```

### Step 3 - channel で結果を受け取る

goroutine の結果を channel 経由で収集する。

```go
ch := make(chan Result, len(urls))
go func(u string) {
    ch <- check(u)
}(url)
```

### Step 4 - テストを書く

`httptest.NewServer` でモックサーバーを立てることで、外部URLに依存しないテストが書ける。

```go
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}))
defer srv.Close()
// srv.URL を使ってテスト
```

## 進め方

Step 1 の直列実装から始める。`checker.go` を新規作成する形で良い。
