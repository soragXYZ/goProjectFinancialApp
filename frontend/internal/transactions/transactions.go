package transactions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"freenahiFront/internal/helper"
	"freenahiFront/internal/settings"
	"io"
	"math"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Const which are used in order to increase code readability
// These consts are the name of the columns and used later in a switch case
const (
	pinnedColumn int = iota
	dateColumn
	valueColumn
	typeColumn
	detailsColumn
)

// The struct which is returned by the backend
type Transaction struct {
	Id               int     `json:"id"`
	Pinned           bool    `json:"pinned"`
	Date             string  `json:"date"`
	Value            float32 `json:"value"`
	Transaction_type string  `json:"type"`
	Original_wording string  `json:"original_wording"`
}

// Create the transaction screen
func NewTransactionScreen(app fyne.App, win fyne.Window) fyne.CanvasObject {

	var (
		pinnedLabel          = widget.NewLabel(lang.L("Pinned"))
		dateLabel            = widget.NewLabel(lang.L("Date"))
		valueLabel           = widget.NewLabel(lang.L("Value"))
		typeLabel            = widget.NewLabel(lang.L("Type"))
		detailsLabel         = widget.NewLabel(lang.L("Details"))
		testDetailsLabelSize = widget.NewLabel("CB DEBIT IMMEDIAT UBER EATS").MinSize().Width
	)

	// Fill txs with the first page of txs. The first tx is a special item only used for the table header (no real data)
	txs := []Transaction{{
		Original_wording: "columnHeader",
	}}
	txs = append(txs, getTransactions(1, app)...)
	var txsPerPage = 50 // Default number of txs returned by the backend when querrying the endpoint "/transaction"
	var reachedDataEnd = false
	var threshold = 5

	txList := widget.NewTable(
		func() (int, int) {
			return len(txs), 5 // The number of column to display, ie the number of iota const value (icon, date, value, type, details)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("template"))
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {

			// Clean the cell from the previous value
			o.(*fyne.Container).RemoveAll()

			// If we are on the first row, we set special values and we will pin this row to create a header with it
			if i.Row == 0 {
				switch i.Col {
				case pinnedColumn:
					AddHAligned(o, pinnedLabel)

				case dateColumn:
					AddHAligned(o, dateLabel)

				case valueColumn:
					AddHAligned(o, valueLabel)

				case typeColumn:
					AddHAligned(o, typeLabel)

				case detailsColumn:
					AddHAligned(o, detailsLabel)

				default:
					helper.Logger.Fatal().Msg("Too much column in the grid")
				}
				return
			}

			// Update the cell by adding content according to its "type" (icon, date, value, type, details)
			switch i.Col {
			case pinnedColumn:
				if txs[i.Row].Pinned {
					AddHAligned(o, widget.NewIcon(theme.RadioButtonCheckedIcon()))
				} else {
					AddHAligned(o, widget.NewIcon(theme.RadioButtonIcon()))
				}

			case dateColumn:
				parsedTxDate, err := time.Parse("2006-01-02 15:04:05", txs[i.Row].Date)
				if err != nil {
					helper.Logger.Error().Err(err).Msgf("Cannot parse date %s", txs[i.Row].Date)
				}
				AddHAligned(o, widget.NewLabel(parsedTxDate.Format("2006-01-02")))

			case valueColumn:
				AddHAligned(o, widget.NewLabel(fmt.Sprintf("%.2f", txs[i.Row].Value)))

			case typeColumn:
				AddHAligned(o, widget.NewLabel(txs[i.Row].Transaction_type))

			case detailsColumn:
				scroller := container.NewHScroll(widget.NewLabel(txs[i.Row].Original_wording))
				scroller.SetMinSize(fyne.NewSize(testDetailsLabelSize, scroller.MinSize().Height))
				AddHAligned(o, scroller)

			default:
				helper.Logger.Fatal().Msg("Too much column in the transaction grid")
			}

			// Load new items in the list when the user scrolled near the bottom of the page => infinite scrolling
			// We ask more data from the backend if we only have less than "threshold" txs left to display
			if i.Row > len(txs)-threshold && !reachedDataEnd {
				pageRequested := len(txs)/txsPerPage + 1
				newTxs := getTransactions(pageRequested, app)

				// We have retrieved every transaction if the backend sent less txs than the default number per page
				if len(newTxs) < txsPerPage {
					reachedDataEnd = true
				}
				txs = append(txs, newTxs...)
			}
		},
	)

	txList.OnSelected = func(id widget.TableCellID) {

		// Do nothing if the user selected the column header
		if id.Row == 0 {
			return
		}

		detailsItem := widget.NewEntry()
		detailsItem.SetText(txs[id.Row].Original_wording) // ToDo: add regex and validator

		pinnedItem := widget.NewCheck("", func(value bool) {
			txs[id.Row].Pinned = value
		})
		pinnedItem.Checked = txs[id.Row].Pinned

		items := []*widget.FormItem{
			widget.NewFormItem("Details", detailsItem),
			widget.NewFormItem("Pinned", pinnedItem),
		}

		d := dialog.NewForm("Edit transaction", "Update", "Cancel", items, func(b bool) {
			if !b {
				return
			}
			txs[id.Row].Original_wording = detailsItem.Text // replaced by the user input

			// Refresh the whole row with the new data set by the user
			txList.RefreshItem(widget.TableCellID{Row: id.Row, Col: pinnedColumn})
			txList.RefreshItem(widget.TableCellID{Row: id.Row, Col: dateColumn})
			txList.RefreshItem(widget.TableCellID{Row: id.Row, Col: valueColumn})
			txList.RefreshItem(widget.TableCellID{Row: id.Row, Col: valueColumn})
			txList.RefreshItem(widget.TableCellID{Row: id.Row, Col: detailsColumn})

			// Call the backend to apply the change
			updateTransaction(txs[id.Row], app)
		}, win)

		d.Resize(fyne.NewSize(d.MinSize().Width*2, d.MinSize().Height))
		d.Show()

	}

	// We set the width of the columns, ie the max between the language name header size and actual value
	// For example, the max between "Value" and "-123456123.00", or "Montant" and "-123456123.00" in french
	txList.SetColumnWidth(pinnedColumn, float32(math.Max(
		float64(widget.NewIcon(theme.RadioButtonCheckedIcon()).MinSize().Width),
		float64(pinnedLabel.MinSize().Width))),
	)
	txList.SetColumnWidth(dateColumn, float32(math.Max(
		float64(widget.NewLabel("XXXX-YY-ZZ").MinSize().Width),
		float64(dateLabel.MinSize().Width))),
	)
	txList.SetColumnWidth(valueColumn, float32(math.Max(
		float64(widget.NewLabel("-123456123.00").MinSize().Width),
		float64(valueLabel.MinSize().Width))),
	)
	txList.SetColumnWidth(typeColumn, float32(math.Max(
		float64(widget.NewLabel("loan_repayment").MinSize().Width),
		float64(typeLabel.MinSize().Width))),
	)
	txList.SetColumnWidth(detailsColumn, float32(math.Max(
		float64(testDetailsLabelSize),
		float64(detailsLabel.MinSize().Width))),
	)

	txList.StickyRowCount = 1 // Basically, we are setting a table header because the first row contains special data

	return txList
}

// Center align the objectToAdd by adding 2 spacers. To be used with an horizontal box
func AddHAligned(object fyne.CanvasObject, objectToAdd fyne.CanvasObject) {
	object.(*fyne.Container).Add(layout.NewSpacer())
	object.(*fyne.Container).Add(objectToAdd)
	object.(*fyne.Container).Add(layout.NewSpacer())
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
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helper.Logger.Error().Err(err).Msg("ReadAll error")
		return nil
	}

	var txs []Transaction
	if err := json.Unmarshal(body, &txs); err != nil {
		helper.Logger.Error().Err(err).Msg("Cannot unmarshal transactions")
		return nil

	}

	return txs
}

// Call the backend endpoint "/transaction" and update the specified tx
func updateTransaction(tx Transaction, app fyne.App) {

	backendIp := app.Preferences().StringWithFallback(settings.PreferenceBackendIP, settings.BackendIPDefault)
	backendProtocol := app.Preferences().StringWithFallback(settings.PreferenceBackendProtocol, settings.BackendProtocolDefault)
	backendPort := app.Preferences().StringWithFallback(settings.PreferenceBackendPort, settings.BackendPortDefault)

	url := fmt.Sprintf("%s://%s:%s/transaction/%d", backendProtocol, backendIp, backendPort, tx.Id)

	jsonBody, err := json.Marshal(tx)
	if err != nil {
		helper.Logger.Error().Err(err).Msg("Cannot marshal tx")
		return
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		helper.Logger.Error().Err(err).Msg("Cannot create new request")
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		helper.Logger.Error().Err(err).Msg("Cannot run http put request")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		helper.Logger.Error().Msg(resp.Status)
		return
	}
}
