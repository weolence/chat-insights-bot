package controller

import (
	"github.com/cornelk/hashmap"
	"gopkg.in/telebot.v3"
)

type BtnRouteController struct {
	routes *hashmap.Map[string, func(telebot.Context) error]
}

func NewBtnRouteController() *BtnRouteController {
	return &BtnRouteController{
		routes: hashmap.New[string, func(telebot.Context) error](),
	}
}

func (btnRouteController *BtnRouteController) SetRoute(key string, handler func(telebot.Context) error) {
	btnRouteController.routes.Set(key, handler)
}

func (btnRouteController *BtnRouteController) GetRoute(key string) (func(telebot.Context) error, bool) {
	return btnRouteController.routes.Get(key)
}

func (btnRouteController *BtnRouteController) DeleteRoute(key string) {
	btnRouteController.routes.Del(key)
}
