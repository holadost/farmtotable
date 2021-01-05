import '../models/auction_item.dart';

var DUMMY_AUCTIONS = <AuctionItem>[
  AuctionItem(
      auctionID: 1,
      itemID: "Item1",
      itemName: "Rice grains",
      itemQty: 300,
      auctionDurationSecs: Duration(seconds: 3600),
      auctionStartTime: DateTime.now(),
      minBid: 10.00),
  AuctionItem(
      auctionID: 2,
      itemID: "Item2",
      itemName: "Whole Wheat",
      itemQty: 200,
      auctionDurationSecs: Duration(seconds: 3600),
      auctionStartTime: DateTime.now(),
      minBid: 15.00),
  AuctionItem(
      auctionID: 3,
      itemID: "Item3",
      itemName: "Peas and beans",
      itemQty: 300,
      auctionDurationSecs: Duration(seconds: 3600),
      auctionStartTime: DateTime.now(),
      minBid: 22.00),
  AuctionItem(
      auctionID: 4,
      itemID: "Item4",
      itemName: "Carrots",
      itemQty: 400,
      auctionDurationSecs: Duration(seconds: 3600),
      auctionStartTime: DateTime.now(),
      minBid: 14.00),
  AuctionItem(
      auctionID: 5,
      itemID: "Item5",
      itemName: "Others",
      itemQty: 500,
      auctionDurationSecs: Duration(seconds: 3600),
      auctionStartTime: DateTime.now(),
      minBid: 20.00),
];