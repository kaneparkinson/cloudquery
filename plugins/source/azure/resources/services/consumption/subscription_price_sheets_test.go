package consumption

import (
	"encoding/json"
	"net/http"

	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/consumption/armconsumption"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/v3/faker"
	"github.com/gorilla/mux"
)

func createSubscriptionPriceSheets(router *mux.Router) error {
	var resp armconsumption.PriceSheetClientGetResponse
	if err := faker.FakeObject(&resp); err != nil {
		return err
	}
	resp.PriceSheetResult.Properties.NextLink = nil

	router.HandleFunc("/subscriptions/{subscriptionId}/providers/Microsoft.Consumption/pricesheets/default", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(&resp)
		if err != nil {
			http.Error(w, "unable to marshal request: "+err.Error(), http.StatusBadRequest)
			return
		}
		if _, err := w.Write(b); err != nil {
			http.Error(w, "failed to write", http.StatusBadRequest)
			return
		}
	})

	return nil
}

func TestSubscriptionPriceSheets(t *testing.T) {
	client.MockTestHelper(t, SubscriptionPriceSheets(), createSubscriptionPriceSheets)
}
