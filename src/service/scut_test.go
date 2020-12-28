package service

import (
	"log"
	"testing"
)

func TestScut(t *testing.T) {
	scut := &ScutService{&HttpService{}}
	items, err := scut.GetJwItems()
	if err != nil {
		log.Fatalln(err)
	}
	for _, item := range items {
		log.Println(item)
	}
	log.Println()
	items, err = scut.GetSeItems()
	if err != nil {
		log.Fatalln(err)
	}
	for _, item := range items {
		log.Println(item)
	}
}
