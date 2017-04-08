package helpers

import (
	"github.com/coldbrewcloud/go-shippo/client"
	tcm "github.com/jakelong95/TownCenter/models"
	"github.com/coldbrewcloud/go-shippo/models"
)

/*TODO: Pass given user and roaster information here*/
func createShipment(c *client.Client, user models.User, roaster models.Roaster) {
	//Roaster address
	addressFromInput := &models.AddressInput{
		ObjectPurpose: models.ObjectPurposePurchase,
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
	addressToInput := &models.AddressInput{
		ObjectPurpose: models.ObjectPurposePurchase,
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
	//TODO: dynamically generate sizing based on size
	parcelInput := &models.ParcelInput{
		Length: "5",
		Width: "5",
		Height: "5",
		DistanceUnit: models.DistanceUnitInch,
		Weight: "5",
		MassUnit: models.MassUnitPound,
	}
	parcel, err := c.CreateParcel(parcelInput)
	if err != nil {
		panic(err)
	}

	shipmentInput := &models.ShipmentInput{
		ObjectPurpose: models.ObjectPurposePurchase,
		AddressFrom: addressFrom.ObjectID,
		AddressTo: addressTo,
		Parcel: parcel.ObjectID
		Async: false,
	}
	shipment, err := c.CreateShipment(shipmentInput)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Shipment:\n%s\n", dump(shipment))

	return shipment
}

func purchaseShippingLabel(c *client.Client, shipment *models.Shipment) *models.Transaction{
	transactionInput := &models.TransactionInput{
		Rate: shipment.RateList[0].ObjectID, 	//TODO pick the cheapest option for Rate and pick file
		LabelFileType: models.LabelFileTypePDF,
		Async: false,
	}
	transaction, err := purchaseShippingLabel(transactionInput)
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
