package sugoroku

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Effect interface {
	Apply(player *Player, game *Game, choice any) error // 効果の適用
	RequiresUserInput() bool                            // ユーザからの入力が必要かどうか
	GetOptions(tile *Tile) any                          // ユーザの入力の選択肢
}

type effectWithType struct {
	Type TileKind `json:"type"`
}

// 収入マス
type ProfitEffect struct {
	Amount int `json:"amount"`
}

func (e ProfitEffect) RequiresUserInput() bool { return false }

func (e ProfitEffect) GetOptions(tile *Tile) any { return nil }

// 指定されたお金分増やす
func (e ProfitEffect) Apply(p *Player, g *Game, choice any) error {
	err := p.Profit(e.Amount)
	return err
}

// 支出マス
type LossEffect struct {
	Amount int `json:"amount"`
}

func (e LossEffect) RequiresUserInput() bool { return false }

func (e LossEffect) GetOptions(tile *Tile) any { return nil }

// 指定されたお金分減らす
func (e LossEffect) Apply(p *Player, g *Game, choice any) error {
	err := p.Loss(e.Amount)
	return err
}

// クイズマス

// クイズの問題集のパス
var QuizJSONPath = "./quizzes.json"

// クイズ問題の構造体
type Quiz struct {
	ID                int      `json:"id"`
	Question          string   `json:"question"`
	Options           []string `json:"options"`
	AnswerIndex       int      `json:"answerIndex"`
	AnswerDescription string   `json:"answer_description"`
}

// グローバル変数にキャッシュしておく
var quizzes []Quiz

// クイズタイルに必要なオプション
type QuizEffect struct {
	QuizID int `json:"quiz_id"`
	Amount int `json:"amount"`
}

// クイズマスにユーザからの入力が必要かどうか
func (e QuizEffect) RequiresUserInput() bool { return true }

// クイズIDからクイズを取ってきて、そのクイズを返す
func (e QuizEffect) GetOptions(tile *Tile) any {
	for _, quiz := range quizzes {
		if quiz.ID == e.QuizID {
			return quiz
		}
	}
	return nil
}

// クイズの実際の処理
func (e QuizEffect) Apply(p *Player, g *Game, choice any) error {
	// テスト時にfloat64型でもらうため、int型に変換している
	var selectedOptionIndex int
	switch v := choice.(type) {
	case int:
		selectedOptionIndex = v
	case float64:
		selectedOptionIndex = int(v)
	default:
		return fmt.Errorf("invalid choice for quiz: unexpected type %T", v)
	}

	var targetQuiz *Quiz
	for i := range quizzes {
		if quizzes[i].ID == e.QuizID {
			targetQuiz = &quizzes[i]
			break
		}
	}

	if targetQuiz == nil {
		return fmt.Errorf("quiz with ID %d not found", e.QuizID)
	}

	if selectedOptionIndex == targetQuiz.AnswerIndex {
		p.Profit(e.Amount)
	} else {
		p.Loss(e.Amount)
	}

	return nil
}

// 分かれ道マス
type BranchEffect struct {
}

// 分かれ道にユーザからの入力が必要かどうか
func (e BranchEffect) RequiresUserInput() bool { return true }

// ユーザの選択肢。次のマスを取得して、それを戻り値にしている。
func (e BranchEffect) GetOptions(tile *Tile) any {
	options := make([]int, len(tile.nexts))
	for i, nextTile := range tile.nexts {
		options[i] = nextTile.id
	}
	return options
}

// 選ばれたマスの方へ進めている。
func (e BranchEffect) Apply(p *Player, g *Game, choice any) error {
	var chosenTileID int
	// float64 と int の両方の型に対応
	switch v := choice.(type) {
	case int:
		chosenTileID = v
	case float64:
		chosenTileID = int(v)
	default:
		return fmt.Errorf("invalid choice for branch: unexpected type %T", v)
	}

	// 選択肢が現在の分岐マスの次のタイルとして有効か検証する
	isValidChoice := false
	for _, nextTile := range p.GetPosition().nexts {
		if nextTile.GetID() == chosenTileID {
			isValidChoice = true
			break
		}
	}

	if !isValidChoice {
		return fmt.Errorf("invalid choice for branch: tile %d is not a valid next tile", chosenTileID)
	}

	// プレイヤーの位置を選択されたタイルに更新
	nextTile, exists := g.tileMap[chosenTileID]
	if !exists {
		// このエラーは上のバリデーションにより通常発生しないはず
		return errors.New("chosen tile does not exist")
	}

	p.position = nextTile
	return nil
}

// 全体効果
type OverallEffect struct {
	ProfitAmount int `json:"profit_amount"`
	LossAmount   int `json:"loss_amount"`
}

func (e OverallEffect) RequiresUserInput() bool { return false }

func (e OverallEffect) GetOptions(tile *Tile) any { return nil }

