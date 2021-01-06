import 'package:flutter/foundation.dart';

class AuctionItem {
  final int auctionID;
  final String itemID;
  final String itemName;
  final String itemDescription;
  final int itemQty;
  final DateTime auctionStartTime;
  final Duration auctionDurationSecs;
  final double minBid;
  final double maxBid;
  final String imageURL;

  AuctionItem({
    @required this.itemDescription,
    @required this.auctionID,
    @required this.itemID,
    @required this.itemName,
    @required this.itemQty,
    @required this.auctionDurationSecs,
    @required this.auctionStartTime,
    @required this.minBid,
    @required this.maxBid,
    @required this.imageURL,
  });
}
