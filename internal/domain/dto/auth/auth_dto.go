package auth

import userDTO "github.com/haily-id/engine/internal/domain/dto/user"

type LoginResponse struct {
	Token string          `json:"token"`
	User  userDTO.UserDTO `json:"user"`
}

type RegisterResponse struct {
	User    userDTO.UserDTO `json:"user"`
	Message string          `json:"message"`
}
