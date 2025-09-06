package repository

import (
	"context"
	"fmt"

	"github.com/SussyaPusya/L0/internal/dto"
	"github.com/SussyaPusya/L0/pkg/redis"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	cache *redis.RedisChache
	pg    *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool, c *redis.RedisChache) *repository {

	return &repository{pg: pg, cache: c}
}

func (r *repository) CreateOrder(ctx context.Context, order *dto.Order) error {

	go r.cache.Set(ctx, order)
	tx, err := r.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)
	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid = $1)", order.OrderUID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check order existence: %w", err)
	}
	if exists {
		return err
	}
	_, err = tx.Exec(ctx, `
        INSERT INTO "orders" (
            order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO deliveries (
            order_uid, name, phone, zip, city, address, region, email
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO payments (
            order_uid, transaction, request_id, currency, provider,
            amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	for _, item := range order.Items {
		_, err = tx.Exec(ctx, `
            INSERT INTO items (
                order_uid, chrt_id, track_number, price, rid, name,
                sale, size, total_price, nm_id, brand, status
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        `,
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to insert item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *repository) GetOrder(ctx context.Context, orderID string) (*dto.Order, error) {

	order, err := r.cache.Get(ctx, orderID)
	if err != nil {
		var order dto.Order

		err := r.pg.QueryRow(ctx, `
        SELECT order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        FROM orders WHERE order_uid = $1
    `, orderID).Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get order: %w", err)
		}

		err = r.pg.QueryRow(ctx, `
        SELECT name, phone, zip, city, address, region, email
        FROM deliveries WHERE order_uid = $1
    `, orderID).Scan(
			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get delivery: %w", err)
		}

		err = r.pg.QueryRow(ctx, `
        SELECT transaction, request_id, currency, provider,
            amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
        FROM payments WHERE order_uid = $1
    `, orderID).Scan(
			&order.Payment.Transaction,
			&order.Payment.RequestID,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDt,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get payment: %w", err)
		}

		rows, err := r.pg.Query(ctx, `
        SELECT chrt_id, track_number, price, rid, name,
            sale, size, total_price, nm_id, brand, status
        FROM items WHERE order_uid = $1
    `, orderID)
		if err != nil {
			return nil, fmt.Errorf("failed to get items: %w", err)
		}
		defer rows.Close()

		var items []dto.Item
		for rows.Next() {
			var item dto.Item
			err := rows.Scan(
				&item.ChrtID,
				&item.TrackNumber,
				&item.Price,
				&item.Rid,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.TotalPrice,
				&item.NmID,
				&item.Brand,
				&item.Status,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan item: %w", err)
			}
			items = append(items, item)
		}
		order.Items = items

		go r.cache.Set(ctx, &order)
		fmt.Println("взял из бд")
		return &order, nil
	}
	fmt.Println("взял из кеша")
	return order, nil
}
