package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	api "adtech.simple/internal/app/adtechsimpleapi"
	"adtech.simple/internal/app/jobhandlers"
	"adtech.simple/internal/config"
	"adtech.simple/internal/pkg/adserverclient"
	"adtech.simple/internal/pkg/dbquery"
	"adtech.simple/internal/pkg/jobscheduler"
	"adtech.simple/internal/pkg/model"
	"adtech.simple/internal/pkg/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vgarvardt/gue/v5"
	guepgxv5 "github.com/vgarvardt/gue/v5/adapter/pgxv5"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/********************************************************
				CONFIG SET UP
	 ********************************************************/
	var cfg config.Config
	b, err := os.ReadFile(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatalf("os.ReadFile error: %v", err)
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		log.Fatalf("json.Unmarshal error: %v", err)
	}

	/********************************************************
				PGX SET UP
	 ********************************************************/
	pgxCfg, err := pgxpool.ParseConfig(cfg.PGDSN)
	if err != nil {
		log.Fatalf("pgxpool.ParseConfig error: %v", err)
	}

	pgxPool, err := pgxpool.NewWithConfig(context.Background(), pgxCfg)
	if err != nil {
		log.Fatalf("pgxpool.NewWithConfig error: %v", err)
	}
	defer pgxPool.Close()

	if err := pgxPool.Ping(ctx); err != nil {
		log.Fatalf("poolAdapter.Ping error: %v", err)
	}

	/********************************************************
				AD SERVER CLIENT SET UP
	 ********************************************************/
	adServerClient := adserverclient.NewClient()

	/********************************************************
				GUE SET UP
	 ********************************************************/
	gueClient, err := gue.NewClient(guepgxv5.NewConnPool(pgxPool))
	if err != nil {
		log.Fatalf("gue.NewClient error: %v", err)
	}

	dbQuerier := dbquery.New(pgxPool)
	campaignCreationHandler := jobhandlers.NewCampaignCreationHandler(adServerClient, dbQuerier)

	wm := gue.WorkMap{
		fmt.Sprint(model.JobTypeCampaignCreation): campaignCreationHandler.MakeHandler(),
	}
	workers, err := gue.NewWorkerPool(gueClient, wm, 2, gue.WithPoolQueue(fmt.Sprint(model.QueueTypeCampaignCreation)))
	if err != nil {
		log.Fatalf("gue.NewWorkerPool error: %v", err)
	}

	go func() {
		if err := workers.Run(ctx); err != nil {
			log.Printf("workers.Run error: %v", err)
		}
	}()

	/********************************************************
				STORAGE SET UP
	 ********************************************************/
	storage := store.NewStorage(pgxPool)

	/********************************************************
				SCHEDULER SET UP
	 ********************************************************/
	jobScheduler := jobscheduler.NewScheduler(gueClient)

	/********************************************************
				HTTP
	 ********************************************************/
	createCampaignHandler := api.NewCreateCampaignHandler(storage, adServerClient, jobScheduler)
	addViewHandler := api.NewAddViewHandler(adServerClient)

	mux := http.NewServeMux()
	mux.Handle("/create-campaign", createCampaignHandler)
	mux.Handle("/add-view", addViewHandler)

	log.Println("adtech-simple is running on port 9090")
	log.Fatal(http.ListenAndServe(":9090", mux))
}
