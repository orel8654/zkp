package service

import (
	"crypto/rand"
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
	serverSecret = big.NewInt(42) // Секретное значение сервера
	p, _         = new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639747", 10)
	g, _         = new(big.Int).SetString("2", 10)
	h            = new(big.Int).Exp(g, serverSecret, p)
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
	r2, _ := rand.Int(rand.Reader, p)
	alpha := req.R1
	authID := fmt.Sprintf("%s_%d", req.User, time.Now().UnixNano())
	beta := new(big.Int).Add(r2, new(big.Int).Mul(big.NewInt(alpha), serverSecret))
	beta.Mod(beta, p)
	if err := m.Rb.SaveVal(authID+"r2", r2.Int64()); err != nil {
		log.Printf("redis save failed - %s", err) // debug
		return nil, errors.New("can't created challenge")
	}
	if err := m.Rb.SaveVal(authID+"r1", alpha); err != nil {
		log.Printf("redis save failed - %s", err) // debug
		return nil, errors.New("can't created challenge")
	}
	return &serv.AuthenticationChallengeResponse{
		AuthId: authID,
		C:      beta.Int64(),
	}, nil
}

func (m *MyAuthServer) VerifyAuthentication(ctx context.Context, req *serv.AuthenticationAnswerRequest) (*serv.AuthenticationAnswerResponse, error) {
	log.Print("VerifyAuthentication method")
	r2, err := m.Rb.GetVal(req.AuthId + "r2")
	if err != nil {
		log.Printf("redis get failed - %s", err) // debug
		return nil, errors.New("data not found")
	}
	r1, err := m.Rb.GetVal(req.AuthId + "r1")
	if err != nil {
		log.Printf("redis get failed - %s", err) // debug
		return nil, errors.New("data not found")
	}
	authID := req.AuthId
	sValue := new(big.Int).Exp(g, big.NewInt(serverSecret.Int64()), p)
	alpha := new(big.Int).Exp(g, big.NewInt(req.S), p)

	s_r := new(big.Int).Sub(alpha, new(big.Int).Mul(sValue, big.NewInt(r2)))
	s_r.Div(s_r, big.NewInt(r1))

	expectedResult := new(big.Int).Exp(g, s_r, p)
	expectedResult.Mul(expectedResult, new(big.Int).Exp(h, big.NewInt(r2), p))
	expectedResult.Mod(expectedResult, p)

	if expectedResult.Cmp(h) == 0 {
		sessionID := fmt.Sprintf("session_%s_%d", authID, time.Now().UnixNano())
		return &serv.AuthenticationAnswerResponse{
			SessionId: sessionID,
		}, nil
	}
	return nil, errors.New("auth failed")
}
