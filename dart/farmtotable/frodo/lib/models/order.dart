import 'package:flutter/foundation.dart';

class Order {
  final String orderID;
  final String itemID;
  final String itemName;
  final String itemDescription;
  final String itemImageURL;
  final int orderedQty;
  final double price;

  Order({
    @required this.itemDescription,
    @required this.orderID,
    @required this.itemID,
    @required this.itemName,
    @required this.orderedQty,
    @required this.price,
    @required this.itemImageURL,
  });
}
