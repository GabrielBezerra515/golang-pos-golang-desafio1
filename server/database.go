package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DolarRealModel struct {
	ID         int    `gorm:"primaryKey"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
	gorm.Model
}

func Database(ctx context.Context, dr *DolarReal) error {
	ctxDatabase, cancel := context.WithTimeout(ctx, 10000000*time.Nanosecond)
	defer cancel()

	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gorm.Open: Error to connect database [%s]", err)
	}

	if err := db.AutoMigrate(&DolarRealModel{}); err != nil {
		return fmt.Errorf("db.AutoMigrate: Error to create or identify table [%s]", err)
	}

	var drModel DolarRealModel
	req, err := json.Marshal(&dr.USDBRL)
	if err != nil {
		return fmt.Errorf("json.Marshal: Error marshal data [%s]", err)
	}

	if err := json.Unmarshal(req, &drModel); err != nil {
		return fmt.Errorf("json.UnMarshal: Error UnMarshal data [%s]", err)
	}

	db.WithContext(ctxDatabase).Create(&drModel)
	return nil
}
