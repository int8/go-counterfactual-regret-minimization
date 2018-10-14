package rhodeisland

import . "github.com/int8/gopoker"

var CheckAction = PlayerAction{Check}
var BetAction = PlayerAction{Bet}
var CallAction = PlayerAction{Call}
var RaiseAction = PlayerAction{Raise}
var FoldAction = PlayerAction{Fold}

const PreFlopBetSize float32 = 10.
const PostFlopBetSize float32 = 20.

const MaxRaises = 0
const Ante float32 = 5.0
