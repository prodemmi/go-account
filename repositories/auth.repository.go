package repositories

type (
	AuthRepository  struct{}
	IAuthRepository interface{}
)

func NewAuthRepository() IAuthRepository {
	return AuthRepository{}
}
