package gateways

import (
	"fmt"
	"strconv"

	"github.com/coldbrewcloud/go-shippo"
	"github.com/coldbrewcloud/go-shippo/client"
	shipm "github.com/coldbrewcloud/go-shippo/models"
	"github.com/ghmeier/bloodlines/config"
	tcm "github.com/jakelong95/TownCenter/models"
	"github.com/lcollin/warehouse/models"
)

type Shippo interface {
	CreateShipment(*tcm.User, *tcm.Roaster, *models.Dimensions) (*shipm.Shipment, error)
	PurchaseShippingLabel(*shipm.Shipment) (*shipm.Transaction, error)
	CreateAddress(string, string, string, string, string, string, string, string) (*shipm.Address, error)
}

type ship struct {
	c *client.Client
}

func NewShippo(cfg config.Shippo) Shippo {
	return &ship{
		c: shippo.NewClient(cfg.Token),
	}
}

/*CreateShipment creates a shipment object, consisting of address from, address to, and parcel*/
func (s *ship) CreateShipment(user *tcm.User, roaster *tcm.Roaster, dimensions *models.Dimensions) (*shipm.Shipment, error) {
	addressFrom, err := s.CreateAddress(
		roaster.Name,
		roaster.AddressLine1,
		roaster.AddressCity,
		roaster.AddressState,
		roaster.AddressZip,
		roaster.AddressCountry,
		roaster.Phone,
		roaster.Email,
	)
	if err != nil {
		return nil, err
	}

	addressTo, err := s.CreateAddress(
		user.FirstName+" "+user.LastName,
		user.AddressLine1,
		user.AddressCity,
		user.AddressState,
		user.AddressZip,
		user.AddressCountry,
		user.Phone,
		user.Email,
	)
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
	parcel, err := s.c.CreateParcel(parcelInput)
	if err != nil {
		return nil, err
	}

	shipmentInput := &shipm.ShipmentInput{
		AddressFrom: addressFrom.ObjectID,
		AddressTo:   addressTo.ObjectID,
		Parcels:     []string{parcel.ObjectID},
		Async:       false,
	}
	shipment, err := s.c.CreateShipment(shipmentInput)
	if err != nil {
		return nil, err
	}
	return shipment, nil
}

func (s *ship) CreateAddress(name, street, city, state, zip, country, phone, email string) (*shipm.Address, error) {
	input := &shipm.AddressInput{
		Name:     name,
		Street1:  street,
		City:     city,
		State:    state,
		Zip:      zip,
		Country:  country,
		Phone:    phone,
		Email:    email,
		Validate: true,
	}
	address, err := s.c.CreateAddress(input)
	if err != nil {
		return nil, err
	}
	if !address.ValidationResults.IsValid {
		return nil, fmt.Errorf("Invalid address: %s", address.ValidationResults.Messages[0].Text)
	}

	return address, nil
}

/*PurchaseShippingLabel creates a transaction object*/
func (s *ship) PurchaseShippingLabel(shipment *shipm.Shipment) (*shipm.Transaction, error) {
	transactionInput := &shipm.TransactionInput{
		Rate:          shipment.Rates[0].ObjectID, //TODO: pick the cheapest option for Rate and pick file
		LabelFileType: shipm.LabelFileTypePDF,     //TODO: offer option of which file type?
		Async:         false,
	}
	transaction, err := s.c.PurchaseShippingLabel(transactionInput)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
