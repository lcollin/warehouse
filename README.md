# warehouse
Go service handling Expresso inventory

## API

### Item

#### `POST /api/item` creates and adds a new item to the database.

Example:

*Request:*
```
POST localhost:8080/api/item
{
    "roasterID": "dd82cc65-d79d-11e6-9d4c-0242ac120004",
    "name": "Spring Blend",
    "coffeeType": "dark roast",
    "inStockBags": 50,
    "providerPrice": 4.00,
    "consumerPrice": 9.00,
    "ozInBag": 12.50
}
```

*Response:*
```
{
  "data": {
    "ID": "deb0a02f-f159-11e6-9563-acbc32977aaf",
    "RoasterID": "dd82cc65-d79d-11e6-9d4c-0242ac120004",
    "Name": "Spring Blend",
    "Picture": "",
    "CoffeeType": "dark roast",
    "InStockBags": 50,
    "ProviderPrice": 4,
    "ConsumerPrice": 9,
    "OzInBag": 12.5
  },
  "success": true
}
```

#### `GET /api/item?offset=0&limit=20` returns up to `limit` items starting from `offset` when ordered by itemID

Example:

*Request:*
```
GET localhost:8080/api/item?offset=0&limit=20
```

*Response:*
```
{
  "data": [
    {
      "ID": "deb0a02f-f159-11e6-9563-acbc32977aaf",
      "RoasterID": "dd82cc65-d79d-11e6-9d4c-0242ac120004",
      "Name": "Spring Blend",
      "Picture": "",
      "CoffeeType": "dark roast",
      "InStockBags": 50,
      "ProviderPrice": 4,
      "ConsumerPrice": 9,
      "OzInBag": 12.5
    }
  ],
  "success": true
}
```

#### `GET /api/item/:itemID` returns the item with the given itemID

Example:

*Request:*
```
GET localhost:8080/api/item/deb0a02f-f159-11e6-9563-acbc32977aaf
```

*Response:*
```
{
  "data": {
    "ID": "deb0a02f-f159-11e6-9563-acbc32977aaf",
    "RoasterID": "dd82cc65-d79d-11e6-9d4c-0242ac120004",
    "Name": "Spring Blend",
    "Picture": "",
    "CoffeeType": "dark roast",
    "InStockBags": 50,
    "ProviderPrice": 4,
    "ConsumerPrice": 9,
    "OzInBag": 12.5
  },
  "success": true
}
```
#### `PUT /api/item/:itemID` updates the item with the given itemID to match the provided data. This just overrides values, so anything not present in the request will be set to NULL

Example:

*Request:*
```
PUT localhost:8080/api/item/1a7c63af-f15c-11e6-bad7-acbc32977aaf
{
    "id": "1a7c63af-f15c-11e6-bad7-acbc32977aaf",
    "roasterID": "dd82cc65-d79d-11e6-9d4c-0242ac120004",
    "name": "Spring Blend",
    "pictureURL": "imgur.com",
    "coffeeType": "dark roast",
    "inStockBags": 50,
    "providerPrice": 4.00,
    "consumerPrice": 9.00,
    "ozInBag": 12.50
}
```
*Response:*
```
{
  "data": {
    "ID": "1a7c63af-f15c-11e6-bad7-acbc32977aaf",
    "RoasterID": "dd82cc65-d79d-11e6-9d4c-0242ac120004",
    "Name": "Spring Blend",
    "PictureURL": "imgur.com",
    "CoffeeType": "dark roast",
    "InStockBags": 50,
    "ProviderPrice": 4,
    "ConsumerPrice": 9,
    "OzInBag": 12.5
  },
  "success": true
}
```
#### `DELETE /api/item/:itemID` removes the item

Example:

*Request:*
```
DELETE localhost:8080/api/item/1a7c63af-f15c-11e6-bad7-acbc32977aaf
```

*Response:*
```
{
  "data": null,
  "success": true
}
```

### Order

#### `POST /api/order` creates and adds a new order to the database.

Example:

*Request:*
```
POST localhost:8080/api/order
{
    "userID": "dd82cc65-d79d-11e6-9d4c-0242ac120006",
    "subscriptionID": "dd82cc65-d79d-11e6-9d4c-0242ac120005"
}
```

