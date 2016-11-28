package main

import (
	"github.com/bwmarrin/discordgo"
)

type (
	method func(*discordgo.Session, *discordgo.Message)
	router struct {
		routes map[string]method
	}
)

func newRouter() (*router) {
	return &router {
		routes: make(map[string]method),
	}
}

func (r *router) addRoute(command string, action method) {
	r.routes[command] = action
}

func (r *router) getRoute(s *discordgo.Session, m *discordgo.Message) {
	command := getMethod(m)
	if r.routes[command] != nil {
		go r.routes[command](s, m)
	}
}
