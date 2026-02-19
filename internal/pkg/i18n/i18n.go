package i18n

import (
	"context"
	"fmt"
	"strings"
)

type contextKey string

const langKey contextKey = "lang"

const (
	LangEN = "en"
	LangID = "id"
)

func Detect(acceptLanguage string) string {
	if acceptLanguage == "" {
		return LangEN
	}
	for _, part := range strings.Split(acceptLanguage, ",") {
		tag := strings.ToLower(strings.TrimSpace(strings.Split(part, ";")[0]))
		lang := strings.Split(tag, "-")[0]
		switch lang {
		case LangID:
			return LangID
		case LangEN:
			return LangEN
		}
	}
	return LangEN
}

func WithLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, langKey, lang)
}

func FromContext(ctx context.Context) string {
	lang, ok := ctx.Value(langKey).(string)
	if !ok || lang == "" {
		return LangEN
	}
	return lang
}

type OTPEmailContent struct {
	Subject string
	Body    string
}

func OTPEmail(name, otp, purpose, lang string) OTPEmailContent {
	if lang == LangID {
		return otpEmailID(name, otp, purpose)
	}
	return otpEmailEN(name, otp, purpose)
}

func otpEmailEN(name, otp, purpose string) OTPEmailContent {
	return OTPEmailContent{
		Subject: otpSubjectEN(purpose),
		Body: fmt.Sprintf(
			"Hi %s,\n\nYour verification code is:\n\n%s\n\nThis code will expire in 10 minutes.\n\nIf you did not request this, please ignore this email.\n\nRegards,\nHaily Team",
			name, otp,
		),
	}
}

func otpEmailID(name, otp, purpose string) OTPEmailContent {
	return OTPEmailContent{
		Subject: otpSubjectID(purpose),
		Body: fmt.Sprintf(
			"Halo %s,\n\nKode verifikasi kamu adalah:\n\n%s\n\nKode ini akan kedaluwarsa dalam 10 menit.\n\nJika kamu tidak merasa meminta kode ini, abaikan email ini.\n\nSalam,\nTim Haily",
			name, otp,
		),
	}
}

func otpSubjectEN(purpose string) string {
	switch purpose {
	case "EMAIL_VERIFICATION":
		return "Email Verification Code"
	case "PASSWORD_RESET":
		return "Password Reset Code"
	default:
		return "Verification Code"
	}
}

func otpSubjectID(purpose string) string {
	switch purpose {
	case "EMAIL_VERIFICATION":
		return "Kode Verifikasi Email"
	case "PASSWORD_RESET":
		return "Kode Reset Password"
	default:
		return "Kode Verifikasi"
	}
}

func RegisterSuccessMessage(lang string) string {
	if lang == LangID {
		return "Registrasi berhasil. Silakan cek email kamu untuk kode OTP verifikasi."
	}
	return "Registration successful. Please check your email for OTP verification."
}

func ResendOTPSuccessMessage(lang string) string {
	if lang == LangID {
		return "OTP berhasil dikirim. Silakan cek email kamu."
	}
	return "OTP sent successfully. Please check your email."
}
