package account

import (
	"io"
	"net/http"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog"
)

type BankAccount struct {
	Account_id    int     `json:"id"`
	User_id       int     `json:"id_user"`
	Number        string  `json:"number"`
	Original_name string  `json:"original_name"`
	Balance       float32 `json:"balance"`
	Last_update   string  `json:"last_update"`
	Iban          string  `json:"iban"`
	Currency      string  `json:"currency"`
	Account_type  string  `json:"type"`
	Error         string  `json:"error"` // not needed ?
	Usage         string  `json:"usage"`
}

var logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}).With().Timestamp().Logger()

func AccountScreen(_ fyne.Window) fyne.CanvasObject {

	resp, err := http.Get("http://localhost:8080/version/")
	if err != nil {
		logger.Error().Err(err).Msg("Cannot run http get request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg("ReadAll error")
	}

	sb := string(body)
	logger.Info().Msg("Successfully ping backend")

	return widget.NewLabel(sb)
}
