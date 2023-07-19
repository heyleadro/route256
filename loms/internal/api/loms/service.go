package loms

import (
	"route256/loms/internal/service/loms"
	"route256/loms/pkg/loms_v1"
)

type Service struct {
	loms_v1.UnimplementedLomsServer
	impl *loms.Service
}

func NewLomsServer(impl *loms.Service) loms_v1.LomsServer {
	return &Service{impl: impl}
}
