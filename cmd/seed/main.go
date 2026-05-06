package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/rating"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/rewards"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalln("Error: DB_URL no encontrada")
	}

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalln("Error: No se pudo conectar con la base de datos")
	}
	defer dbPool.Close()

	store := db.NewStore(dbPool)
	accountService := account.NewService(store)
	ratingService := rating.NewService(store)

	// INSERCIONES CORE (Siempre)
	log.Println("--- Insertando datos CORE ---")
	seedShopItems(ctx, dbPool)
	seedAIUser(ctx, accountService)

	// INSERCIONES DEBUG (Solo si SEED_DEBUG=true)
	isDebug := os.Getenv("SEED_DEBUG") == "true"
	if isDebug {
		log.Println("--- Insertando datos DEBUG ---")
		seedDebugUsers(ctx, accountService, ratingService)
		seedDebugMatchesAndFriends(ctx, dbPool)
	} else {
		log.Println("--- Saltando datos DEBUG ---")
	}

	log.Println("Seeding completado con éxito.")
}

func seedAIUser(ctx context.Context, accSvc *account.AccountService) {
	log.Println("Iniciando la inserción del usuario IA...")

	username := fmt.Sprintf("AI")
	password := fmt.Sprintf(os.Getenv("AI_PASSWORD"))
	mail := fmt.Sprintf("ai@ai.ai")

	// Comprobar si ya existe
	_, err := accSvc.GetIDByUsername(ctx, username)
	if err == nil {
		return
	}

	dto := &account.CreateAccountDTO{
		Username: username,
		Mail:     mail,
		Password: password,
	}

	// Crear la cuenta
	_, err = accSvc.Create(ctx, dto)
	if err != nil {
		log.Printf("Error creando a %s: %v", username, err)
		return
	}

	if err != nil {
		log.Printf("Error actualizando stats (nivel/xp) para %s: %v", username, err)
	}

	log.Println("Usuarios IA validado/creado correctamente.")
}

func seedShopItems(ctx context.Context, dbPool *pgxpool.Pool) {
	// Comprobar si ya hay items para no duplicar
	var count int
	err := dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM shop_item").Scan(&count)
	if err == nil && count > 0 {
		log.Println("Los items ya existen, saltando inserción...")
		return
	}

	// Insertamos los items explícitamente con su ID para mantener el orden exacto.
	query := `
	INSERT INTO shop_item (item_id, price, level_requisite, item_type, is_default) VALUES 
		(1, 0, 0, 'PIECE_SKIN', true),      -- Classic
		(2, 500, 5, 'PIECE_SKIN', false),   -- Soretro (Nivel 5)
		(3, 1000, 10, 'PIECE_SKIN', false), -- Cats (Nivel 10)
		(4, 0, 0, 'BOARD_SKIN', true),      -- Classic
		(5, 500, 5, 'BOARD_SKIN', false),   -- Soretro (Nivel 5)
		(6, 1000, 10, 'BOARD_SKIN', false), -- Cats (Nivel 10)
		(7, 0, 0, 'WIN_ANIMATION', true),   -- Classic
		(8, 500, 5, 'WIN_ANIMATION', false),-- Soretro (Nivel 5)
		(9, 1000, 10, 'WIN_ANIMATION', false),-- Cats (Nivel 10)
		(10, 0, 0, 'AVATAR', true),         -- bot1_lila (Defecto)
		(11, 100, 2, 'AVATAR', false),      -- bot2_amarillo
		(12, 150, 3, 'AVATAR', false),      -- bot3_magenta
		(13, 200, 4, 'AVATAR', false),      -- bot4_naranja
		(14, 250, 6, 'AVATAR', false),      -- bot5_amarillo
		(15, 300, 7, 'AVATAR', false),      -- bot6_magenta
		(16, 400, 8, 'AVATAR', false),      -- bot7_rojo
		(17, 500, 10, 'AVATAR', false),     -- bot8_rojo
		(18, 600, 11, 'AVATAR', false),     -- bot9_verde
		(19, 750, 13, 'AVATAR', false),     -- bot10_lila
		(20, 900, 14, 'AVATAR', false),     -- bot11_amarillo
		(21, 1200, 15, 'AVATAR', false)     -- bot12_verde (Nivel Máximo 15)
	ON CONFLICT (item_id) DO NOTHING;
	`
	_, err = dbPool.Exec(ctx, query)
	if err != nil {
		log.Printf("Error insertando items: %v", err)
	} else {
		log.Println("Items base insertados correctamente en orden.")

		_, err = dbPool.Exec(ctx, "SELECT setval('shop_item_item_id_seq', (SELECT MAX(item_id) FROM shop_item));")
		if err != nil {
			log.Printf("Error actualizando la secuencia de item_id: %v", err)
		}
	}
}

