package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/omarelshopky/craftlist/internal/app"
)

func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    if err := app.Run(ctx); err != nil {
        os.Exit(1)
    }
}