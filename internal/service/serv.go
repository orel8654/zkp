package service

import (
	"errors"
	"fmt"
	"log"
	"service/internal/logic"
	"service/pkg/serv"
	"time"

	red "service/internal/db/redis"

	"golang.org/x/net/context"
)

type MyAuthServer struct {
	serv.UnimplementedAuthServer
	Rb *red.MyRedis
}

func (m *MyAuthServer) Register(ctx context.Context, req *serv.RegisterRequest) (*serv.RegisterResponse, error) {
	log.Printf("User %s registered with y1: %d, y2: %d\n", req.User, req.Y1, req.Y2)
	if err := m.Rb.SaveVal(req.User, req.Y2); err != nil {
		log.Printf("redis save failed - %s", err) // debug
	}
	return &serv.RegisterResponse{}, nil
}

func (m *MyAuthServer) CreateAuthenticationChallenge(ctx context.Context, req *serv.AuthenticationChallengeRequest) (*serv.AuthenticationChallengeResponse, error) {
	log.Print("CreateAuthenticationChallenge method")
	c, authID := logic.CreateChalleng(req.R1, req.R2, req.User)
	if err := m.Rb.SaveVal(authID, req.R2); err != nil {
		log.Printf("redis save failed - %s", err) // debug
		return nil, errors.New("can't created challenge")
	}
	return &serv.AuthenticationChallengeResponse{
		AuthId: authID,
		C:      c,
	}, nil
}

func (m *MyAuthServer) VerifyAuthentication(ctx context.Context, req *serv.AuthenticationAnswerRequest) (*serv.AuthenticationAnswerResponse, error) {
	log.Print("VerifyAuthentication method")
	challengVal, err := m.Rb.GetVal(req.AuthId)
	if err != nil {
		log.Printf("redis get failed - %s", err) // debug
		return nil, errors.New("data not found")
	}
	r := logic.CreateVerify(req.S, challengVal)
	if r {
		authID := fmt.Sprintf("session_%s_%d", req.AuthId, time.Now().UnixNano())
		return &serv.AuthenticationAnswerResponse{
			SessionId: authID,
		}, nil
	}
	return nil, errors.New("auth failed")
}
