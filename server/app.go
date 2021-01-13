package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PECHIVKO/task-manager/column"
	colhttp "github.com/PECHIVKO/task-manager/column/delivery/http"
	columnrepo "github.com/PECHIVKO/task-manager/column/repository/postgres"
	columnusecase "github.com/PECHIVKO/task-manager/column/usecase"
	"github.com/PECHIVKO/task-manager/config"
	"github.com/PECHIVKO/task-manager/project"
	phttp "github.com/PECHIVKO/task-manager/project/delivery/http"
	projectrepo "github.com/PECHIVKO/task-manager/project/repository/postgres"
	projectusecase "github.com/PECHIVKO/task-manager/project/usecase"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type App struct {
	httpServer *http.Server

	uc *Usecase
}

type Usecase struct {
	projectUC project.UseCase
	columnUC  column.UseCase
}

func NewApp() *App {
	db := initDB()

	return &App{
		uc: NewUC(db),
	}
}

func NewUC(db *sql.DB) *Usecase {
	projectRepo := projectrepo.NewProjectRepository(db)
	columnRepo := columnrepo.NewColumnRepository(db)

	return &Usecase{
		projectUC: projectusecase.NewProjectUseCase(projectRepo),
		columnUC:  columnusecase.NewColumnUseCase(columnRepo),
	}
}

func (a *App) Run(port string) error {
	router := a.uc.Routes()

	a.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	err := a.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to listen and serve: %+v", err)
		return err
	}

	return nil
}

func initDB() *sql.DB {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	mainConfig, openCfgErr := config.NewConfig(configPath)
	if openCfgErr != nil {
		panic("cannot open config: " + openCfgErr.Error())
	}

	conn, err := sql.Open("postgres", mainConfig.Database.DbSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}

func (uc *Usecase) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)
	router.Route("/", func(r chi.Router) {
		r.Mount("/projects", phttp.Routes(uc.projectUC))
		r.Mount("/columns", colhttp.Routes(uc.columnUC))
	})
	return router
}
