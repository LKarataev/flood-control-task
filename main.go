package main

import (
	"context"

	ctrl "github.com/LKarataev/flood-control-task/internal/flood_control"
)

func main() {
	id := 123
	ctx := ctx.Background()
	fmt.Println(Check(ctx, id))
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}