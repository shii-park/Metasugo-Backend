# API

# WebSocket API仕様書

このドキュメントは、クライアント（フロントエンド）とサーバー（バックエンド）がWebSocketを介してリアルタイム通信を行う際に使用するJSONメッセージのフォーマットを定義します。

## 基本的なメッセージフォーマット

送受信されるすべてのメッセージは、以下の共通フォーマットに従います。

```json
{  
	"type": "メッセージ種別名",  
	"payload": {    ...  }
}
```

- `type`: 送信されるメッセージの種類を識別するための文字列です。
- `payload`: メッセージ種別に応じたデータが含まれるオブジェクトです。

---

## Websocketコネクション

## `/ws/connection`

でサーバとのコネクションを確立します。

Websocketには認証が必要なので

### `/ws/connection?token=Firebaseのトークン`

のようにアクセスしないと拒否られると思います。

---

## === クライアント → サーバーへのメッセージ ===

クライアントからサーバーへ送信されるメッセージです。

### `ROLL_DICE`

現在のプレイヤーがサイコロを振って駒を動かす際に送信します。

- **`type`**: `ROLL_DICE`
- **`payload`**: (空)特になんのデータもいりません。

**例:**

```json
{  
	"type": "ROLL_DICE",  
	"payload": {}
}
```

### `SUBMIT_CHOICE`

プレイヤーが分岐マスで進む方向を選択した際に送信します。これは、サーバーからの `BRANCH_CHOICE_REQUIRED` メッセージへの応答です。

- **`type`**: `SUBMIT_CHOICE`
- **`payload`**:
    - `selection` (数値): プレイヤーが移動先として選択したタイルのID。

**例:**

```json
{  
	"type": "SUBMIT_CHOICE",  
	"payload": {    
		"selection": 6  
	}
}
```

### `SUBMIT_QUIZ`

プレイヤーがクイズの回答を選択した際に送信します。これは、サーバーからの `QUIZ_REQUIRED` メッセージへの応答です。

- **`type`**: `SUBMIT_QUIZ`
- **`payload`**:
    - `selection` (数値): プレイヤーが選択した選択肢のインデックス（0から始まる）。

**例:**

```json
{  
	"type": "SUBMIT_QUIZ",  
	"payload": {    
		"selection": 0  
	}
}
```

### `SUBMIT_GAMBLE`

プレイヤーがギャンブルマスでの賭けの内容を決定した際に送信します。これは、サーバーからの `GAMBLE_REQUIRED` メッセージへの応答です。

- **`type`**: `SUBMIT_GAMBLE`
- **`payload`**:
    - `bet` (数値): プレイヤーが賭ける金額。
    - `choice` (文字列): プレイヤーの選択 (`"High"` または `"Low"`)。

**例:**

```json
{  
	"type": "SUBMIT_GAMBLE",  
	"payload": {    
		"bet": 50,    
		"choice": "High"  
	}
}
```

---

## === サーバー → クライアントへのメッセージ ===

サーバーから一人または複数のクライアントへ送信されるメッセージです。

### `PLAYER_MOVED`

プレイヤーがマスからマスへ移動した際に、すべてのクライアントに通知されます。

- **`type`**: `PLAYER_MOVED`
- **`payload`**:
    - `userID` (文字列): 移動したプレイヤーのID。
    - `newPosition` (数値): プレイヤーが新たに移動した先のタイルID。

**例:**

```json
{  
	"type": "PLAYER_MOVED",  
	"payload": {    
		"userID": "player1",    
		"newPosition": 5  
	}
}
```

### `MONEY_CHANGED`

プレイヤーの所持金が変動した際に、すべてのクライアントに通知されます。

- **`type`**: `MONEY_CHANGED`
- **`payload`**:
    - `userID` (文字列): 所持金が変動したプレイヤーのID。
    - `newMoney` (数値): プレイヤーの新しい所持金総額。

**例:**

```json
{  
	"type": "MONEY_CHANGED",  
	"payload": {    
		"userID": "player1",    
		"newMoney": 150  
	}
}
```

### `DICE_RESULT`

プレイヤーがサイコロを振った結果を通知します。

- **`type`**: `DICE_RESULT`
- **`payload`**:
    - `userID` (文字列): サイコロを振ったプレイヤーのID。
    - `diceResult` (数値): サイコロの目の数。

**例:**

```json
{  
	"type": "DICE_RESULT",  
	"payload": {    
		"userID": "player1",    
		"diceResult": 4  
	}
}
```

### `BRANCH_CHOICE_REQUIRED`

プレイヤーが分岐マスに止まり、選択が必要になった際に、対象のクライアントに送信されます。

- **`type`**: `BRANCH_CHOICE_REQUIRED`
- **`payload`**:
    - `tileID` (数値): プレイヤーがいる分岐マスのタイルID。
    - `options` (数値の配列): プレイヤーが選択可能な移動先のタイルIDのリスト。

