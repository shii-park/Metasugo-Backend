package test

import (
	"errors"
	"testing"

	"github.com/shii-park/Metasugo-Backend/internal/hub"
)

// SendJSONのバッファテスト
func TestClientSendJSONBuffer(t *testing.T) {
	c := hub.NewClient(nil, nil, "test-user")
	// バッファサイズ（256）を超えて送信を試み、少なくとも1回は失敗することを期待
	succeeded := 0
	var lastErr error
	for i := 0; i < 300; i++ {
		err := c.SendJSON(map[string]interface{}{"i": i})
		if err == nil {
			succeeded++
		} else {
			lastErr = err
			break
		}
	}
	if succeeded == 0 {
		t.Fatalf("成功した送信が期待されましたが、%d回でした", succeeded)
	}
	if lastErr == nil {
		t.Fatalf("送信バッファが満杯になることを期待しましたが、%d回成功しました", succeeded)
	}
	t.Logf("SendJSONは%d回成功した後、エラーを返しました: %v", succeeded, lastErr)
}

// SendJSONの正常動作テスト
func TestClientSendJSON_Success(t *testing.T) {
	c := hub.NewClient(nil, nil, "test-user")
	
	// 1つのメッセージを送信
	err := c.SendJSON(map[string]interface{}{
		"type":    "test",
		"message": "hello",
	})
	
	if err != nil {
		t.Errorf("SendJSONが失敗しました: %v", err)
	}
}

// SendJSONの不正なデータテスト
func TestClientSendJSON_InvalidData(t *testing.T) {
	c := hub.NewClient(nil, nil, "test-user-invalid")
	
	// JSONにシリアライズできないデータを送信
	invalidData := make(chan int)
	err := c.SendJSON(invalidData)
	
	if err == nil {
		t.Error("不正なデータでSendJSONがエラーを返さなかった")
	}
}

// SendErrorのテスト
func TestClientSendError(t *testing.T) {
	c := hub.NewClient(nil, nil, "test-user-error")
	
	// エラーメッセージを送信
	testErr := errors.New("テストエラー")
	err := c.SendError(testErr)
	
	if err != nil {
		t.Errorf("SendErrorが失敗しました: %v", err)
	}
}

// SendErrorのnilテスト
func TestClientSendError_Nil(t *testing.T) {
	c := hub.NewClient(nil, nil, "test-user-nil-error")
	
	// nilエラーを送信
	err := c.SendError(nil)
	
	if err != nil {
		t.Errorf("nilエラーでSendErrorが失敗しました: %v", err)
	}
}

// NewClientのテスト
func TestNewClient(t *testing.T) {
	testHub := hub.NewHub()
	client := hub.NewClient(testHub, nil, "test-user-new")
	
	if client == nil {
		t.Fatal("NewClientがnilを返しました")
	}
	
	// 複数のメッセージを送信してクライアントが正常に動作することを確認
	for i := 0; i < 5; i++ {
		err := client.SendJSON(map[string]interface{}{"count": i})
		if err != nil {
			t.Errorf("メッセージ%dの送信に失敗: %v", i, err)
		}
	}
}
