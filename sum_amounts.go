package sumamounts

type Obligation struct {
	ID        uint64
	AmountDue uint64
	Priority  uint64
}

type Line struct {
	ID      uint64
	Applied uint64
}

type Settlement struct {
	Lines     []Line
	Unapplied uint64
}

func Settle(payment uint64, obligations []Obligation) Settlement {
	panic("not implemented")
}
