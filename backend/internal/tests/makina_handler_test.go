package tests

// func TestGetPozisyonlar(t *testing.T) {

// 	config.LoadConfig("/Users/enesaygun/Desktop/sasa-repositories/sasa-elterminali-service/cmd/api/config")
// 	app := fiber.New()

// 	postgresDB, err := postgres.NewPostgresDB()
// 	if err != nil {
// 		log.Fatalf("Failed to connect to Postgres: %v", err)
// 	}

// 	DB := repository.NewDB(postgresDB.Pool, nil)

// 	services := services.NewMakinaService(DB)
// 	handler := handler.NewMakinaHandler(*services)
// 	app.Post("/pozisyonlarGet", handler.GetPozisyonlar)

// 	var requestBody models.PozisyonlarGetRequest

// 	requestBody.BlokIDs = []int{1, 2}

// 	jsonData, err := json.Marshal(requestBody)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal request body: %v", err)
// 	}

// 	reqBody := bytes.NewBuffer(jsonData)

// 	req := httptest.NewRequest("POST", "/pozisyonlarGet", reqBody)
// 	resp, _ := app.Test(req, 10000)

// 	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
// }
