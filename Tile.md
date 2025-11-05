# `tiles.json` 仕様書

このドキュメントは、すごろくゲームの盤面を構成するタイル（マス）のデータを定義する `tiles.json` ファイルの形式について説明します。

## 1. 基本構造

`tiles.json` は、以下の構造を持つオブジェクトの配列です。各オブジェクトが一つのタイルを表します。

```json
{
  "id": 1,
  "kind": "profit",
  "detail": "給料日！",
  "effect": {
    "type": "profit",
    "amount": 100
  },
  "prev_ids": [0],
  "next_ids": [2]
}
```

### フィールド説明

- `id` (number, required): タイルを一意に識別するためのID。
- `kind` (string, required): タイルの大まかな種類を示す識別子。詳細は `3. kindの種類` を参照。
- `detail` (string, required): タイルの説明文。UIに表示されます。
- `effect` (object, required): プレイヤーがこのタイルに止まった際に発生する効果を定義するオブジェクト。詳細は `4. effectオブジェクトの詳細` を参照。
- `prev_ids` (array of numbers, required): このタイルにつながる前のタイルのIDの配列。
- `next_ids` (array of numbers, required): このタイルからつながる次のタイルのIDの配列。

---

## 2. `kind`の種類

`kind` はタイルの種類を分類するための識別子です。`effect` オブジェクトの `type` と密接に関連しています。

| `kind`        | 説明                               |
| :------------ | :--------------------------------- |
| `profit`      | プレイヤーの所持金が増えるマス。   |
| `loss`        | プレイヤーの所持金が減るマス。     |
| `quiz`        | クイズが出題されるマス。           |
| `branch`      | プレイヤーが次に進む道を選択するマス。 |
| `gamble`      | ギャンブルに挑戦するマス。         |
| `conditional` | プレイヤーの状態で効果が変わるマス。 |
| `overall`     | 全員に影響する効果が発生するマス。 |
| `neighbor`    | 周囲のプレイヤーに影響するマス。   |
| `require`     | 特定の条件が必要なマス。           |
| `goal`        | ゴールマス。                       |
| (その他)      | `effect` が `null` または `{}` の場合は効果なしマス。 |

---

## 3. `effect`オブジェクトの詳細

`effect` オブジェクトは、タイルの具体的な効果を定義します。必ず `type` フィールドを持ち、その値に応じて必要なフィールドが異なります。`type` は通常 `kind` と同じ値を指定します。

### 3.1. `profit`

所持金を増やします。

- `type`: `"profit"`
- `amount` (number): 増える金額。

```json
"effect": {
  "type": "profit",
  "amount": 100
}
```

### 3.2. `loss`

所持金を減らします。

- `type`: `"loss"`
- `amount` (number): 減る金額。

```json
"effect": {
  "type": "loss",
  "amount": 100
}
```

### 3.3. `quiz`

クイズを出題します。正解すると `amount` がプラスされ、不正解だとマイナスされます。

- `type`: `"quiz"`
- `quiz_id` (number): `quizzes.json` に定義されているクイズのID。
- `amount` (number): 賞金または罰金の額。

```json
"effect": {
  "type": "quiz",
  "quiz_id": 1,
  "amount": 50
}
```

### 3.4. `branch`

プレイヤーに進む道を選択させます。選択肢は `next_ids` に基づいて自動的に生成されるため、`effect` オブジェクトに追加のフィールドは不要です。

- `type`: `"branch"`

```json
"effect": {
  "type": "branch"
}
```

### 3.5. `gamble`

High & Low ギャンブルを発生させます。

- `type`: `"gamble"`

```json
"effect": {
  "type": "gamble"
}
```

### 3.6. `overall`

自分以外の全プレイヤーに影響を与えます。

- `type`: `"overall"`
- `profit_amount` (number, optional): 他のプレイヤーから徴収する金額。
- `loss_amount` (number, optional): 他のプレイヤーに配る金額。

```json
"effect": {
  "type": "overall",
  "profit_amount": 10
}
```

### 3.7. `neighbor`

前後1マスにいるプレイヤーに影響を与えます。

- `type`: `"neighbor"`
- `profit_amount` (number, optional): 周囲のプレイヤーから徴収する金額。
- `loss_amount` (number, optional): 周囲のプレイヤーに配る金額。

```json
"effect": {
  "type": "neighbor",
  "loss_amount": 20
}
```

### 3.8. `require`

（現在この `effect` の具体的な動作は未定義です）

- `type`: `"require"`
- `require_value` (number)
- `amount` (number)

### 3.9. `goal` / `NoEffect`

ゴールマスや、何も効果がないマスです。追加のフィールドは不要です。

- `type`: `"goal"` または `null`

```json
"effect": {
  "type": "goal"
}
// または
"effect": null
```

### 3.10. `conditional`

プレイヤーの特定の状態によって、適用される効果が変わります。

- `type`: `"conditional"`
- `condition` (string): 判定する条件。以下のいずれかを指定します。
  - `"isMarried"`: 結婚しているか
  - `"hasChildren"`: 子供がいるか
  - `"isProfessor"`: 職業がコース長か
  - `"isLecturer"`: 職業が平教員か
- `true_effect` (object): 条件が真の場合に適用される `effect` オブジェクト。
- `false_effect` (object): 条件が偽の場合に適用される `effect` オブジェクト。`null` を指定すると何も起きません。

`true_effect` と `false_effect` には、このドキュメントで説明されている任意の `effect` オブジェクトを指定でき、`conditional` を入れ子にすることも可能です。

**例1: 結婚している場合のみ効果が発生**
```json
"effect": {
  "type": "conditional",
  "condition": "isMarried",
  "true_effect": {
    "type": "profit",
    "amount": 1000
  },
  "false_effect": null
}
```

**例2: 職業によって効果が変わる（入れ子）**
```json
"effect": {
  "type": "conditional",
  "condition": "isProfessor",
  "true_effect": {
    "type": "profit",
    "amount": 500
  },
  "false_effect": {
    "type": "conditional",
    "condition": "isLecturer",
    "true_effect": {
      "type": "loss",
      "amount": 200
    },
    "false_effect": null
  }
}
```