**例:**

```json
{  
	"type": "BRANCH_CHOICE_REQUIRED",  
	"payload": {    
		"tileID": 5,    
		"options": [6, 7]  
	}
}
```

### `QUIZ_REQUIRED`

プレイヤーがクイズマスに止まった際に、対象のクライアントにクイズ情報を送信します。

- **`type`**: `QUIZ_REQUIRED`
- **`payload`**:
    - `tileID` (数値): プレイヤーがいるクイズマスのタイルID。
    - `quizData` (オブジェクト): クイズの詳細情報。
        - `id` (数値): クイズID。
        - `question` (文字列): 問題文。
        - `options` (文字列の配列): 選択肢のリスト。
        - `answer_description` (文字列): 正解・不正解時に表示する解説文。

**例:**

```json
{  
	"type": "QUIZ_REQUIRED",  
	"payload": {    
		"tileID": 9,    
		"quizData": {      
		"id": 1,      
		"question": "日本の首都は？",      
		"options": ["大阪", "京都", "東京"],      
		"answer_description": "正解は東京です。"    }  }}
```

### `GAMBLE_REQUIRED`

プレイヤーがギャンブルマスに止まった際に、対象のクライアントに選択を要求します。

- **`type`**: `GAMBLE_REQUIRED`
- **`payload`**:
    - `tileID` (数値): プレイヤーがいるギャンブルマスのタイルID。
    - `referenceValue` (数値): High/Lowの基準となる値。

**例:**

```json
{  
	"type": "GAMBLE_REQUIRED",  
	"payload": {    
		"tileID": 10,    
		"referenceValue": 3  
	}
}
```

### `GAMBLE_RESULT`

ギャンブルの結果を全クライアントに通知します。

- **`type`**: `GAMBLE_RESULT`
- **`payload`**:
    - `userID` (文字列): ギャンブルを行ったプレイヤーのID。
    - `diceResult` (数値): サイコロの目の合計。
    - `choice` (文字列): プレイヤーの選択 (`"High"` または `"Low"`)。
    - `won` (真偽値): プレイヤーが勝ったかどうか。
    - `amount` (数値): 変動した金額。
    - `newMoney` (数値): ギャンブル後の最終的な所持金。

**例:**

```json
{  
	"type": "GAMBLE_RESULT",  
	"payload": {    
		"userID": "player1",    
		"diceResult": 5,    
		"choice": "High",    
		"won": true,    
		"amount": 50,    
		"newMoney": 200  
	}
}
```

### `PLAYER_FINISHED`

プレイヤーがゴールした際に、全クライアントに通知します。

- **`type`**: `PLAYER_FINISHED`
- **`payload`**:
    - `userID` (文字列): ゴールしたプレイヤーのID。
    - `money` (数値): ゴール時の最終所持金。

**例:**

```json
{  
	"type": "PLAYER_FINISHED",  
	"payload": {    
		"userID": "player1",    
		"money": 500  
	}
}
```

### `PLAYER_STATUS_CHANGED`

プレイヤーのステータス（結婚、子供、職業など）が変化した際に、全クライアントに通知します。

- **`type`**: `PLAYER_STATUS_CHANGED`
- **`payload`**:
    - `userID` (文字列): ステータスが変化したプレイヤーのID。
    - `status` (文字列): 変化したステータスの種類 (`"isMarried"`, `"hasChildren"`, `"job"`)。
    - `value` (任意): 変化後の新しい値 (`true`, `"professor"` など)。

**例:**

```json
{
  "type": "PLAYER_STATUS_CHANGED",
  "payload": {
    "userID": "player1",
    "status": "isMarried",
    "value": true
  }
}
```

### `ERROR`

プレイヤーのアクションがエラーになったり、不正なメッセージを送信したりした場合に、対象のクライアントに送信されます。

- **`type`**: `ERROR`
- **`payload`**:
    - `message` (文字列): 発生したエラーの内容を説明するメッセージ。

**例:**

```json
{  
	"type": "ERROR",  
	"payload": {    
		"message": "無効なリクエストです。"  
	}
}
```

---

# HTTP API仕様書

このセクションでは、WebSocket以外の方法で提供されるAPIについて記述します。

## ランキングAPI (`/ranking`)

### `GET /ranking`

- **説明:** 全プレイヤーのランキングを取得します。
- **認証:** 必要
- **レスポンス:**
    - `200 OK`:
    `json [ { "playerID": "player1", "money": 1000, "finishedAt": "2023-10-27T10:00:00Z" }, { "playerID": "player2", "money": 900, "finishedAt": "2023-10-27T10:01:00Z" } ]`

// 全プレイヤーのデータは流石にグロいので何かしら対策するかも

工場の排煙や自動車の廃棄バスから発生したNOxと揮発性有機化合物が太陽光の紫外線によって反応し、二次的に生成されるオゾンやPANなどの酸化生成物の総称である。