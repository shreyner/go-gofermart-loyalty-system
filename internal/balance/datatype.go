package balance

type Balance struct {
	Current   int
	Withdrawn int
}

type BalanceEntity struct {
	UserID    string
	Current   int
	Withdrawn int
	//CreatedAt time.Time
}
