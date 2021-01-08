import 'package:flutter/foundation.dart';

class Item {
  final String itemID;
  final String itemName;
  final String itemDescription;
  final int    itemQty;
  final String imageURL;
  final double minBidPrice;
  final int minBidQty;
  final int maxBidQty;
  final DateTime auctionStartTime;
  final Duration auctionDurationSecs;
  final String itemUnit;

  Item({
    @required this.itemDescription,
    @required this.itemID,
    @required this.itemName,
    @required this.itemQty,
    @required this.imageURL,
    @required this.minBidPrice,
    @required this.minBidQty,
    @required this.maxBidQty,
    @required this.auctionStartTime,
    @required this.auctionDurationSecs,
    @required this.itemUnit,
  });
}
