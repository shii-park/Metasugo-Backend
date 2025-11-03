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
	assert.Equal(t, "test_player", player.Id)

	// Test adding the same player again
	_, err = game.AddPlayer("test_player")
	assert.Error(t, err, "should return an error when adding a player with an existing ID")
}

func TestQuizEffect(t *testing.T) {
	game := NewGameWithTilesForTest("../../tiles.json")
	player, _ := game.AddPlayer("test_player")
	player.Money = 0 // Reset money for test

	quizTile := game.tileMap[9]
	assert.NotNil(t, quizTile)
	assert.Equal(t, TileKind("quiz"), quizTile.kind)

	player.Position = quizTile

	effect, ok := quizTile.Effect.(QuizEffect)
	assert.True(t, ok)

	// Test GetOptions
	quiz, ok := effect.GetOptions(quizTile).(Quiz)
	assert.True(t, ok)
	expectedOptions := []string{"1", "2", "3", "4"}
	assert.Equal(t, expectedOptions, quiz.Options)

	// Test Apply with correct answer
	initialMoney := player.Money
	err := effect.Apply(player, game, 1) // Correct answer index is 1
	assert.NoError(t, err)
	assert.Equal(t, initialMoney+10, player.Money)

	// Test Apply with incorrect answer (no penalty)
	initialMoney = player.Money
	err = effect.Apply(player, game, 2) // Incorrect answer
	assert.NoError(t, err)
	assert.Equal(t, initialMoney, player.Money)
}
