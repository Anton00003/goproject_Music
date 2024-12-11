package main

import (
	"fmt"
	"goproject_Music/api"
	"goproject_Music/client"
	"goproject_Music/repository"
	"goproject_Music/service"

	"goproject_Music/config"

	"github.com/sirupsen/logrus"
)

func main() {
	conf, err := config.New()
	if err != nil {
		fmt.Println(err, "config init failed")
		return
	}
	lvl, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug(err, "parse level failed")
	}
	logrus.SetLevel(lvl)
	client := client.NewClient(conf.MusicSrorageHost, conf.Path)
	repo, err := repository.NewRepo(conf.DSN)
	if err != nil {
		logrus.Error(err)
		return
	}
	serv := service.NewServ(repo, client)
	api := api.NewApi(serv)
	api.Run(conf.Host)

}
