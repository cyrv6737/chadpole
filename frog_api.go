package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type FrogJSON struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Link  string `json:"link"`
	Image string `json:"imageURL"`
}

func FrogAPIHandler(w http.ResponseWriter, r *http.Request) {
	response := []FrogJSON{
		{
			Name:  "Red-eyed Tree Frog",
			Desc:  "Known for its vibrant green coloration, the red-eyed tree frog is native to the rainforests of Central and South America. It has striking red eyes and orange feet, and it spends most of its life in trees near water bodies.",
			Link:  "https://en.wikipedia.org/wiki/Agalychnis_callidryas",
			Image: "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e3/Red-eyed_Tree_Frog_%28Agalychnis_callidryas%29_1.png/220px-Red-eyed_Tree_Frog_%28Agalychnis_callidryas%29_1.png",
		},
		{
			Name:  "Poison Dart Frog",
			Desc:  "Poison dart frogs are small and brightly colored frogs found in Central and South America. They are known for their toxic skin secretions, which have been used by indigenous people to poison the tips of blowgun darts for hunting.",
			Link:  "https://en.wikipedia.org/wiki/Poison_dart_frog",
			Image: "https://upload.wikimedia.org/wikipedia/commons/thumb/0/0e/Blue-poison.dart.frog.and.Yellow-banded.dart.frog.arp.jpg/220px-Blue-poison.dart.frog.and.Yellow-banded.dart.frog.arp.jpg",
		},
		{
			Name:  "African Clawed Frog",
			Desc:  "Native to sub-Saharan Africa, the African clawed frog is an aquatic species that lacks tongue and teeth. Its distinctive \"claws\" on its hind feet are used for digging and defense.",
			Link:  "https://en.wikipedia.org/wiki/African_clawed_frog",
			Image: "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b4/Xenopus_laevis_02.jpg/220px-Xenopus_laevis_02.jpg",
		},
		{
			Name:  "Goliath Frog",
			Desc:  "The Goliath frog is the largest frog species, found in Cameroon and Equatorial Guinea in Africa. It can grow to be over a foot long and has a unique appearance with robust body proportions.",
			Link:  "https://en.wikipedia.org/wiki/Goliath_frog",
			Image: "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d7/Goliath_Frog.jpg/220px-Goliath_Frog.jpg",
		},
		{
			Name:  "Wood Frog",
			Desc:  "Wood frogs are found in North America and are known for their remarkable adaptation to cold environments. They can survive freezing temperatures by entering a state of suspended animation and then thawing back to life when temperatures rise.",
			Link:  "https://en.wikipedia.org/wiki/Wood_frog",
			Image: "https://images.saymedia-content.com/.image/t_share/MTc0NjQ2MTIzODkyNTgxNzU0/frozen-wood-frogs-and-adaptations-for-survival.jpg",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("[FATAL] API Error")
	}
}

func StartFrogAPI() {
	http.HandleFunc("/frog", FrogAPIHandler)
	http.ListenAndServe(":8081", nil)
}
