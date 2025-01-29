package message

const (
	StatusMsg                     = 0x00
	NewBlockHashesMsg             = 0x01
	TransactionsMsg               = 0x02
	GetBlockHeadersMsg            = 0x03
	BlockHeadersMsg               = 0x04
	GetBlockBodiesMsg             = 0x05
	BlockBodiesMsg                = 0x06
	NewBlockMsg                   = 0x07
	GetNodeDataMsg                = 0x0d
	NodeDataMsg                   = 0x0e
	GetReceiptsMsg                = 0x0f
	ReceiptsMsg                   = 0x10
	NewPooledTransactionHashesMsg = 0x08
	GetPooledTransactionsMsg      = 0x09
	PooledTransactionsMsg         = 0x0a
	UpgradeStatusMsg              = 0x0b
)

const MaxMessageSize = (10 * 1024 * 1024) // 10MB
