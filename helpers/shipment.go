package helpers

import (
	"encoding/json"
	tcm "github.com/jakelong95/TownCenter/models"
	"github.com/lcollin/warehouse/models"
	"github.com/yuderekyu/go-shippo/client"
	shipm "github.com/yuderekyu/go-shippo/models"
	"strconv"
)

/*CreateShipment creates a shipment object, consisting of address from, address to, and parcel*/
func CreateShipment(c *client.Client, user *tcm.User, roaster *tcm.Roaster, dimensions *models.Dimensions) (*shipm.Shipment, error) {
	addressFromInput := &shipm.AddressInput{
		Name:    roaster.Name,
		Street1: roaster.AddressLine1,
		City:    roaster.AddressCity,
		State:   roaster.AddressState,
		Zip:     roaster.AddressZip,
		Country: roaster.AddressCountry,
		Phone:   roaster.Phone,
		Email:   roaster.Email,
	}
	addressFrom, err := c.CreateAddress(addressFromInput)
	if err != nil {
		return nil, err
	}

	addressToInput := &shipm.AddressInput{
		Name:    user.FirstName + " " + user.LastName,
		Street1: user.AddressLine1,
		City:    user.AddressCity,
		State:   user.AddressState,
		Zip:     user.AddressZip,
		Country: user.AddressCountry,
		Phone:   user.Phone,
		Email:   user.Email,
	}
	addressTo, err := c.CreateAddress(addressToInput)
	if err != nil {
		return nil, err
	}

	parcelInput := &shipm.ParcelInput{
		Length:       strconv.FormatFloat(dimensions.Length, 'f', 2, 64),
		Width:        strconv.FormatFloat(dimensions.Width, 'f', 2, 64),
		Height:       strconv.FormatFloat(dimensions.Height, 'f', 2, 64),
		DistanceUnit: dimensions.DistanceUnit,
		Weight:       strconv.FormatFloat(dimensions.OzInBag, 'f', 2, 64),
		MassUnit:     dimensions.MassUnit,
	}
	parcel, err := c.CreateParcel(parcelInput)
	if err != nil {
		return nil, err
	}

	shipmentInput := &shipm.ShipmentInput{
		AddressFrom: addressFrom.ObjectID,
		AddressTo:   addressTo.ObjectID,
		Parcels:     []string{parcel.ObjectID},
		Async:       false,
	}
	shipment, err := c.CreateShipment(shipmentInput)
	if err != nil {
		return nil, err
	}

	return shipment, nil
}

/*PurchaseShippingLabel creates a transaction object*/
func PurchaseShippingLabel(c *client.Client, shipment *shipm.Shipment) (*shipm.Transaction, error) {
	transactionInput := &shipm.TransactionInput{
		Rate:          shipment.Rates[0].ObjectID, //TODO: pick the cheapest option for Rate and pick file
		LabelFileType: shipm.LabelFileTypePDF,     //TODO: offer option of which file type?
		Async:         false,
	}
	transaction, err := c.PurchaseShippingLabel(transactionInput)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func dump(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(data)
}
