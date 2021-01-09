import 'package:flutter/foundation.dart';

class Order {
  final String orderID;
  final String itemID;
  final String itemName;
  final int orderedQty;
  final double price;
  final String imageURL;
  final int status;

  Order({
    @required this.orderID,
    @required this.itemID,
    @required this.itemName,
    @required this.orderedQty,
    @required this.price,
    @required this.imageURL,
    @required this.status,
  });
}
