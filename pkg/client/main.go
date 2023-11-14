package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"service/pkg/serv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	p          = big.NewInt(123)
	g          = big.NewInt(321)
	privateKey = big.NewInt(456)
	publicKey  = new(big.Int).Exp(g, privateKey, p)
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
	r2 := rand.Int63()
	challengeReq := &serv.AuthenticationChallengeRequest{
		User: "Alice",
		R2:   r2,
	}
	challengeRes, err := client.CreateAuthenticationChallenge(ctx, challengeReq)
	if err != nil {
		return err
	}

	// Authentication Answer
	s := new(big.Int).Exp(g, privateKey, p)
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
