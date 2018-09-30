package gocfr
import (
	"testing"
)

func TestGameCreation(t *testing.T) {
	root := createRhodeIslandGameRoot()
	if root.CausingAction != nil {
		t.Error("Root node should not have causing action")
	}

	if root.parent != nil {
		t.Error("Root node should not have nil parent")
	}

	if root.round != PreFlop {
		t.Error("First round of the game should be PreFlop")
	}

	if root.IsTerminal() == true {
		t.Error("Game root should not be terminal")
	}

	actions, err := root.GetAvailableActions()
	if err != nil {
		t.Error(err)
	}

	if actions == nil {
		t.Error("Game root should have one action available, no actions available")
	}

	if len(actions) != 1 {
		t.Errorf("Game root should have one action available, %v actions available", len(actions))
	}

}

