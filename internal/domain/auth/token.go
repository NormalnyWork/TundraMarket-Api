package auth

const (
	RoleNomad          = "nomad"
	RoleTradingStation = "trading_station"
)

type TokenClaims struct {
	Role             string
	Phone            string
	NomadID          *int32
	TradingStationID *int32
}

type TokenIssuer interface {
	Issue(claims TokenClaims) (string, error)
}
