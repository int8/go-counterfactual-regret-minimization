## Counterfactual regret minimization in Go  


Implementation of chance sampling [counterfactual regret minimization](https://int8.io/counterfactual-regret-minimization-for-poker-ai/) in Go.

 #### Installation 
 Install with 
 ```bash
 go get github.com/int8/go-counterfactual-regret-minimization/cfr
 ```
 
 
 
#### Usage 

The central struct is ```cfr.ComputingRoutine``` which holds GameState interface and global Strategy maps (keeping data necessary to compute nash equilibrium at the end) 

To use it for your imperfect-information sum-zero strictly competitive two players game (like various kinds of poker) you need to provide implementation of ```GameState``` interface.


```go
// GameState - state of the game interface
type GameState interface {
	Parent() GameState
	Act(action acting.Action) GameState
	InformationSet() InformationSet
	Actions() []acting.Action
	IsTerminal() bool
	CurrentActor() acting.Actor
	Evaluate() float32
}
```

#### Rhode Island Poker example 
Example implementations of Rhode Island Poker and Kuhn Poker are included in repository.

```go 
import (	
	"github.com/int8/go-counterfactual-regret-minimization/acting"
	"github.com/int8/go-counterfactual-regret-minimization/cards"
	"github.com/int8/go-counterfactual-regret-minimization/games/rhodeisland"
	"github.com/int8/go-counterfactual-regret-minimization/cfr"
)

playerA := &rhodeisland.Player{Id: acting.PlayerA, Actions: nil, Card: nil, Stack: 1000.}
playerB := &rhodeisland.Player{Id: acting.PlayerB, Actions: nil, Card: nil, Stack: 1000.}
root := rhodeisland.Root(playerA, playerB, cards.CreateLimitedDeck(cards.C10, true))
routine := cfr.CreateComputingRoutine(root)
nashEquilibrium := routine.ComputeNashEquilibriumViaCFR(10000)
```