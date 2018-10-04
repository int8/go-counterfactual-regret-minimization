package gocfr


type TwoPlayersGameNode interface {
	IsTerminal() bool
	GetAvailableActions() []Action
	Play(Action) TwoPlayersGameNode
	NextToMove() Player
}

//TODO: Remember to model players and include them
type RhodeIslandGameState struct {
	round         Round
	deck          *FullDeck
	parent        *RhodeIslandGameState
	causingAction *Action
	availableActionsCache *ActionsCache
}

func (node *RhodeIslandGameState) NextToMove() Player {
	availableActions := node.GetAvailableActions()
	if availableActions == nil {
		return Environment
	} else {
		return availableActions[0].player
	}
}

func (node *RhodeIslandGameState) Play(action Action, table *PokerTable) RhodeIslandGameState {

	round := node.round
	if action.move == DealPrivateCards {
		// TODO: deal private cards here
		round = round.NextRound()
	}

	if action.move == DealPublicCard {
		round = round.NextRound()
		// TODO: deal public cards
	}

	child := RhodeIslandGameState{round, node.deck, node, &action, nil}
	return child
}

func (node *RhodeIslandGameState) computeAvailableActionsCache() {
	if node.availableActionsCache != nil {
		return
	}

	node.availableActionsCache = &ActionsCache{ nil}

	if node.round == Start {
		dealPrivateCards := Action{player: Environment, move: DealPrivateCards}
		node.availableActionsCache.actions = []Action{dealPrivateCards}
		return
	}

	// if any player folds - no further action - game ends
	if node.causingAction.move == Fold {
		node.availableActionsCache.actions = nil
		return
	}

	// whenever betting roung is over (CALL OR CHECK->CHECK) deal public cards or end if turn
	bettingRoundEnded := node.causingAction.move == Call  || (node.causingAction.move == Check && node.parent.causingAction.move == Check)
	if bettingRoundEnded {
		if node.round != Turn {
			dealPublicCard := Action{Environment, DealPublicCard}
			node.availableActionsCache.actions = []Action{dealPublicCard}
		} else {
			node.availableActionsCache.actions = nil
		}
		return
	}
	// single check implies BET or CHECK
	if node.causingAction.move == Check && node.parent.causingAction.move != Check {
		bet := Action{-node.causingAction.player, Bet}
		check := Action{-node.causingAction.player, Check}
		node.availableActionsCache.actions = []Action{bet, check}
		return
	}

	// you can only FOLD, RAISE or CALL on BET
	if node.causingAction.move == Bet {
		call := Action{-node.causingAction.player, Call}
		fold := Action{-node.causingAction.player, Fold}
		raise := Action{-node.causingAction.player, Raise}
		node.availableActionsCache.actions = []Action{call, fold, raise}
		return
	}

	// if RAISE, you can CALL FOLD or RAISE (unless there has been 6 prior raises - 3 for each player)
	if node.causingAction.move == Raise {
		if countPriorRaises(*node) < 6 {
			// allow raise if there has been less than 6 raises so far
			call := Action{-node.causingAction.player, Call}
			fold := Action{-node.causingAction.player, Fold}
			raise := Action{-node.causingAction.player, Raise}
			node.availableActionsCache.actions = []Action{call, fold, raise}
		} else {
			call := Action{-node.causingAction.player, Call}
			fold := Action{-node.causingAction.player, Fold}
			node.availableActionsCache.actions = []Action{call, fold}
		}
		return
	}

	if node.causingAction.move == DealPrivateCards || node.causingAction.move == DealPublicCard {
		bet := Action{PlayerA, Bet}
		check := Action{PlayerA, Check}
		node.availableActionsCache.actions = []Action{bet, check}
		return
	}

	node.availableActionsCache.actions = nil
}

func (node *RhodeIslandGameState) GetAvailableActions() []Action {
	node.computeAvailableActionsCache()
	return node.availableActionsCache.actions
}

func (node *RhodeIslandGameState) IsTerminal() bool {
	actions := node.GetAvailableActions()
	return actions == nil
}
