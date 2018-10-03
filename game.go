package gocfr

import (
	"fmt"
)

type Player int8
type Move int8
type Round int8


func (round Round) NextRound() Round {
	switch round {
	case Start:
		return PreFlop
	case PreFlop:
		return Flop
	case Flop:
		return Turn
	}
	return End
}

type Action struct {
	player Player
	move Move
}

type TwoPlayersGameNode interface {
	IsTerminal() bool
	GetAvailableActions() []Action
	Play(Action) TwoPlayersGameNode
}

//TODO: Remember to model players and include them
type RhodeIslandGameState struct {
	round Round
	deck *FullDeck
	parent *RhodeIslandGameState
	causingAction *Action
}


func (node *RhodeIslandGameState) Play(action Action) RhodeIslandGameState {
	round := node.round
	if action.move == DealPrivateCards {
		// TODO: deal private cards here
		round = round.NextRound()
	}

	if action.move == DealPublicCard {
		round = round.NextRound()
		// TODO: deal public cards
	}

	child := RhodeIslandGameState{round, node.deck,node, &action }
	return child
}

func (node *RhodeIslandGameState) GetAvailableActions() []Action {

	// actions for root of the game tree - first action is a chance action: dealing private cards
	if node.round == Start {
		dealPrivateCards := Action{player: Chance, move: DealPrivateCards}
		return []Action{dealPrivateCards}
	}

	// whenever betting roung is over (CALL OR FOLD OR CHECK->CHECK) deal public cards or end
	bettingRoundEnded := node.causingAction.move == Call || node.causingAction.move == Fold
	bettingRoundEnded = bettingRoundEnded || (node.causingAction.move == Check && node.parent.causingAction.move == Check)
 	if bettingRoundEnded {
 		if node.round != Turn {
 			dealPublicCard := Action{Chance, DealPublicCard}
			return []Action{dealPublicCard}
		} else {
			return nil
		}
	}
	// single check implies BET or CHECK
	if node.causingAction.move == Check && node.parent.causingAction.move != Check {
		bet := Action{-node.causingAction.player, Bet}
		check := Action{-node.causingAction.player, Check}
		return []Action{bet, check}
	}

	// you can only FOLD, RAISE or CALL on BET
	if node.causingAction.move == Bet {
		call := Action{-node.causingAction.player, Call}
		fold := Action{-node.causingAction.player, Fold}
		raise := Action{-node.causingAction.player, Raise}
		return []Action{call, fold, raise}
	}

	// if RAISE, you can CALL FOLD or RAISE (unless there has been 6 prior raises - 3 for each player)
	if node.causingAction.move == Raise {
		if countPriorRaises(*node) < 6 {
			// allow raise if there has been less than 6 raises so far
			call := Action{-node.causingAction.player, Call}
			fold := Action{-node.causingAction.player, Fold}
			raise := Action{-node.causingAction.player, Raise}
			return []Action{call, fold, raise}
		} else {
			call := Action{-node.causingAction.player, Call}
			fold := Action{-node.causingAction.player, Fold}
			return []Action{call, fold}
		}
	}

	if node.causingAction.move == DealPrivateCards || node.causingAction.move == DealPublicCard {
		bet := Action{PlayerA, Bet}
		check := Action{PlayerA, Check}
		return []Action{bet, check}
	}

	return nil
}


func (node *RhodeIslandGameState) IsTerminal() bool {
	actions := node.GetAvailableActions()
	return actions == nil
}


func (a Action) String() string {
	return fmt.Sprintf("%v:%v", a.player, a.move)
}

func (p Player) String() string {
	if p == PlayerA {
		return "A"
	}
	if p == PlayerB {
		return "B"
	}
	return "Chance"
}

func (m Move) String() string {
	switch m {
	case Check:
		return "Check"
	case Bet:
		return "Bet"
	case Call:
		return "Call"
	case Fold:
		return "Fold"
	case Raise:
		return "Raise"
	case DealPrivateCards:
		return "DealPrivateCards"
	case DealPublicCard:
		return "DealPublicCard"
	}
	return "Undefined"
}

func (r Round) String() string {
	switch r {
	case Start:
		return "Start"
	case End:
		return "End"
	case PreFlop:
		return "Preflop"
	case Turn:
		return "Turn"
	case Flop:
		return "Flop"
	}
	return "(?)"
}


