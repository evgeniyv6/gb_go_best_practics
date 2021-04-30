package main

import (
	"bufio"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"math/rand"
	"net"
	"strconv"
	"time"
)

const (
	timeout = 10 * time.Millisecond
	addr    = "127.0.0.1:5433"
	count   = 250
)

var logger, _ = zap.NewProduction()

func init() {
	rand.Seed(time.Now().UnixNano())
}

func StandartRand(min, max int) string {
	return strconv.Itoa(rand.Intn(max-min) + min)
}

func main() {
	defer logger.Sync()

	debug := flag.Bool("debug", false, "set log level to debug")
	flag.Parse()
	if *debug {
		zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	logger = logger.With(zap.String("name", "Circuit Breaker"))

	orders := make([]string, 0)
	for i := 0; i < count; i++ {
		orders = append(orders, StandartRand(1, 100))
	}

	cb := NewDefaultCB()
	logger.With(zap.String("CB", fmt.Sprintf("%+v", cb))).Info("self Circuit Breaker")

	for _, order := range orders {
		err := cb.Exec(pay(order))
		time.Sleep(1 * time.Second)
		if err != nil {
			logger.Error(fmt.Sprintf("%v", err))
		}
	}

}

func pay(order string) error {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		logger.Error(fmt.Sprintf("could not connect to server - %v", err))
		return err
	}
	fmt.Fprintf(conn, "please process order %q\n", order)
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return err
	}
	logger.With(zap.String("msg from srv", msg))
	return nil
}
