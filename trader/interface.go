package trader

import (
	"nofx/trader/types"
)

// Re-export types for backward compatibility
type (
	ClosedPnLRecord = types.ClosedPnLRecord
	TradeRecord     = types.TradeRecord
	Trader          = types.Trader
	OpenOrder       = types.OpenOrder
)
