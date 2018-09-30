package gocfr

import (
	"fmt"
)

type Player int8
type Move int8
type Round int8


type Action struct {
	player Player
	move Move
}

type TwoPlayerGameNode interface {
	IsTerminal() (bool, error)
	GetAvailableActions() ([]Action, error)
	Play(Action) TwoPlayerGameNode
}

type RhodeIslandGameNode struct {
	round Round
	deck FullDeck
	parent *RhodeIslandGameNode
	CausingAction *Action
}


func createRhodeIslandGameRoot() RhodeIslandGameNode {
	deck := CreateFullDeck()
	root := RhodeIslandGameNode{PreFlop, deck, nil, nil}
	return root
}

func (node RhodeIslandGameNode) GetAvailableActions() ([]Action, error){

	if node.parent == nil {
		dealPrivateCards := Action{player: Chance, move: DealPrivateCards}
		return []Action{dealPrivateCards}, nil
	}

	roundEnded := (node.CausingAction.move == Call || node.CausingAction.move == Fold)
	roundEnded = roundEnded && (node.CausingAction.move == Check && node.parent.CausingAction.player != Chance)
 	if roundEnded {
 		if node.round != Turn {
 			dealPublicCard := Action{Chance, DealPublicCard}
			return []Action{dealPublicCard}, nil
		} else {
			return nil, nil
		}
	}

	if node.CausingAction.move == Check && node.parent.CausingAction.player == Chance {
		bet := Action{-node.CausingAction.player, Bet}
		check := Action{-node.CausingAction.player, Check}
		return []Action{bet, check}, nil
	}

	if node.CausingAction.move == Bet {
		call := Action{-node.CausingAction.player, Call}
		fold := Action{-node.CausingAction.player, Fold}
		raise := Action{-node.CausingAction.player, Raise}
		return []Action{call, fold, raise}, nil
	}

	if node.CausingAction.move == Raise {
		previousRaises := func(node RhodeIslandGameNode) int{
			count := 0
			for (node.CausingAction.move != Raise) {
				node = *node.parent
				count++
			}
			return count
		}(node)

		if previousRaises < 6 {
			call := Action{-node.CausingAction.player, Call}
			fold := Action{-node.CausingAction.player, Fold}
			raise := Action{-node.CausingAction.player, Raise}
			return []Action{call, fold, raise}, nil
		} else {
			call := Action{-node.CausingAction.player, Call}
			fold := Action{-node.CausingAction.player, Fold}
			return []Action{call, fold}, nil
		}
	}

	if node.CausingAction.move == DealPrivateCards {
		bet := Action{PlayerA, Bet}
		check := Action{PlayerA, Check}
		return []Action{bet, check}, nil
	}

	return nil, fmt.Errorf("No available actions computed, something is very wrong :)")
}


func (node RhodeIslandGameNode) IsTerminal() bool {
	actions, err := node.GetAvailableActions()
	return (actions == nil && err == nil)
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


