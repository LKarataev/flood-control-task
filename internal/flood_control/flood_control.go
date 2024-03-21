package flood_control


type Api struct {

}

func Check(ctx context.Context, userID int64) (bool, error) {
	fmt.Println("userID: ", userID)
	fmy.Printf("%#v\n", ctx)
	return true, nil
}

// type FloodControl interface {
// 	// Check возвращает false если достигнут лимит максимально разрешенного
// 	// кол-ва запросов согласно заданным правилам флуд контроля.
// 	Check(ctx context.Context, userID int64) (bool, error)
// }
