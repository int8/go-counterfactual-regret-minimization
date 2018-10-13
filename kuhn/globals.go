package kuhn

import . "github.com/int8/gopoker"

var CheckAction = PlayerAction{Check}
var BetAction = PlayerAction{Bet}
var CallAction = PlayerAction{Call}
var FoldAction = PlayerAction{Fold}

const BetSize float64 = 1.0
const Ante = 1.0
