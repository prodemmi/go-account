package outbound

// OTPGenerator is the outbound port implemented by the otpgen adapter.
// Length is clamped by the caller (see app.Normalize) and charset is one of
// numeric, alpha or alphanumeric.
type OTPGenerator interface {
	GenerateOTP(length int, charset string) (string, error)
}
