package main

import (
	"example.com/greetings/internal/app/forum"
	"example.com/greetings/internal/app/post"
	"example.com/greetings/internal/app/service"
	"example.com/greetings/internal/app/thread"
	"example.com/greetings/internal/app/user"
	"fmt"
	//fasthttpprom "github.com/701search/fasthttp-prometheus-middleware" //added
	fasthttpprom "example.com/greetings/internal/app/metrics"
	"github.com/fasthttp/router"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
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
	r := router.New()

	//added
	p := fasthttpprom.NewPrometheus("")
	p.Use(r)

	//r.GET("/health", func(ctx *fasthttp.RequestCtx) {
	//	ctx.SetStatusCode(200)
	//	ctx.SetBody([]byte(`{"status": "pass"}`))
	//	log.Println(string(ctx.Request.URI().Path()))
	//})
	//

	DBConf := DBConfig{
		Host: "127.0.0.1",
		Port: "5432",
		Username: "docker",
		Password: "docker",
		DBName: "docker",
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

	user.SetServiceRouting(r, &user.Handlers{
		UserRepo: user.NewPgxRepository(db),
	})

	forum.SetForumRouting(r, &forum.Handlers{
		ForumRepo: forum.NewPgxRepository(db),
	})

	thread.SetThreadRouting(r, &thread.Handlers{
		ThreadRepo: thread.NewPgxRepository(db),
	})

	post.SetPostRouting(r, &post.Handlers{
		PostRepo: post.NewPgxRepository(db),
	})

	service.SetServiceRouting(r, &service.Handlers{
		ServiceRepo: service.NewPgxRepository(db),
	})
	fmt.Printf("Start server on port :5000")

	//log.Fatal(fasthttp.ListenAndServe(":5000", r.Handler))
	log.Fatal(fasthttp.ListenAndServe(":5000", p.Handler))
}
