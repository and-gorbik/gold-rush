package earners

// func Test_TreasuresEarner_Run_Success(t *testing.T) {
// 	goodDigger := testmocks.GoodProvider{}

// 	licenses := make(chan int, 100)

// 	go func() {
// 		for i := 0; i < 100; i++ {
// 			<-time.After(time.Millisecond)
// 			licenses <- i
// 		}
// 	}()

// 	te := NewTreasuresEarner(
// 		config.Entity{},
// 		goodDigger,
// 		testmocks.GetFakeQueue(),
// 		licenses,
// 	)

// 	i := 0
// 	for treasures := range te.Treasures() {
// 		log.Println(i, treasures)
// 		i++
// 	}

// 	<-time.After(5 * time.Second)
// }