*Response:*
```
{
  "data": {
    "ID": "9772d7ea-f15e-11e6-bad7-acbc32977aaf",
    "UserID": "dd82cc65-d79d-11e6-9d4c-0242ac120006",
    "SubscriptionID": "dd82cc65-d79d-11e6-9d4c-0242ac120005",
    "RequestDate": "2017-02-12T14:05:25.148465376-06:00",
    "ShipDate": "2020-01-01T01:00:00Z"
  },
  "success": true
}
```

#### `GET /api/order?offset=0&limit=20` returns up to `limit` orders starting from `offset` when ordered by orderID
#### `GET /api/roaster/order/:roasterId?offset=0&limit=20` returns up to `limit` orders starting from `offset` that belong to the roaster with the given id
#### `GET /api/user/order/:userId?offset=0&limit=20` returns up to `limit` orders starting from `offset` that belong to the user with the given id

Example:

*Request:*
```
GET localhost:8080/api/order?offset=0&limit=20
```

*Response:*
```
{
  "data": [
    {
      "ID": "9772d7ea-f15e-11e6-bad7-acbc32977aaf",
      "UserID": "dd82cc65-d79d-11e6-9d4c-0242ac120006",
      "SubscriptionID": "dd82cc65-d79d-11e6-9d4c-0242ac120005",
      "RequestDate": "2017-02-12T20:05:25Z",
      "ShipDate": "2020-01-01T01:00:00Z"
    }
  ],
  "success": true
}
```

#### `GET /api/order/:orderID` returns the order with the given orderID

Example:

*Request:*
```
GET localhost:8080/api/order/9772d7ea-f15e-11e6-bad7-acbc32977aaf
```

*Response:*
```
{
  "data": {
    "ID": "9772d7ea-f15e-11e6-bad7-acbc32977aaf",
    "UserID": "dd82cc65-d79d-11e6-9d4c-0242ac120006",
    "SubscriptionID": "dd82cc65-d79d-11e6-9d4c-0242ac120005",
    "RequestDate": "2017-02-12T20:05:25Z",
    "ShipDate": "2020-01-01T01:00:00Z"
  },
  "success": true
}
```
#### `PUT /api/order/:orderID` updates the order with the given orderID to match the provided data. This just overrides values, so anything not present in the request will be set to NULL

Example:

*Request:*
```
PUT localhost:8080/api/order/9772d7ea-f15e-11e6-bad7-acbc32977aaf
{
    "id": "9772d7ea-f15e-11e6-bad7-acbc32977aaf",
    "userID": "dd82cc65-d79d-11e6-9d4c-0242ac120006",
    "subscriptionID": "dd82cc65-d79d-11e6-9d4c-0242ac120005",
    "requestDate": "2017-02-12T20:05:25Z",
    "shipDate": "2021-01-01T01:00:00Z"
}
```
*Response:*
```
{
  "data": {
    "ID": "9772d7ea-f15e-11e6-bad7-acbc32977aaf",
    "UserID": "dd82cc65-d79d-11e6-9d4c-0242ac120006",
    "SubscriptionID": "dd82cc65-d79d-11e6-9d4c-0242ac120005",
    "RequestDate": "2017-02-12T20:05:25Z",
    "ShipDate": "2021-01-01T01:00:00Z"
  },
  "success": true
}
```
#### `DELETE /api/order/:orderID` removes the order

Example:

*Request:*
```
DELETE localhost:8080/api/order/9772d7ea-f15e-11e6-bad7-acbc32977aaf
```

*Response:*
```
{
  "data": null,
  "success": true
}
```
#### `POST api/label` creates a shipping label for the specified order

Example:

*Request:*
```
POST localhost:8080/api/label
{
	"orderId":"025138a7-2231-11e7-a6a0-0242ac13000b",
	"userId": "69c68317-f7d3-11e6-b525-0242ac13000b",
	"roasterId": "247a0ba0-1c91-11e7-938e-0242ac13000a",
	"quantity": 1,
	"ozInBag": 5.5,
	"length": 5,
	"width": 5,
	"height": 5,
	"distanceUnit": "in",
	"massUnit": "oz"
}

```

