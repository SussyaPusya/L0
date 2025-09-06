package test

import (
	"github.com/SussyaPusya/L0/internal/dto"
	"github.com/brianvoe/gofakeit/v7"

	"strings"
)

func GenerateOrder() *dto.Order {
	return &dto.Order{
		OrderUID:    strings.ReplaceAll(gofakeit.UUID(), "-", ""),
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: dto.Delivery{
			Name:    gofakeit.Name(),
			Phone:   "+1" + gofakeit.Numerify("##########"),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Address().Street,
			Region:  gofakeit.RandomString([]string{"North", "South", "West", "East"}),
			Email:   gofakeit.Email(),
		},
		Payment: dto.Payment{
			Transaction:  strings.ReplaceAll(gofakeit.UUID(), "-", ""),
			RequestID:    gofakeit.Numerify("####"),
			Currency:     gofakeit.CurrencyShort(),
			Provider:     gofakeit.RandomString([]string{"wbpay", "sberpay", "alipay"}),
			Amount:       gofakeit.Float64(),
			PaymentDt:    int64(gofakeit.Number(1231231233, 1637907727)),
			Bank:         gofakeit.RandomString([]string{"alpha", "sber", "tbank", "pspb"}),
			DeliveryCost: float64(gofakeit.Number(0, 2000)),
			GoodsTotal:   gofakeit.Number(0, 500),
			CustomFee:    gofakeit.Number(0, 10),
		},
		Items: []dto.Item{
			{
				ChrtID:      gofakeit.Number(100000, 999999),
				TrackNumber: "WBILMTESTTRACK",
				Price:       gofakeit.Float64(),
				Rid:         strings.ReplaceAll(gofakeit.UUID(), "-", ""),
				Name:        gofakeit.ProductName(),
				Sale:        gofakeit.Number(0, 50),
				Size:        gofakeit.Numerify("#"),
				TotalPrice:  float64(gofakeit.Number(0, 5000)),
				NmID:        gofakeit.Number(100000, 999999),
				Brand:       gofakeit.Company(),
				Status:      gofakeit.HTTPStatusCode(),
			},
		},
		Locale:            gofakeit.LanguageAbbreviation(),
		InternalSignature: gofakeit.RandomString([]string{"warehouse", "", "pvz"}),
		CustomerID:        strings.ReplaceAll(gofakeit.UUID(), "-", ""),
		DeliveryService:   gofakeit.RandomString([]string{"wb", "ali", "ozon"}),
		ShardKey:          gofakeit.Numerify("##"),
		SmID:              gofakeit.Number(0, 100),
		DateCreated:       gofakeit.Date(),
		OofShard:          gofakeit.Numerify("#"),
	}
}
