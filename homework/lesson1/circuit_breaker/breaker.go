// образец https://github.com/sony/gobreaker (thanks google halfopen search)

package main

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type State int

const (
	Closed State = iota  // сообщения проходят напрямую к сервису через CB, счетчик ошибок = 0
	Open // CB сразу отбивает сообщения с ошибкой
	HalfOpen // к сервису проходит только часть сообщений
)

var (
	ErrTooManyRequests = errors.New("too many requests")
	ErrOpenState = errors.New("circuit breaker is open")
)

// счетчики для CB
type Counters struct {
	Req int
	TotalGood int
	TotalFail int
	ConsecGood int
	ConsecFail int
}

func (c *Counters) incReq()  {
	c.Req++
	logger.With(zap.String("requests count", strconv.Itoa(c.Req))).Debug("requests")
}

func (c *Counters) incGood() {
	c.TotalGood++
	c.ConsecGood++
	c.ConsecFail = 0
	logger.Debug("increase in success")
	fmt.Printf("c.TotalGood - %s, c.ConsecGood - %s, c.ConsecFail - %s\n", strconv.Itoa(c.TotalGood),
		strconv.Itoa(c.ConsecGood), strconv.Itoa(c.ConsecFail))
	logger.Info(fmt.Sprintf("Total Good Requests - %s, Сonsecutive Good Req - %s, Сonsecutive Fail Req - %s\n",
		strconv.Itoa(c.TotalGood), strconv.Itoa(c.ConsecGood), strconv.Itoa(c.ConsecFail)))
}

func (c *Counters) incFail() {
	c.ConsecFail++
	c.TotalFail++
	c.ConsecGood = 0
	logger.Debug("increase on fail")
	logger.Info(fmt.Sprintf("Сonsecutive Fail Req - %s, Total Fail Req - %s, Сonsecutive Good Req - %s\n",
		strconv.Itoa(c.ConsecFail), strconv.Itoa(c.TotalFail), strconv.Itoa(c.ConsecGood)))
}

func (c *Counters) clear() {
	logger.Debug("counters were cleared")
	c.Req, c.TotalGood, c.TotalFail, c.ConsecGood, c.ConsecFail = 0, 0, 0, 0, 0
}


// сам CB
type CircuitBreaker struct {
	maxReqForHalfOpen int
	timeoutFromOpenToHalfOpen time.Duration
	fromCloseToOpenCounts func(c Counters) bool
	state State
	counts Counters
	expiry time.Time
}

func NewDefaultCB() *CircuitBreaker {
	return &CircuitBreaker{
		maxReqForHalfOpen: 3, // количество запросов CB в состоянии HalfOpen
		timeoutFromOpenToHalfOpen: 2 * time.Second, // время полного локдауна
		fromCloseToOpenCounts: func(c Counters) bool {  // счетчик после которого CB переходит в Open и запускает таймер для ожидания
														// после таймаута CB переходит в HalfOpen
			return c.ConsecFail > 3
		},
		state: Closed,      // init статус Closed, сообщения свободно проходят к сервису
		counts: Counters{},
	}
}


// state rules
func (cb *CircuitBreaker) successReq(st State) {
	switch st {
	case Closed:
		logger.Info("Success req, state is Closed")
		cb.counts.incGood() // если в состоянии Closed, просто увеличиваем счетчик успешных вызовов
	case HalfOpen:
		logger.Info("Success req, state is HalfOpen")
		cb.counts.incGood() // считаем последовательный успешные вызовы
		if cb.counts.ConsecGood >= cb.maxReqForHalfOpen {
			logger.With(zap.Int("count", cb.counts.ConsecGood)).Info("number of consecutive success calls")
			logger.With(zap.Int("count", cb.maxReqForHalfOpen)).Info("halfopen threshold")
			logger.Info("go to Closed state")
			cb.setState(Closed) // переходим в Closed, если не было сбоев
		}

	}
}

func (cb *CircuitBreaker) failReq(st State) {
	switch st {
	case Closed:
		logger.Info("Fail req, state is Closed")
		cb.counts.incFail()
		if cb.fromCloseToOpenCounts(cb.counts) { // после указанного числа ошибок - переходим в Open
			logger.Info("Fail. Go to Open state")
			cb.expiry = time.Now().Add(cb.timeoutFromOpenToHalfOpen)
			cb.setState(Open)
		}
	case HalfOpen:
		logger.Info("Success req, state is HalfOpen")
		cb.expiry = time.Now().Add(cb.timeoutFromOpenToHalfOpen)
		cb.setState(Open) // если были в HalfOpen и получили ошибку, снова переходим в Open

	}
}

// при смене статуса сбрасываем счетчики
func (cb *CircuitBreaker) setState(st State) {
	if cb.state == st {
		return
	}
	cb.state = st
	cb.counts.clear()
}

// CB worker
func (cb *CircuitBreaker) Exec(excErr error) error {
	if cb.state == Open && cb.expiry.Before(time.Now()) {
		cb.expiry = time.Time{}
		logger.Info("set halfopen state")
		cb.setState(HalfOpen)
	}

	if cb.state == Open {
		return ErrOpenState
	} else if cb.state == HalfOpen && cb.counts.Req >= cb.maxReqForHalfOpen {
		return ErrTooManyRequests
	}

	cb.counts.incReq()

	err := excErr
	if err != nil {
		cb.failReq(cb.state)
	} else {
		cb.successReq(cb.state)
	}

	return err
}