// 　全員にお金を配るもしくはお金をもらう
func (e OverallEffect) Apply(p *Player, g *Game, choice any) error {
	allPlayers := g.GetAllPlayers()

	//自分以外のプレイヤーを取得する
	otherPlayers := make([]*Player, 0, len(allPlayers)-1)
	for _, player := range allPlayers {
		if player.id != p.id {
			otherPlayers = append(otherPlayers, player)
		}
	}

	if e.ProfitAmount > 0 {
		// 全体にお金をもらう
		p.Profit(len(otherPlayers) * e.ProfitAmount)
		LossForTargetPlayers(otherPlayers, e.ProfitAmount)
	} else if e.LossAmount > 0 {
		// 全員にお金を配る
		p.Loss(len(otherPlayers) * e.LossAmount)
		ProfitForTargetPlayers(otherPlayers, e.LossAmount)
	} else {
		return errors.New("invalid amount for overall effect")
	}
	return nil
}

// 隣人効果
type NeighborEffect struct {
	ProfitAmount int `json:"profit_amount"`
	LossAmount   int `json:"loss_amount"`
}

func (e NeighborEffect) RequiresUserInput() bool { return false }

func (e NeighborEffect) GetOptions(tile *Tile) any { return nil }

// 周辺(前後1マス)のプレイヤーからお金をもらうもしくは配る
func (e NeighborEffect) Apply(p *Player, g *Game, choice any) error {
	targetPlayers := g.GetNeighbors(p)
	if e.ProfitAmount > 0 {
		// 全体にお金をもらう
		p.Profit(len(targetPlayers) * e.ProfitAmount)
		LossForTargetPlayers(targetPlayers, e.ProfitAmount)
	} else if e.LossAmount > 0 {
		// 全員にお金を配る
		p.Loss(len(targetPlayers) * e.LossAmount)
		ProfitForTargetPlayers(targetPlayers, e.LossAmount)
	} else {
		return errors.New("invalid amount for overall effect")
	}
	return nil
}

// 条件分岐
type RequireEffect struct {
	RequireValue int `json:"require_value"`
	Amount       int `json:"amount"`
}

func (e RequireEffect) RequiresUserInput() bool { return false }

func (e RequireEffect) GetOptions(tile *Tile) any { return nil }

func (e RequireEffect) Apply(p *Player, g *Game, choice any) error {

	return nil
}

// ギャンブル効果
type GambleEffect struct {
}

// ConditionalEffect はプレイヤーのステータスに基づいて異なる効果を適用します。

type ConditionalEffect struct {
	Condition   string          `json:"condition"`    // "isMarried", "hasChildren"
	TrueEffect  json.RawMessage `json:"true_effect"`  // 条件がtrueの場合の効果
	FalseEffect json.RawMessage `json:"false_effect"` // 条件がfalseの場合の効果
}

func (e ConditionalEffect) RequiresUserInput() bool {
	// TODO: 将来的に、内部のEffectがユーザー入力を要求する可能性を考慮し拡張する
	return false
}

func (e ConditionalEffect) GetOptions(tile *Tile) any {
	return nil
}

func (e ConditionalEffect) Apply(p *Player, g *Game, choice any) error {
	var conditionMet bool
	switch e.Condition {
	case "isMarried":
		conditionMet = p.GetIsMarried()
	case "children":
		conditionMet = p.GetChildren() > 0
	case "isProfessor":
		conditionMet = p.GetJob() == JobProfessor
	case "isLecturer":
		conditionMet = p.GetJob() == JobLecturer
	default:
		return fmt.Errorf("unknown condition: %s", e.Condition)
	}

	var effectJSON json.RawMessage
	if conditionMet {
		effectJSON = e.TrueEffect
	} else {
		effectJSON = e.FalseEffect
	}

	if effectJSON == nil || string(effectJSON) == "null" || string(effectJSON) == "{}" {
		return nil // 適用する効果がない場合は何もしない
	}

	// effectJSONからEffectインスタンスを動的に生成して適用
	effect, err := CreateEffectFromJSON(effectJSON)
	if err != nil {
		return err
	}

	return effect.Apply(p, g, choice)
}

func (e GambleEffect) RequiresUserInput() bool { return true }

func (e GambleEffect) GetOptions(tile *Tile) any { return nil }

// ギャンブルの入力の有効か検証している
// 本当はここにギャンブルの処理を書いて、returnでギャンブル結果を返したほうが良いのだろうが、時間がなかったので呼び出し先でギャンブルの判定を行っている。TODO: リファクタリングが必要
func (e GambleEffect) Apply(p *Player, g *Game, choice any) error {
	userInput, ok := choice.(map[string]interface{})
	if !ok {
		return errors.New("invalid input format for gamble")
	}

	bet, ok := userInput["bet"].(float64)
	if !ok {
		return errors.New("bet is missing or not a number")
	}
	if int(bet) <= 0 {
		return errors.New("bet must be positive")
	}

	choiceStr, ok := userInput["choice"].(string)
	if !ok || (choiceStr != "High" && choiceStr != "Low") {
		return errors.New("choice must be 'High' or 'Low'")
	}

	// 検証のみ行い、エラーがなければnilを返す
	return nil
}

