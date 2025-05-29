package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	db "github.com/saran-pt/livetrade-engine/pkg/sql/dbal"
	"github.com/saran-pt/livetrade-engine/pkg/utils"
)

type ApiConfig struct {
	DB *db.Queries
}

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Name string `json:"name"`
	}

	params := parameter{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
	}

	user, err := cfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Balance:   0.0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error while creating User: %v", err))
		return
	}
	utils.RespondWithJson(w, 200, user)
}

func (cfg *ApiConfig) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	// type parameter struct {
	// 	Side     string  `json:"side"`
	// 	Price    float64 `json:"price"`
	// 	Quantity int     `json:"quantity"`
	// 	UserId   string  `json:"userid"`
	// }

	params := order{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
	}

	remainingQuantity := fillOrders(params.Side, params.Price, params.Quantity, params.UserId)

	if remainingQuantity == 0 {
		utils.RespondWithJson(w, 200, map[string]interface{}{"filledQuantity": params.Quantity})
		return
	}

	if params.Side == "bid" {
		bids = append(bids, order{
			Price:    params.Price,
			Quantity: remainingQuantity,
			UserId:   params.UserId,
			Side:     params.Side,
		})
		sort.Slice(bids, func(i, j int) bool {
			return bids[i].Price < bids[j].Price
		})
	} else {
		asks = append(asks, order{
			Price:    params.Price,
			Quantity: remainingQuantity,
			UserId:   params.UserId,
			Side:     params.Side,
		})
		sort.Slice(asks, func(i, j int) bool {
			return asks[i].Price < asks[j].Price
		})
	}
	utils.RespondWithJson(w, 200, map[string]interface{}{"filledQuantity": (params.Quantity - remainingQuantity)})
}

func (cfg *ApiConfig) GetDepth(w http.ResponseWriter, r *http.Request) {
	type values struct {
		Type     string
		Quantity int
	}

	depth := map[string]values{}

	for _, bid := range bids {
		if _, ok := depth[fmt.Sprint(bid.Price)]; ok {
			depth[fmt.Sprint(bid.Price)] = values{
				Type:     "bid",
				Quantity: bid.Quantity,
			}
		} else {
			item := depth[fmt.Sprint(bid.Price)]
			item.Quantity += bid.Quantity
			depth[fmt.Sprint(bid.Price)] = item
		}
	}
	utils.RespondWithJson(w, 200, map[string]interface{}{"depth": depth})
}

func (cfg *ApiConfig) GetBalance(w http.ResponseWriter, r *http.Request) {}

func (cfg *ApiConfig) GetQuote(w http.ResponseWriter, r *http.Request) {}

// Structs and Helper Functions
type balance struct {
	Values map[string]int
}

type user struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	Balance balance
}

type order struct {
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	UserId   string  `json:"userid"`
	Side     string  `json:"side"`
}

var bids []order
var asks []order
var users []user
var TICKER = "GOOGLE"

func flipBalance(userId1 string, userId2 string, quantity int, price float64) {
	user1, user2 := user{}, user{}
	for _, user := range users {
		if user.Id == userId1 {
			user1 = user
		}
		if user.Id == userId2 {
			user2 = user
		}
	}
	if user1.Id == "" || user2.Id == "" {
		return
	}
	user1.Balance.Values[TICKER] -= quantity
	user2.Balance.Values[TICKER] += quantity
	user1.Balance.Values["USD"] += (quantity * int(price))
	user2.Balance.Values["USD"] -= (quantity * int(price))
}

func fillOrders(side string, price float64, quantity int, userId string) int {
	remainingQuantity := quantity
	if side == "bid" {
		for _, ask := range asks {
			if ask.Price > price {
				continue
			}
			if ask.Quantity > remainingQuantity {
				ask.Quantity -= remainingQuantity
				flipBalance(ask.UserId, userId, remainingQuantity, ask.Price)
				return 0
			} else {
				remainingQuantity -= ask.Quantity
				flipBalance(ask.UserId, userId, ask.Quantity, ask.Price)
				asks = bids[:len(bids)-1]
			}
		}
	} else {
		for _, bid := range bids {
			if bid.Price < price {
				continue
			}
			if bid.Quantity > remainingQuantity {
				bid.Quantity -= remainingQuantity
				flipBalance(userId, bid.UserId, remainingQuantity, price)
				return 0
			} else {
				remainingQuantity -= bid.Quantity
				flipBalance(userId, bid.UserId, bid.Quantity, price)
				bids = bids[:len(bids)-1]
			}
		}
	}
	return remainingQuantity
}
