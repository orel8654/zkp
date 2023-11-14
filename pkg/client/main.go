package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"service/pkg/serv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	x = new(big.Int).SetInt64(3)
	q = new(big.Int).SetInt64(11)
	g = new(big.Int).SetInt64(5)
	h = new(big.Int).SetInt64(7)
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
		Y1:   new(big.Int).Exp(g, x, nil).Int64(),
		Y2:   new(big.Int).Exp(g, x, nil).Int64(),
	}
	regRes, err := client.Register(ctx, regReq)
	_ = regRes
	if err != nil {
		return err
	}
	log.Print("Registration successful")

	// Authentication Challenge
	k := new(big.Int).SetInt64(2)
	challengeReq := &serv.AuthenticationChallengeRequest{
		User: "Alice",
		R1:   new(big.Int).Exp(g, x, nil).Int64(),
		R2:   new(big.Int).Exp(g, x, nil).Int64(),
	}
	challengeRes, err := client.CreateAuthenticationChallenge(ctx, challengeReq)
	if err != nil {
		return err
	}

	// Authentication Answer
	c := challengeRes.C
	s := new(big.Int).Sub(k, new(big.Int).Mul(big.NewInt(c), x))
	s.Mod(s, q)
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