func seedDebugUsers(ctx context.Context, accSvc *account.AccountService, ratingSvc *rating.RatingService) {
	log.Println("Iniciando la inserción de usuarios con Elo, Nivel y XP random...")

	for i := 1; i <= 150; i++ { // (Ajusta a 10 o 150 según necesites)
		username := fmt.Sprintf("user%d", i)
		password := fmt.Sprintf("%s%s", username, username)
		mail := fmt.Sprintf("user%d@gmail.com", i)

		// Comprobar si ya existe
		_, err := accSvc.GetIDByUsername(ctx, username)
		if err == nil {
			continue
		}

		dto := &account.CreateAccountDTO{
			Username: username,
			Mail:     mail,
			Password: password,
		}

		// Crear la cuenta
		accDTO, err := accSvc.Create(ctx, dto)
		if err != nil {
			log.Printf("Error creando a %s: %v", username, err)
			continue
		}

		randomBlitz := int32(800 + rand.Intn(1201))
		randomRapid := int32(800 + rand.Intn(1201))
		randomClassic := int32(800 + rand.Intn(1201))
		randomExtended := int32(800 + rand.Intn(1201))

		ratingSvc.UpdateEloByID(ctx, &rating.RatingDTO{UserID: *accDTO.AccountID, EloType: db.EloTypeBLITZ, Value: randomBlitz, Deviation: 350, Volatility: 0.06})
		ratingSvc.UpdateEloByID(ctx, &rating.RatingDTO{UserID: *accDTO.AccountID, EloType: db.EloTypeRAPID, Value: randomRapid, Deviation: 350, Volatility: 0.06})
		ratingSvc.UpdateEloByID(ctx, &rating.RatingDTO{UserID: *accDTO.AccountID, EloType: db.EloTypeCLASSIC, Value: randomClassic, Deviation: 350, Volatility: 0.06})
		ratingSvc.UpdateEloByID(ctx, &rating.RatingDTO{UserID: *accDTO.AccountID, EloType: db.EloTypeEXTENDED, Value: randomExtended, Deviation: 350, Volatility: 0.06})

		// XP aleatoria entre 0 y 99999
		randomXp := int32(rand.Intn(100000))
		randomLevel := rewards.GetLevel(randomXp)

		// Dinero aleatorio entre 0 y 5000
		randomMoney := int32(rand.Intn(5001))

		err = accSvc.UpdateStats(ctx, *accDTO.AccountID, &account.AccountStatsDTO{
			Level: randomLevel,
			Xp:    randomXp,
			Money: randomMoney,
		})

		if err != nil {
			log.Printf("Error actualizando stats (nivel/xp) para %s: %v", username, err)
		}

		log.Printf("Usuario %s creado | Nivel: %d | XP: %d", username, randomLevel, randomXp)
	}

	log.Println("Usuarios validados/creados correctamente.")
}

func seedDebugMatchesAndFriends(ctx context.Context, dbPool *pgxpool.Pool) {
	// Amistades
	friendshipQuery := `
	INSERT INTO public."friendship" (user1_id, user2_id, accepted_1, accepted_2) VALUES 
		(2, 3, true, true),
		(2, 4, false, true),
		(2, 5, true, false),
		(4, 5, true, true),
		(4, 6, true, true),
		(4, 7, true, true),
		(4, 8, true, true),
		(4, 9, true, true),
		(4, 10, true, true),
		(4, 11, true, true),
		(4, 12, true, true),
		(4, 13, true, true),
		(4, 14, true, true),
		(4, 15, true, true),
		(4, 16, true, true),
		(4, 17, true, true),
		(4, 18, true, true),
		(4, 18, true, true),
		(4, 20, true, true),
		(4, 21, true, true),
		(4, 22, true, true),
		(4, 23, true, true),
		(4, 24, true, true),
		(4, 25, true, true),
		(4, 26, true, true),
		(4, 27, true, true),
		(4, 28, true, true),
		(4, 29, true, true),
		(4, 30, true, true),
		(4, 31, true, true),
		(4, 32, true, true),
		(4, 33, true, true),
		(4, 34, true, true),
		(4, 35, true, true),
		(4, 36, true, true),
		(4, 37, true, true),
		(4, 38, true, true),
		(4, 39, true, true),
		(4, 40, true, true)
	ON CONFLICT (user1_id, user2_id) DO NOTHING;
	`
	dbPool.Exec(ctx, friendshipQuery)

	// Partidas
	var matchCount int
	dbPool.QueryRow(ctx, `SELECT COUNT(*) FROM public."match" WHERE p1_id = 2`).Scan(&matchCount)

	if matchCount == 0 {
		matchQuery := `
		INSERT INTO public."match" (
			p1_id, p2_id, p1_elo, p2_elo, "date", "winner", "termination", 
			"match_type", board, movement_history, time_base, time_increment
		) VALUES 
		(2, 3, 1500, 1600, '2026-02-22T15:04:05Z', 'P1_WINS', 'LASER', 'RANKED', 'ACE', '', 300, 5),
		(2, 3, 1500, 1600, '2026-02-22T15:04:05Z', 'NONE', 'UNFINISHED', 'RANKED', 'CURIOSITY', 'Rf1%j1,j4,i4,i5,j5,j9%{300};Tg6:f6%a8,a5,b5,b4,a4,a0%{300};Rb4%j1,j4,i4,i5,j5,j9%{295};Ri5xf6%a8,a5,b5,b4,e4,e5,f5,f6%{290};', 300, 0);
		`
		_, err := dbPool.Exec(ctx, matchQuery)
		if err != nil {
			log.Printf("Error insertando partidas: %v", err)
		} else {
			log.Println("Partidas insertadas correctamente.")
		}
	}
}
