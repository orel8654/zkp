package main

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	r "math/rand"
	"service/pkg/serv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// p          = big.NewInt(123)
	// g          = big.NewInt(321)
	privateKey = big.NewInt(456)
	publicKey  = new(big.Int).Exp(g, privateKey, p)
	x          = "MySecretPassword123"
)

var (
	clientPrivateKey = big.NewInt(23) // Секретное значение клиента
	p, _             = new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639747", 10)
	g, _             = new(big.Int).SetString("2", 10)
)

func main() {
	if err := run(context.Background(), "localhost:8001"); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, address string) error {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := serv.NewAuthClient(conn)

	// Registration
	regReq := &serv.RegisterRequest{
		User: "Alice",
		Y1:   publicKey.Int64(),
		Y2:   new(big.Int).Exp(g, privateKey, p).Int64(),
	}
	regRes, err := client.Register(ctx, regReq)
	_ = regRes
	if err != nil {
		return err
	}
	log.Print("Registration successful")

	// Authentication Challenge
	alpha := r.Int63()
	challengeReq := &serv.AuthenticationChallengeRequest{
		User: "Alice",
		R1:   alpha,
	}
	challengeRes, err := client.CreateAuthenticationChallenge(ctx, challengeReq)
	if err != nil {
		return err
	}

	// Authentication Answer
	beta := challengeRes.C
	alpha = challengeReq.R1
	hashed := sha256.Sum256([]byte(x))
	int64Value := int64(binary.BigEndian.Uint64(hashed[:8]))
	s := new(big.Int).Sub(big.NewInt(alpha), new(big.Int).Mul(big.NewInt(int64Value), big.NewInt(beta)))
	s.Div(s, big.NewInt(alpha))

	authReq := &serv.AuthenticationAnswerRequest{
		AuthId: challengeRes.AuthId,
		S:      s.Int64(),
	}
	authRes, err := client.VerifyAuthentication(ctx, authReq)
	if err != nil {
		return err
	}

	fmt.Printf("Authentication successful. Session ID: %s\n", authRes.SessionId)
	return nil
}
