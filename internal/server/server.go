package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"guthub.com/serge64/joffer/internal/config"
	"guthub.com/serge64/joffer/internal/storage/storagepg"

	"github.com/sirupsen/logrus"
)

func Start(c *config.Config) error {
	store, err := storagepg.New(c.DatabaseURL)
	if err != nil {
		return err
	}

	defer store.Close()

	handler := newHandler(store)
	srv := &http.Server{
		Addr:         c.HostAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logrus.Error(err)
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)
	logrus.Info("shutting server")

	return nil
}
