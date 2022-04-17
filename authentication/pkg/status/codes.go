package status

import "strconv"

type Status int

func (status Status) Equals(anotherStatus Status) bool {
	return status == anotherStatus
}
func (status Status) Serialize() string {
	return strconv.Itoa(int(status))
}

var (
	YES Status = 1
	NO  Status = 0
)

var (
	SUCCESS               Status = 100
	INVALID_REQUEST       Status = 101
	INTERNAL_SERVER_ERROR Status = 102
	AUTHORIZATION_FAILED  Status = 103
	ALREADY_EXISTS        Status = 104
	NOT_FOUND             Status = 105
)
