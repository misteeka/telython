package status

type Status int

var (
	YES Status = 1
	NO  Status = 0
)

var (
	SUCCESS                Status = 100
	INVALID_REQUEST        Status = 101
	INTERNAL_SERVER_ERROR  Status = 102
	AUTHORIZATION_FAILED   Status = 103
	INVALID_CURRENCY_CODE  Status = 104
	CURRENCY_CODE_MISMATCH Status = 105
	NOT_FOUND              Status = 106
	WRONG_AMOUNT           Status = 107
	INSUFFICIENT_FUNDS     Status = 108
	TOO_MANY_REQUESTS      Status = 109
)
