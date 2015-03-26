package action

const (
	None int = iota
	Add
	Update
	Del
	Wait
)

func ParseAction(t string) int {
	switch t {
	case "A":
		return Add
	case "U":
		return Update
	case "D":
		return Del
	case "W":
		return Wait
	default:
		return None
	}
}
