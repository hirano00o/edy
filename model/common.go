package model

type AttributeType int

const (
	S AttributeType = iota
	N
	B
)

func ConvertToAttributeType(s string) AttributeType {
	switch s {
	case "S":
		return S
	case "N":
		return N
	case "B":
		return B
	default:
		return S
	}
}

func ConvertAttributeTypeToString(aType AttributeType) string {
	switch aType {
	case S:
		return "S"
	case N:
		return "N"
	case B:
		return "B"
	default:
		return "S"
	}
}
