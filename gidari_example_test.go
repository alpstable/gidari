package gidari_test

//func ExampleNewIterator() {
//	ctx := context.Background()
//
//	const api = "https://anapioficeandfire.com/api"
//
//	// Create the HTTP Requests to iterate over.
//	bookReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, api+"/books", nil)
//	houseReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, api+"/houses", nil)
//	charReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, api+"/characters", nil)
//
//	// Create the iterator object.
//	iter, _ := gidari.NewIterator(ctx, &gidari.Config{
//		Requests: []*gidari.Request{
//			{
//				HttpRequest: bookReq,
//				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 5),
//			},
//			{
//				HttpRequest: houseReq,
//				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 5),
//			},
//			{
//				HttpRequest: charReq,
//				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 5),
//			},
//		},
//	})
//
//	defer iter.Close(ctx)
//
//	// byteSize will keep track of the sum of bytes for each HTTP Response's body.
//	var byteSize int
//
//	for iter.Next(ctx) {
//		// Get the byte slice from the response body.
//		body, err := io.ReadAll(iter.Current.Body)
//		if err != nil {
//			log.Fatalf("failed to read response body: %v", err)
//		}
//
//		// Add the number of bytes to the sum.
//		byteSize += len(body)
//	}
//
//	fmt.Println("Total number of bytes:", byteSize)
//	// Output:
//	// Total number of bytes: 256179
//}
//
//func ExampleTransport() {
//	ctx := context.Background()
//
//	const api = "https://anapioficeandfire.com/api"
//
//	// Create the HTTP Requests to iterate over.
//	bookReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, api+"/books", nil)
//	houseReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, api+"/houses", nil)
//	charReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, api+"/characters", nil)
//
//	// Initiate the transport
//	err := gidari.Transport(ctx, &gidari.Config{
//		Requests: []*gidari.Request{
//			{
//				HttpRequest: bookReq,
//				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 5),
//			},
//			{
//				HttpRequest: houseReq,
//				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 5),
//			},
//			{
//				HttpRequest: charReq,
//				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 5),
//			},
//		},
//	})
//
//	if err != nil {
//		log.Fatalf("failed to transport: %v", err)
//	}
//
//	fmt.Println("Transported successfully")
//	// Output:
//	// Transported successfully
//}
