package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	supa "github.com/supabase-community/supabase-go"
)

type DB struct {
	*sql.DB
	SupabaseClient *supa.Client
}

func NewConnection(databaseURL, supabaseURL, supabaseKey string) (*DB, error) {
	// Traditional SQL connection
	sqlDB, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Supabase client - Updated syntax
	supabaseClient, err := supa.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create supabase client: %w", err)
	}

	log.Println("Database connection established")

	return &DB{
		DB:             sqlDB,
		SupabaseClient: supabaseClient,
	}, nil
}
