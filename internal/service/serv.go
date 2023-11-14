package service

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"service/pkg/serv"
	"time"

	red "service/internal/db/redis"

	"golang.org/x/net/context"
)

var (
	g = new(big.Int).SetInt64(5)
	h = new(big.Int).SetInt64(7)
	q = new(big.Int).SetInt64(11)
)

type MyAuthServer struct {
	serv.UnimplementedAuthServer
	Rb *red.MyRedis
}

func (m *MyAuthServer) Register(ctx context.Context, req *serv.RegisterRequest) (*serv.RegisterResponse, error) {
	log.Printf("User %s registered with y1: %d, y2: %d\n", req.User, req.Y1, req.Y2)
	if err := m.Rb.SaveVal(req.User+"y", req.Y1, req.Y2); err != nil {
		log.Printf("redis save failed - %s", err)
		return nil, err
	}
	return &serv.RegisterResponse{}, nil
}

func (m *MyAuthServer) CreateAuthenticationChallenge(ctx context.Context, req *serv.AuthenticationChallengeRequest) (*serv.AuthenticationChallengeResponse, error) {
	log.Print("CreateAuthenticationChallenge method")
	c := new(big.Int).SetInt64(4)
	if err := m.Rb.SaveVal(req.User+"r", req.R1, req.R2); err != nil {
		log.Printf("redis save failed - %s", err) // debug
		return nil, errors.New("can't created challenge")
	}
	if err := m.Rb.SaveVal(req.User+"c", c.Int64()); err != nil {
		log.Printf("redis save failed - %s", err) // debug
		return nil, errors.New("can't created challenge")
	}
	return &serv.AuthenticationChallengeResponse{
		AuthId: req.User,
		C:      c.Int64(),
	}, nil
}

func (m *MyAuthServer) VerifyAuthentication(ctx context.Context, req *serv.AuthenticationAnswerRequest) (*serv.AuthenticationAnswerResponse, error) {
	log.Print("VerifyAuthentication method")
	authUser := req.AuthId
	consts, err := m.Rb.GetVal(authUser+"r1", authUser+"r2", authUser+"y1", authUser+"y2", authUser+"c1")
	if err != nil {
		return nil, errors.New("auth failed")
	}

	checkR1 := new(big.Int).Exp(g, big.NewInt(req.S), nil)
	checkR1.Mul(checkR1, new(big.Int).Exp(big.NewInt(consts[authUser+"y1"]), big.NewInt(consts[authUser+"c1"]), nil))
	checkR1.Mod(checkR1, q)

	checkR2 := new(big.Int).Exp(h, big.NewInt(req.S), nil)
	checkR2.Mul(checkR2, new(big.Int).Exp(big.NewInt(consts[authUser+"y1"]), new(big.Int).SetInt64(3), nil))
	checkR2.Mod(checkR2, q)

	if big.NewInt(consts[authUser+"r1"]).Cmp(checkR1) == 0 && big.NewInt(consts[authUser+"r2"]).Cmp(checkR2) == 0 {
		sessionID := fmt.Sprintf("session_%s_%d", authUser, time.Now().UnixNano())
		return &serv.AuthenticationAnswerResponse{
			SessionId: sessionID,
		}, nil
	}
	return nil, errors.New("auth failed")
}
