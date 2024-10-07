package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/stonoy/my_remainder/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
}

type Remainder struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	HasPriority bool      `json:"has_priority"`
	Timing      time.Time `json:"timing"`
	Userid      uuid.UUID `json:"user_id"`
}

func dbToRespUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
	}
}

func dbToRespRemainder(remainder database.Remainder) Remainder {
	return Remainder{
		ID:          remainder.ID,
		CreatedAt:   remainder.CreatedAt,
		UpdatedAt:   remainder.UpdatedAt,
		Subject:     remainder.Subject,
		Description: remainder.Description,
		Timing:      remainder.Timing,
		HasPriority: remainder.HasPriority,
		Userid:      remainder.Userid,
	}
}

func dbToRespRemainders(remainders []database.Remainder) []Remainder {
	final := []Remainder{}

	for _, remainder := range remainders {
		final = append(final, dbToRespRemainder(remainder))
	}

	return final
}
