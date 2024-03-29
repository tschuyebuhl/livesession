package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tschuyebuhl/livesession/src/cache"
	"github.com/tschuyebuhl/livesession/src/data"
	"github.com/tschuyebuhl/livesession/src/service"
	"log/slog"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), "postgres://livesession:livesession@db:5432/livesession")
	if err != nil {
		panic(err)
	}
	err = dbpool.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	c := cache.NewInMemory()
	userRepo := data.NewPostgresRepository(dbpool)
	userService := service.NewUserService(userRepo, c)

	user, cached, err := userService.GetUser("PeterGonzalesisfna")
	if err != nil {
		slog.Error("error getting user", "error", err)
		return
	}

	println(cached)
	user, cached, err = userService.GetUser("PeterGonzalesisfna")
	if err != nil {
		slog.Error("error getting user", "error", err)
		return
	}
	println(cached)
	println(user.Name)
	user, cached, err = userService.GetUser("KelseyElliswvnxi")
	if err != nil {
		slog.Error("error getting user", "error", err)
		return
	}
}
