package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/totsumaru/gacha-bot-backend/expose"
	"github.com/totsumaru/gacha-bot-backend/expose/api"
	"github.com/totsumaru/gacha-bot-backend/expose/bot/handler"
	"github.com/totsumaru/gacha-bot-backend/lib/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	location = "Asia/Tokyo"
)

func init() {
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

	// .envが存在している場合は読み込み
	if _, err = os.Stat(".env"); err == nil {
		if err = godotenv.Load(); err != nil {
			panic(err)
		}
	}
}

func main() {
	var Token = "Bot " + os.Getenv("APP_BOT_TOKEN")

	session, err := discordgo.New(Token)
	session.Token = Token
	if err != nil {
		panic(err)
	}

	//イベントハンドラを追加
	handler.Handler(session)

	if err = session.Open(); err != nil {
		panic(err)
	}

	// 直近の関数（main）の最後に実行される
	defer func() {
		if err = session.Close(); err != nil {
			panic(err)
		}
		return
	}()

	// DBの設定
	dialector := postgres.Open(os.Getenv("DB_URL"))
	db, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(errors.NewError("DBに接続できません", err))
	}

	expose.DB = db

	// テーブルが存在していない場合のみテーブルを作成します
	// 存在している場合はスキーマを同期します
	//if err = db.AutoMigrate(&applicationDB.Application{}); err != nil {
	//	panic(errors.NewError("テーブルのスキーマが一致しません", err))
	//}

	// Ginの設定
	{
		engine := gin.Default()

		// CORSの設定
		// ここからCorsの設定
		engine.Use(cors.New(cors.Config{
			// アクセスを許可したいアクセス元
			AllowOrigins: []string{
				"http://localhost:3000",
			},
			// アクセスを許可したいHTTPメソッド
			AllowMethods: []string{
				"GET",
				"POST",
				"OPTIONS",
			},
			// 許可したいHTTPリクエストヘッダ
			AllowHeaders: []string{
				"Origin",
				"Content-Length",
				"Content-Type",
				"Authorization",
				"Accept",
				"X-Requested-With",
			},
			ExposeHeaders: []string{"Content-Length"},
			// cookieなどの情報を必要とするかどうか
			AllowCredentials: false,
			// preflightリクエストの結果をキャッシュする時間
			//MaxAge: 24 * time.Hour,
		}))

		// ルートを設定する
		api.RegisterRouter(engine, db)

		if err = engine.Run(":8080"); err != nil {
			log.Fatal("起動に失敗しました", err)
		}
	}

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stopBot
	return
}