// 効果なしマス
type NoEffect struct {
}

func (e NoEffect) RequiresUserInput() bool { return false }

func (e NoEffect) GetOptions(tile *Tile) any { return nil }

func (e NoEffect) Apply(p *Player, g *Game, choice any) error {
	return nil
}

// ゴールマス
type GoalEffect struct {
}

// SetStatusEffect はプレイヤーのステータス（属性）を変更する効果です。
type SetStatusEffect struct {
	Status string `json:"status"` // 変更するステータス名 ("isMarried", "hasChildren", "job")
	Value  any    `json:"value"`  // 設定する値 (true, "professor" など)
}

func (e SetStatusEffect) RequiresUserInput() bool { return false }

func (e SetStatusEffect) GetOptions(tile *Tile) any { return nil }

func (e SetStatusEffect) Apply(p *Player, g *Game, choice any) error {
	switch e.Status {
	case "isMarried":
		if val, ok := e.Value.(bool); ok && val {
			p.marry()
		}
	case "children":
		if val, ok := e.Value.(float64); ok {
			p.changeChildren(int(val))
		}
	case "job":
		if val, ok := e.Value.(string); ok {
			p.setJob(val)
		}
	default:
		return fmt.Errorf("unknown status to set: %s", e.Status)
	}
	return nil
}

func (e GoalEffect) RequiresUserInput() bool { return false }

func (e GoalEffect) GetOptions(tile *Tile) any { return nil }

func (e GoalEffect) Apply(p *Player, g *Game, choice any) error {

	return nil
}

func CreateEffectFromJSON(data json.RawMessage) (Effect, error) {
	var ewt effectWithType
	if err := json.Unmarshal(data, &ewt); err != nil {
		return nil, fmt.Errorf("effect type unmarshal error: %w", err)
	}

	if ewt.Type == "" {
		return nil, errors.New("effect type is missing")
	}

	switch ewt.Type {
	case profit:
		var profitEffect ProfitEffect
		if err := json.Unmarshal(data, &profitEffect); err != nil {
			return nil, fmt.Errorf("ProfitEffect unmarshal error: %w", err)
		}
		return profitEffect, nil
	case loss:
		var lossEffect LossEffect
		if err := json.Unmarshal(data, &lossEffect); err != nil {
			return nil, fmt.Errorf("LossEffect unmarshal error: %w", err)
		}
		return lossEffect, nil
	case quiz:
		var quizEffect QuizEffect
		if err := json.Unmarshal(data, &quizEffect); err != nil {
			return nil, fmt.Errorf("QuizEffect unmarshal error: %w", err)
		}
		return quizEffect, nil
	case branch:
		var branchEffect BranchEffect
		if err := json.Unmarshal(data, &branchEffect); err != nil {
			return nil, fmt.Errorf("BranchEffect unmarshal error: %w", err)
		}
		return branchEffect, nil
	case overall:
		var overallEffect OverallEffect
		if err := json.Unmarshal(data, &overallEffect); err != nil {
			return nil, fmt.Errorf("OverallEffect unmarshal error: %w", err)
		}
		return overallEffect, nil
	case neighbor:
		var neighborEffect NeighborEffect
		if err := json.Unmarshal(data, &neighborEffect); err != nil {
			return nil, fmt.Errorf("NeighborEffect unmarshal error: %w", err)
		}
		return neighborEffect, nil
	case require:
		var requireEffect RequireEffect
		if err := json.Unmarshal(data, &requireEffect); err != nil {
			return nil, fmt.Errorf("RequireEffect unmarshal error: %w", err)
		}
		return requireEffect, nil
	case gamble:
		return GambleEffect{}, nil
	case goal:
		return GoalEffect{}, nil
	case conditional:
		var conditionalEffect ConditionalEffect
		if err := json.Unmarshal(data, &conditionalEffect); err != nil {
			return nil, fmt.Errorf("ConditionalEffect unmarshal error: %w", err)
		}
		return conditionalEffect, nil
	case setStatus:
		var setStatusEffect SetStatusEffect
		if err := json.Unmarshal(data, &setStatusEffect); err != nil {
			return nil, fmt.Errorf("SetStatusEffect unmarshal error: %w", err)
		}
		return setStatusEffect, nil
	default:
		return NoEffect{}, nil
	}
}

func InitQuiz() error {
	file, err := os.Open(QuizJSONPath)
	if err != nil {
		return fmt.Errorf("file open error: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&quizzes); err != nil {
		return fmt.Errorf("JSON decode error: %w", err)
	}
	return nil
}
