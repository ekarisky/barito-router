package router

import (
	"testing"

	. "github.com/BaritoLog/go-boilerplate/testkit"
)

func TestTrader_New(t *testing.T) {
	want := "http://host:port"
	trader := NewTrader(want)

	got := trader.Url()
	FatalIf(t, got != want, "%s != %s", got, want)
}

func TestTrader_Trade(t *testing.T) {
	trader := NewTrader("http://host:port")

	item, err := trader.Trade("secret")
	FatalIfError(t, err)
	FatalIf(t, item == nil, "item can't be nil")

}
