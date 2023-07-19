package integration_test

import (
	"context"
	"log"
	"route256/checkout/internal/model"
	"route256/checkout/internal/repository/postgres"
	"route256/checkout/internal/repository/postgres/tx"
	"route256/checkout/internal/service/checkout"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	suite.Suite
	checkout checkout.CartRepository
}

var (
	URL = "postgres://user:password@localhost:5432/checkout?sslmode=disable"
)

func (s *Suite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool, err := pgxpool.Connect(ctx, URL)
	if err != nil {
		log.Fatalln("connect to db: %w", err)
	}

	provider := tx.New(pool)
	repo := postgres.NewRepo(provider)
	s.checkout = repo
}

func (s *Suite) TestCheckoutDB() {
	//todo: clean db

	const (
		user   = int64(1)
		sku1   = uint32(1)
		count1 = uint16(1)
		sku2   = uint32(2)
		count2 = uint16(2)
	)
	srcItems := []model.Cart{
		{
			SKU:   sku1,
			Count: count1 * 2,
		},
		{
			SKU:   sku2,
			Count: count2,
		},
	}

	//Initial adding to cart
	err := s.checkout.AddToCart(context.Background(), user, sku1, count1)
	s.Require().NoError(err)

	//Secondary adding to cart with same sku (checking for double count in future)
	err = s.checkout.AddToCart(context.Background(), user, sku1, count1)
	s.Require().NoError(err)

	//Secondary adding to cart with different sku
	err = s.checkout.AddToCart(context.Background(), user, sku2, count2)
	s.Require().NoError(err)

	//Getting Items and checking for element matching with srcItems
	dstItems, err := s.checkout.GetUserItems(context.Background(), user)
	s.Require().NoError(err)
	s.Require().ElementsMatch(dstItems, srcItems)

	//Checking the deletion from cart
	countDel := count2 - 1
	err = s.checkout.DeleteFromCart(context.Background(), user, sku2, countDel)
	s.Require().NoError(err)

	//Updating srcItems for deletion checking
	srcItems[1].Count = count2 - countDel
	dstItems, err = s.checkout.GetUserItems(context.Background(), user)
	s.Require().NoError(err)
	s.Require().ElementsMatch(dstItems, srcItems)

	//Checking item Count by Sku request
	dstCount, err := s.checkout.GetUserSkuItemCount(context.Background(), user, sku1)
	s.Require().NoError(err)
	s.Require().Equal(dstCount, srcItems[0].Count)

	dstCount, err = s.checkout.GetUserSkuItemCount(context.Background(), user, sku2)
	s.Require().NoError(err)
	s.Require().Equal(dstCount, srcItems[1].Count)

	//todo cleaning
}
