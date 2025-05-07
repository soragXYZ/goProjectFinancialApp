package transactions

import (
	"encoding/json"
	"fmt"
	"freenahiFront/internal/helper"
	"freenahiFront/internal/settings"
	"io"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// The struct which is returned by the backend
type Transaction struct {
	Date             string  `json:"date"`
	Value            float32 `json:"value"`
	Transaction_type string  `json:"type"`
	Original_wording string  `json:"original_wording"`
}

// Call the backend endpoint "/transaction" and retrieve txs of the selected page
func getTransactions(page int, app fyne.App) []Transaction {

	backendIp := app.Preferences().StringWithFallback(settings.PreferenceBackendIP, settings.BackendIPDefault)
	backendProtocol := app.Preferences().StringWithFallback(settings.PreferenceBackendProtocol, settings.BackendProtocolDefault)
	backendPort := app.Preferences().StringWithFallback(settings.PreferenceBackendPort, settings.BackendPortDefault)

	url := fmt.Sprintf("%s://%s:%s/transaction?page=%d", backendProtocol, backendIp, backendPort, page)
	resp, err := http.Get(url)
	if err != nil {
		helper.Logger.Error().Err(err).Msg("Cannot run http get request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helper.Logger.Error().Err(err).Msg("ReadAll error")
	}

	var txs []Transaction
	if err := json.Unmarshal(body, &txs); err != nil {
		helper.Logger.Error().Err(err).Msg("Cannot unmarshal transactions")

	}

	return txs
}

// Create the transaction screen
func NewTransactionScreen(app fyne.App) fyne.CanvasObject {

	var txs = getTransactions(1, app) // Fill txs with the first page of txs
	var txsPerPage = 50               // Default number of txs returned by the backend when querrying the endpoint "/transaction"
	var reachedDataEnd = false
	var threshold = 5

	list := widget.NewList(
		func() int {
			return len(txs)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("date"),
				widget.NewLabel("value"),
				widget.NewLabel("type"),
				widget.NewLabel("name"), // ToDo: use a scroll container for long text
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			hbox := o.(*fyne.Container).Objects
			txDate := hbox[0].(*widget.Label)
			txValue := hbox[1].(*widget.Label)
			txType := hbox[2].(*widget.Label)
			txName := hbox[3].(*widget.Label)

			parsedTxDate, err := time.Parse("2006-01-02 15:04:05", txs[i].Date)
			if err != nil {
				helper.Logger.Error().Err(err).Msgf("Cannot parse date %s", txs[i].Date)
			}

			txType.SetText(txs[i].Transaction_type)
			txDate.SetText(parsedTxDate.Format("2006-01-02"))
			txValue.SetText(fmt.Sprintf("%.2f", txs[i].Value))
			txName.SetText(txs[i].Original_wording)

			// Load new items in the list when the user scrolled near the bottom of the page => infinite scrolling
			// We ask more data from the backend if we only have less than "threshold" txs left to display
			if i > len(txs)-threshold && !reachedDataEnd {
				pageRequested := len(txs)/txsPerPage + 1
				newTxs := getTransactions(pageRequested, app)

				// We have retrieved every transaction if the backend sent less txs than the default number per page
				if len(newTxs) < txsPerPage {
					reachedDataEnd = true
				}
				txs = append(txs, newTxs...)
			}
		})

	return list
}
