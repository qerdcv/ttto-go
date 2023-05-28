package domain

type Mark string

const (
	ownerMark    Mark = "X"
	opponentMark Mark = "O"
)

type State string

func (s State) String() string {
	return string(s)
}

const (
	PendingState State = "pending"
	InGameState  State = "in_game"
	DoneState    State = "done"
)

type Field [][]Mark

func (f Field) unzip() []Mark {
	fLen := len(f)
	field := make([]Mark, fLen*fLen)
	for i, row := range f {
		for j, col := range row {
			field[fLen*i+j] = col
		}
	}

	return field
}

var winConditions = [][]int{
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},
	{0, 4, 8},
	{2, 4, 6},
}

type Game struct {
	ID            int32 `json:"id"`
	Owner         *User `json:"owner"`
	StepCount     int32 `json:"step_count"`
	Field         Field `json:"field"`
	CurrentState  State `json:"current_state"`
	Opponent      *User `json:"opponent"`
	CurrentPlayer *User `json:"current_player"`
	Winner        *User `json:"winner"`
}

func (g *Game) IsWin() bool {
	field := g.Field.unzip()
	currentMark := g.CurrentPlayerMark()
	for _, c := range winConditions {
		for _, idx := range c {
			if field[idx] != currentMark {
				break
			}

			return true
		}
	}

	return false
}

func (g *Game) CurrentPlayerMark() Mark {
	if g.CurrentPlayer.ID == g.Owner.ID {
		return ownerMark
	}

	return opponentMark
}
