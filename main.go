package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func main() {

	if err := pingKafka(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("consumer is ready!")

	if err := pingMysql(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ping mysql is ok")

	time.Sleep(100 * time.Hour)
}

func pingMysql() error {
	connect := fmt.Sprintf("%v:%v@tcp(%v:%v)/mysql?charset=utf8",
		os.Getenv("MYSQL_USER_NAME"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"))
	fmt.Println(connect)
	_, err := xorm.NewEngine("mysql", connect)
	if err != nil {
		return err
	}

	return nil
}

func pingKafka() error {
	notifyKill()
	fmt.Printf(">> env  addr:%v,port:%v,topic:%v \n", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"), os.Getenv("KAFKA_TOPIC"))
	config := Config{
		Brokers: []string{os.Getenv("KAFKA_HOST") + ":" + os.Getenv("KAFKA_PORT")},
		Topic:   os.Getenv("KAFKA_TOPIC"),
	}
	if err := Consume("consumer-kafka", config,
		EventFruit,
	); err != nil {
		return err
	}

	return nil

}

func notifyKill() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)
	go func() {
		for s := range signals {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				os.Exit(0)
			}
		}
	}()

}
