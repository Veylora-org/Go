package sumamounts

import (
	"fmt"
	"strings"
	"testing"
)

func ob(id, amountDue, priority uint64) Obligation {
	return Obligation{ID: id, AmountDue: amountDue, Priority: priority}
}

func lines(pairs ...[2]uint64) []Line {
	out := make([]Line, len(pairs))
	for i, pair := range pairs {
		out[i] = Line{ID: pair[0], Applied: pair[1]}
	}
	return out
}

func equalSettlement(got, want Settlement) bool {
	if got.Unapplied != want.Unapplied || len(got.Lines) != len(want.Lines) {
		return false
	}
	for i := range want.Lines {
		if got.Lines[i] != want.Lines[i] {
			return false
		}
	}
	return true
}

func TestSettle(t *testing.T) {
	tests := []struct {
		name          string
		payment       uint64
		obligations   []Obligation
		wantLines     []Line
		wantUnapplied uint64
	}{
		{"empty", 100, nil, nil, 100},
		{"pays in full", 500, []Obligation{ob(1, 100, 1), ob(2, 40, 3), ob(3, 60, 2)}, lines([2]uint64{1, 100}, [2]uint64{2, 40}, [2]uint64{3, 60}), 300},
		{"higher priority first", 60, []Obligation{ob(1, 80, 1), ob(2, 50, 9), ob(3, 25, 9)}, lines([2]uint64{1, 0}, [2]uint64{2, 50}, [2]uint64{3, 10}), 0},
		{"same priority", 10, []Obligation{ob(1, 50, 1), ob(2, 30, 1), ob(3, 20, 1)}, lines([2]uint64{1, 10}, [2]uint64{2, 0}, [2]uint64{3, 0}), 0},
		{"spills to next priority", 120, []Obligation{ob(1, 100, 1), ob(2, 40, 5), ob(3, 60, 5)}, lines([2]uint64{1, 20}, [2]uint64{2, 40}, [2]uint64{3, 60}), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Settle(tt.payment, tt.obligations)
			want := Settlement{Lines: tt.wantLines, Unapplied: tt.wantUnapplied}
			if !equalSettlement(got, want) {
				t.Fatalf("got %+v\nwant %+v", got, want)
			}
		})
	}
}

func TestDuplicateID(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected a panic")
		}
		if !strings.Contains(strings.ToLower(fmt.Sprint(recovered)), "duplicate") {
			t.Fatalf("panic %v", recovered)
		}
	}()
	Settle(5, []Obligation{ob(1, 10, 1), ob(1, 4, 2)})
}
