package constant

var CLOTH_TYPE = map[string]string{
	"kemeja": "Kemeja",
	"outer":  "Outer",
	"kain":   "Kain",
}

type SizeVariant struct {
	Size  string
	Price float64
}

var CLOTH_SIZE = []SizeVariant{
	{Size: "S", Price: 120000},
	{Size: "M", Price: 120000},
	{Size: "L", Price: 130000},
	{Size: "XL", Price: 150000},
	{Size: "XXL", Price: 160000},
}

const (
	OrderPending  = "pending"
	OrderProcess  = "proses"
	OrderShipped  = "dikirim"
	OrderDone     = "selesai"
	OrderCanceled = "dibatalkan"
)

const (
	TransactionPending  = "pending"
	TransactionProcess  = "proses"
	TransactionDone     = "selesai"
	TransactionCanceled = "dibatalkan"
)
