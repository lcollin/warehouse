package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/coldbrewcloud/go-shippo/client"
	shipm "github.com/coldbrewcloud/go-shippo/models"
	tcm "github.com/jakelong95/TownCenter/models"
	"github.com/lcollin/warehouse/models"
	"strconv"
)

/*CreateShipment creates a shipment object, consisting of address from, address to, and parcel*/
func CreateShipment(c *client.Client, user *tcm.User, roaster *tcm.Roaster, dimensions *models.Dimensions) *shipm.Shipment {
	//Roaster address
	addressFromInput := &shipm.AddressInput{
		ObjectPurpose: shipm.ObjectPurposePurchase,
		Name:          roaster.Name,
		Street1:       roaster.AddressLine1,
		City:          roaster.AddressCity,
		State:         roaster.AddressState,
		Zip:           roaster.AddressZip,
		Country:       roaster.AddressCountry,
		Phone:         roaster.Phone,
		Email:         roaster.Email,
	}
	addressFrom, err := c.CreateAddress(addressFromInput)
	if err != nil {
		panic(err)
	}
	//Customer address
	addressToInput := &shipm.AddressInput{
		ObjectPurpose: shipm.ObjectPurposePurchase,
		Name:          user.FirstName + " " + user.LastName,
		Street1:       user.AddressLine1,
		City:          user.AddressCity,
		State:         user.AddressState,
		Zip:           user.AddressZip,
		Country:       user.AddressCountry,
		Phone:         user.Phone,
		Email:         user.Email,
	}
	addressTo, err := c.CreateAddress(addressToInput)
	if err != nil {
		panic(err)
	}
	//TODO: use quantity to order multiple parcels?
	parcelInput := &shipm.ParcelInput{
		Length:       strconv.FormatFloat(dimensions.Length, 'f', 2, 64),
		Width:        strconv.FormatFloat(dimensions.Width, 'f', 2, 64),
		Height:       strconv.FormatFloat(dimensions.Height, 'f', 2, 64),
		DistanceUnit: shipm.DistanceUnitInch,
		Weight:       strconv.FormatFloat(dimensions.OzInBag, 'f', 2, 64),
		MassUnit:     shipm.MassUnitOunce,
	}
	parcel, err := c.CreateParcel(parcelInput)
	if err != nil {
		panic(err)
	}

	shipmentInput := &shipm.ShipmentInput{
		ObjectPurpose: shipm.ObjectPurposePurchase,
		AddressFrom:   addressFrom.ObjectID,
		AddressTo:     addressTo.ObjectID,
		Parcel:        parcel.ObjectID,
		Async:         false,
	}
	shipment, err := c.CreateShipment(shipmentInput)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Shipment:\n%s\n", dump(shipment))

	return shipment
}

/*PurchaseShippingLabel*/
func PurchaseShippingLabel(c *client.Client, shipment *shipm.Shipment) *shipm.Transaction {
	transactionInput := &shipm.TransactionInput{
		Rate:          shipment.RatesList[0].ObjectID, //TODO pick the cheapest option for Rate and pick file
		LabelFileType: shipm.LabelFileTypePDF,
		Async:         false,
	}
	transaction, err := c.PurchaseShippingLabel(transactionInput)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transaction:\n%s\n", dump(transaction))
	return transaction
}

func dump(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(data)
}
