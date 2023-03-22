package main

import (
	"example.com/greetings/internal/app/forum"
	"fmt"
	//fasthttpprom "github.com/701search/fasthttp-prometheus-middleware" //added
	//"github.com/fasthttp/router"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"github.com/georgecookeIW/fasthttprouter"
	"log"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

func main() {
	r := fasthttprouter.New()

	r.HandleOPTIONS = true
	r.HandleCORS.Handle = true
	r.HandleCORS.AllowOrigin = "*"
	r.HandleCORS.AllowMethods = []string{"GET", "POST"}

	//added

	//r.GET("/health", func(ctx *fasthttp.RequestCtx) {
	//	ctx.SetStatusCode(200)
	//	ctx.SetBody([]byte(`{"status": "pass"}`))
	//	log.Println(string(ctx.Request.URI().Path()))
	//})
	//

	DBConf := DBConfig{
		Host: "109.120.182.154",
		Port: "5432",
		Username: "yutfut",
		Password: "yutfut",
		DBName: "yutfut",
	}

	dsn := fmt.Sprintf(`user=%s dbname=%s password=%s host=%s port=%s sslmode=disable`,
		DBConf.Username, DBConf.DBName, DBConf.Password, DBConf.Host, DBConf.Port)

	conn, err := pgx.ParseConnectionString(dsn)
	if err != nil {
		log.Fatalln("cant parse config", err)
	}

	db, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     conn,
		MaxConnections: 1000,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})

	if err != nil {
		fmt.Println("db error")
		fmt.Println(err.Error())
		log.Fatalf("Error %s occurred during connection to database", err)
	}
	fmt.Println("db connect done")

	forum.SetForumRouting(r, &forum.Handlers{
		ForumRepo: forum.NewPgxRepository(db),
	})

	fmt.Printf("Start server on port :8000\n")

	//log.Fatal(fasthttp.ListenAndServe(":5000", r.Handler))
	log.Fatal(fasthttp.ListenAndServe(":8000", r.Handler))
}
