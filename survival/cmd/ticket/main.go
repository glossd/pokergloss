package main

import (
	"github.com/glossd/pokergloss/survival/db"
	"log"
)

func main() {
	ctx, dbClient, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Disconnect(ctx)

	//find, err := db.Client.Database("profile").Collection("usernames").Find(nil, bson.D{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//type Profile struct {
	//	Username string `bson:"_id"`
	//	UserID string
	//}
	//for find.Next(nil) {
	//	var p Profile
	//	err := find.Decode(&p)
	//	if err != nil {
	//		log.Println("Failed to decode")
	//		continue
	//	}
	//	err = db.updateCard(nil, p.UserID, 2)
	//	if err != nil {
	//		log.Printf("Failed to insert: %s\n", err)
	//	}
	//}

	err = db.CardIncTwoTickets(nil, "cLhhrCC0sMUdpIHGlngodn1lKLl2")
	if err != nil {
		log.Fatal(err)
	}
}
