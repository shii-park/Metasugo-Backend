package test

import (
"testing"
"time"

"github.com/shii-park/Metasugo-Backend/internal/hub"
)

// NewHubのテスト
func TestNewHub(t *testing.T) {
h := hub.NewHub()

if h == nil {
t.Fatal("NewHubがnilを返しました")
}

t.Log("Hubが正常に作成されました")
}

// Hubのクライアント登録と登録解除のテスト
func TestHub_RegisterUnregister(t *testing.T) {
h := hub.NewHub()
go h.Run()

// テストクライアントを作成
client1 := hub.NewClient(h, nil, "user1")
client2 := hub.NewClient(h, nil, "user2")

// クライアントを登録
h.Register(client1)
h.Register(client2)

// 登録が処理されるまで待機
time.Sleep(50 * time.Millisecond)

t.Log("2つのクライアントを登録しました")

// クライアントを登録解除
h.Unregister(client1)

// 登録解除が処理されるまで待機
time.Sleep(50 * time.Millisecond)

t.Log("1つのクライアントを登録解除しました")

// 残りのクライアントを登録解除
h.Unregister(client2)

// クリーンアップを待機
time.Sleep(50 * time.Millisecond)
}

// Hubの実行テスト
func TestHub_Run(t *testing.T) {
h := hub.NewHub()

// Hubを別のゴルーチンで実行
go h.Run()

// 短時間待機
time.Sleep(100 * time.Millisecond)

t.Log("Hubが正常に実行されています")
}

// 複数クライアントの同時登録テスト
func TestHub_MultipleClients(t *testing.T) {
h := hub.NewHub()
go h.Run()

// 複数のクライアントを作成して登録
clients := make([]*hub.Client, 5)
for i := 0; i < 5; i++ {
clients[i] = hub.NewClient(h, nil, string(rune('A'+i)))
h.Register(clients[i])
}

// 登録が処理されるまで待機
time.Sleep(100 * time.Millisecond)

t.Log("5つのクライアントを登録しました")

// すべてのクライアントを登録解除
for _, client := range clients {
h.Unregister(client)
}

// クリーンアップを待機
time.Sleep(100 * time.Millisecond)

t.Log("すべてのクライアントを登録解除しました")
}
