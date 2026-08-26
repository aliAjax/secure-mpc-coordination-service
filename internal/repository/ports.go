package repository

import (
	"context"
	"github.com/example/027-mpc-coordinator/internal/domain"
)

type Store interface {
	CreateComputation(context.Context, *domain.Computation) error
	GetComputation(context.Context, string) (*domain.Computation, error)
	UpdateComputation(context.Context, *domain.Computation) error
	ListComputations(context.Context) ([]domain.Computation, error)
	PutParticipant(context.Context, *domain.Participant) error
	ListParticipants(context.Context, string) ([]domain.Participant, error)
	CreateRound(context.Context, *domain.Round) error
	GetRound(context.Context, string) (*domain.Round, error)
	UpdateRound(context.Context, *domain.Round) error
	PutEvidence(context.Context, *domain.Evidence) error
	ListEvidence(context.Context, string) ([]domain.Evidence, error)
}
