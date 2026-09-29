package main

import (
	"campus-lost-found-api/app/model"
	"campus-lost-found-api/app/repository"
	"campus-lost-found-api/config"
	"campus-lost-found-api/database"
	"campus-lost-found-api/helper"
	"context"
	"flag"
	"fmt"
	"log"
)

func main() {
	u := flag.String("username", "admin", "username")
	e := flag.String("email", "admin@example.com", "email")
	p := flag.String("password", "Admin123!", "password")
	flag.Parse()
	cfg := config.Load()
	pool, err := database.NewPool(context.Background(), cfg.DatabaseURL, cfg.DBMinConns, cfg.DBMaxConns)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	repo := repository.NewUserRepository(pool)
	h, err := helper.HashPassword(*p)
	if err != nil {
		log.Fatal(err)
	}
	user, err := repo.Create(context.Background(), model.User{Username: *u, Email: *e, Password: h})
	if err != nil {
		log.Fatal(err)
	}
	if err := repo.UpdateRole(context.Background(), user.ID, "admin"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("admin created: id=%d username=%s\n", user.ID, user.Username)
}
