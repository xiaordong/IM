package e

const (
	SUCCESS               = 200
	UpdatePasswordSuccess = 201
	NotExistInentifier    = 202
	ERROR                 = 500
	InvalidParams         = 400
	ErrorDatabase         = 40001

	ErrorHeaderData = 30001
	ErrorValidToken = 30002
	ErrorGetFile    = 30003
	ErrorCreateFile = 30004
	ErrorPlaceFile  = 30005
	DeleteError     = 30006
	ErrorCreateData = 30007
	ErrorGroupFull  = 30008
	ErrorNoData     = 30009
	ExistUser       = 30010

	WebsocketSuccessMessage = 50001
	WebsocketSuccess        = 50002
	WebsocketEnd            = 50003
	WebsocketOnlineReply    = 50004
	WebsocketOfflineReply   = 50005
	WebsocketLimit          = 50006
)