*Response:*
```
{
  "data": {
    "id": "025138a7-2231-11e7-a6a0-0242ac13000b",
    "userId": "69c68317-f7d3-11e6-b525-0242ac13000b",
    "subscriptionId": "ff6e2b6d-2230-11e7-b826-0242ac13000c",
    "requestDate": "2017-04-18T06:41:09Z",
    "shipDate": "0001-01-01T00:00:00Z",
    "quantity": 10,
    "status": "SHIPPED",
    "labelUrl": "https://shippo-delivery-east.s3.amazonaws.com/29b9eed8ece744258d91eb944f817bf1.pdf?Signature=QWdjN7l2HSTIz7XAk%2BiZjboLA5s%3D&Expires=1524261228&AWSAccessKeyId=AKIAJGLCC5MYLLWIG42A",
    "trackingUrl": "https://tools.usps.com/go/TrackConfirmAction_input?origTrackNum=9270190164917307321057"
  },
  "success": true
}
```

### SubOrder

#### `POST /api/suborder` creates and adds a new suborder to the database.

Example:

*Request:*
```
POST localhost:8080/api/suborder
{
    "orderID": "dd82cc65-d79d-11e6-9d4c-0242ac120004",
    "itemID": "dd82cc65-d79d-11e6-9d4c-0242ac120005",
    "quantity": 12
}
```

*Response:*
```
{
  "data": {
    "ID": "b39dd4a7-f160-11e6-bad7-acbc32977aaf",
    "OrderID": "dd82cc65-d79d-11e6-9d4c-0242ac120004",
    "ItemID": "dd82cc65-d79d-11e6-9d4c-0242ac120005",
    "Quantity": 12
  },
  "success": true
}
```

#### `GET /api/suborder?offset=0&limit=20` returns up to `limit` suborders starting from `offset` when subordered by suborderID

Example:

*Request:*
```
GET localhost:8080/api/suborder?offset=0&limit=20
```

*Response:*
```
{
  "data": [
    {
      "ID": "b39dd4a7-f160-11e6-bad7-acbc32977aaf",
      "OrderID": "dd82cc65-d79d-11e6-9d4c-0242ac120004",
      "ItemID": "dd82cc65-d79d-11e6-9d4c-0242ac120005",
      "Quantity": 12
    }
  ],
  "success": true
}
```

#### `GET /api/suborder/:suborderID` returns the suborder with the given suborderID

Example:

*Request:*
```
GET localhost:8080/api/suborder/b39dd4a7-f160-11e6-bad7-acbc32977aaf
```

*Response:*
```
{
  "data": {
    "ID": "b39dd4a7-f160-11e6-bad7-acbc32977aaf",
    "OrderID": "dd82cc65-d79d-11e6-9d4c-0242ac120004",
    "ItemID": "dd82cc65-d79d-11e6-9d4c-0242ac120005",
    "Quantity": 12
  },
  "success": true
}
```
#### `PUT /api/suborder/:suborderID` updates the suborder with the given suborderID to match the provided data. This just overrides values, so anything not present in the request will be set to NULL

Example:

*Request:*
```
PUT localhost:8080/api/suborder/9772d7ea-f15e-11e6-bad7-acbc32977aaf
{
    "id": "b39dd4a7-f160-11e6-bad7-acbc32977aaf",
    "orderID": "dd82cc65-d79d-11e6-9d4c-0242ac120004",
    "itemID": "dd82cc65-d79d-11e6-9d4c-0242ac120005",
    "quantity": 10
}
```
*Response:*
```
{
  "data": {
    "ID": "b39dd4a7-f160-11e6-bad7-acbc32977aaf",
    "OrderID": "dd82cc65-d79d-11e6-9d4c-0242ac120004",
    "ItemID": "dd82cc65-d79d-11e6-9d4c-0242ac120005",
    "Quantity": 10
  },
  "success": true
}
```
#### `DELETE /api/suborder/:suborderID` removes the suborder

Example:

*Request:*
```
DELETE localhost:8080/api/suborder/b39dd4a7-f160-11e6-bad7-acbc32977aaf
```

*Response:*
```
{
  "data": null,
  "success": true
}
```
