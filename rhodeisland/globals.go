package rhodeisland

import . "github.com/int8/gopoker"

var CheckAction = PlayerAction{Check}
var BetAction = PlayerAction{Bet}
var CallAction = PlayerAction{Call}
var RaiseAction = PlayerAction{Raise}
var FoldAction = PlayerAction{Fold}

const PreFlopBetSize = 10.
const PostFlopBetSize = 20.

const MaxRaises = 3
const Ante = 5.0
