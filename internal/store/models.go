package store

type SeatInfo struct {
	Label  string `json:"label"`
	Status string `json:"status"`
}

type Show struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	PricePaise   int64      `json:"price_paise"`
	PerUserLimit int        `json:"per_user_limit"`
	TotalSeats   int        `json:"total_seats"`
	Available    int        `json:"available"`
	Held         int        `json:"held"`
	Confirmed    int        `json:"confirmed"`
	Seats        []SeatInfo `json:"seats"`
}

type Reservation struct {
	ID          string   `json:"reservation_id"`
	ShowID      string   `json:"show_id"`
	UserID      string   `json:"user_id"`
	Seats       []string `json:"seats"`
	AmountPaise int64    `json:"amount_paise"`
	Status      string   `json:"status"`
}
