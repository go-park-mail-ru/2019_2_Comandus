package chat_app

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients"
	chgrpc "github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/chat/delivery/grpc"
	chat_rep "github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/chat/repository"
	chat_ucase "github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/chat/usecase"
	mes_rep "github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/message/repository"
	mes_ucase "github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/message/usecase"
	store "github.com/go-park-mail-ru/2019_2_Comandus/internal/store/create"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Start() error {
	flag.Parse()

	config := NewConfig()
	db, _ := newDB(config.DatabaseURL)

	log.SetFlags(log.Lshortfile)


	mUcase := mes_ucase.NewMessageUsecase(mes_rep.NewMessageRepository(db))
	chUcase := chat_ucase.NewChatUsecase(chat_rep.NewChatRepository(db))

	go func() {
		lis, err := net.Listen("tcp", clients.CHAT_PORT)
		if err != nil {
			log.Fatalln("cant listet port", err)
		}
		server := grpc.NewServer()
		chgrpc.NewChatServerGrpc(server, chUcase)

		fmt.Println("starting server at ", clients.CHAT_PORT)
		server.Serve(lis)
	}()

	server := NewServer("/entry", mUcase, chUcase)
	go server.Listen()

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	fmt.Println("starting chat_app at :8089")
	return http.ListenAndServe(":8089", nil)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	if err := store.CreateChatTables(db); err != nil {
		return nil, err
	}
	return db, nil
}
