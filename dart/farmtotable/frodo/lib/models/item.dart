import 'package:flutter/foundation.dart';

class Item {
  final String itemID;
  final String itemName;
  final String itemDescription;
  final int    itemQty;
  final String imageURL;

  Item({
    @required this.itemDescription,
    @required this.itemID,
    @required this.itemName,
    @required this.itemQty,
    @required this.imageURL,
  });
}
