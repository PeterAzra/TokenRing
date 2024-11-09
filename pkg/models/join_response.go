package models

type JoinResponse struct {
	Ok          bool
	Left, Right string
}

func NewJoinResponse(left string, right string) *JoinResponse {
	return &JoinResponse{
		Ok:    true,
		Left:  left,
		Right: right,
	}
}
