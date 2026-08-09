package resource

// Statistics represents the public transaction statistics returned to the client.
type Statistics struct {
	Target   Target   `json:"target"`
	Baseline Baseline `json:"baseline"`
}

// Target represents the transaction statistics for the requested range.
type Target struct {
	Overview   Overview          `json:"overview"`
	TopFive    []Transaction     `json:"topFive"`
	Categories []CategoryExpense `json:"categories"`
}

// Baseline represents the averaged comparison statistics.
type Baseline struct {
	Overview Overview `json:"overview"`
}

// Overview represents aggregated income and expense statistics.
type Overview struct {
	Expenses  float64 `json:"expenses"`
	Incomes   float64 `json:"incomes"`
	NetChange float64 `json:"netChange"`
}

// Transaction represents the public transaction fields returned inside the top-five list.
type Transaction struct {
	ID          string  `json:"id"`
	Money       float32 `json:"money"`
	Date        string  `json:"date"`
	Direction   int     `json:"direction"`
	Category    int     `json:"category"`
	Description string  `json:"description"`
}

// CategoryExpense represents aggregated expense statistics for a transaction category.
type CategoryExpense struct {
	Absolute float64 `json:"absolute"`
	Category int     `json:"category"`
	Percent  int     `json:"percent"`
}
