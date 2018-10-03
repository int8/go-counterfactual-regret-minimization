package gocfr
import (
"testing"
)

func TestGameCreation(t *testing.T) {
	deck := CreateFullDeck()
	root := RhodeIslandGameState{PreFlop, &deck, nil, nil}
	if root.causingAction != nil {
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

	actions := root.GetAvailableActions()


	if actions == nil {
		t.Error("Game root should have one action available, no actions available")
	}

	if len(actions) != 1 {
		t.Errorf("Game root should have one action available, %v actions available", len(actions))
	}
}


func TestGamePlay(t *testing.T) {
	deck := CreateFullDeck()
	root := RhodeIslandGameState{Start, &deck, nil, nil}

	if root.round != Start {
		t.Errorf("Root should be in Start state, but it is in %v state", root.round)
	}

	actions := root.GetAvailableActions()
	preflopState := root.Play(actions[0])
	if preflopState.round != PreFlop {
		t.Errorf("Game should be in PreFlop state, but it is in %v state", preflopState.round)
	}

	actions = preflopState.GetAvailableActions()
	actionIndex := selectActionByMove(actions, Check)

	preflopStatePlayerAChecks := preflopState.Play(actions[actionIndex])
	if preflopStatePlayerAChecks.round != PreFlop {
		t.Errorf("Game should be in PreFlop state, but it is in %v state", preflopState.round)
	}


	actions = preflopStatePlayerAChecks.GetAvailableActions()
	actionIndex = selectActionByMove(actions, Check)
	preflopStatePlayerBChecks := preflopStatePlayerAChecks.Play(actions[actionIndex])

	if preflopStatePlayerBChecks.round != PreFlop {
		t.Errorf("Game should be in PreFlop state, but it is in %v state", preflopStatePlayerBChecks.round)
	}


	actions = preflopStatePlayerBChecks.GetAvailableActions()
	actionIndex = selectActionByMove(actions,DealPublicCard)
	flopState := preflopStatePlayerBChecks.Play(actions[actionIndex])

	if flopState.round != Flop {
		t.Errorf("Game should be in Flop state, but it is in %v state", flopState.round)
	}

	actions = flopState.GetAvailableActions()
	actionIndex = selectActionByMove(actions, Bet)
	flopStatePlayerABets := flopState.Play(actions[actionIndex])

	if flopStatePlayerABets.round != Flop {
		t.Errorf("Game should be in Flop state, but it is in %v state", flopStatePlayerABets.round)
	}

	actions = flopStatePlayerABets.GetAvailableActions()
	actionIndex = selectActionByMove(actions, Call)
	flopStatePlayerBCalls := flopStatePlayerABets.Play(actions[actionIndex])

	if flopStatePlayerBCalls.round != Flop {
		t.Errorf("Game should be in Flop state, but it is in %v state", flopStatePlayerBCalls.round)
	}

	actions = flopStatePlayerBCalls.GetAvailableActions()
	actionIndex = selectActionByMove(actions,  DealPublicCard)
	turnState := flopStatePlayerBCalls.Play(actions[actionIndex])

	if turnState.round != Turn {
		t.Errorf("Game should be in Turn state, but it is in %v state", turnState.round)
	}

}

