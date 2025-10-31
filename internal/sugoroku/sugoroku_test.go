package sugoroku

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CreateTestFileはtile_test.goに存在するため、ここでは削除しました。

func TestMain(m *testing.M) {
	// Setup
	originalQuizJSONPath := QuizJSONPath
	QuizJSONPath = "../../test/test_quizzes.json"
	err := InitQuiz()
	if err != nil {
		panic("failed to initialize quiz for test: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Teardown
	QuizJSONPath = originalQuizJSONPath
	os.Exit(code)
}

func TestGame_AddPlayer(t *testing.T) {
	game := NewGameWithTilesForTest("../../tiles.json")
	player, err := game.AddPlayer("test_player")

	assert.NoError(t, err)
	assert.NotNil(t, player)
	assert.Equal(t, "test_player", player.GetID())

	// Test adding the same player again
	_, err = game.AddPlayer("test_player")
	assert.Error(t, err, "should return an error when adding a player with an existing ID")
}

func TestQuizEffect(t *testing.T) {
	game := NewGameWithTilesForTest("../../tiles.json")
	player, _ := game.AddPlayer("test_player")
	player.money = 0 // Reset money for test

	quizTile := game.tileMap[9]
	assert.NotNil(t, quizTile)
	assert.Equal(t, TileKind("quiz"), quizTile.kind)

	player.position = quizTile

	effect, ok := quizTile.effect.(QuizEffect)
	assert.True(t, ok)

	// Test GetOptions
	quiz, ok := effect.GetOptions(quizTile).(Quiz)
	assert.True(t, ok)
	expectedOptions := []string{"1", "2", "3", "4"}
	assert.Equal(t, expectedOptions, quiz.Options)

	// Test Apply with correct answer
	player.money = 0 // Reset money for test
	initialMoney := player.GetMoney()
	err := effect.Apply(player, game, 1) // Correct answer index is 1
	assert.NoError(t, err)
	assert.Equal(t, initialMoney+effect.Amount, player.GetMoney())

	// Test Apply with incorrect answer
	player.money = 0 // Reset money for test
	initialMoney = player.GetMoney()
	err = effect.Apply(player, game, 2) // Incorrect answer
	assert.NoError(t, err)
	assert.Equal(t, initialMoney-effect.Amount, player.GetMoney())
}

func TestChildBonusEffect(t *testing.T) {
	game := NewGameWithTilesForTest("../../tiles.json")
	player, _ := game.AddPlayer("test_player")
	player.money = 0      // Reset money for test
	player.Children = 2 // Set number of children

	// Test with profit
	profitEffect := ChildBonusEffect{ProfitAmountPerChild: 100}
	initialMoney := player.GetMoney()
	err := profitEffect.Apply(player, game, nil)
	assert.NoError(t, err)
	assert.Equal(t, initialMoney+(2*100), player.GetMoney())

	// Test with loss
	player.money = 0 // Reset money for test
	lossEffect := ChildBonusEffect{LossAmountPerChild: 50}
	initialMoney = player.GetMoney()
	err = lossEffect.Apply(player, game, nil)
	assert.NoError(t, err)
	assert.Equal(t, initialMoney-(2*50), player.GetMoney())
}