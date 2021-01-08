import 'dart:convert';
import 'package:http/http.dart' as http;

import '../models/auction_item.dart';

class RestApiClient {
  static const baseRoute = "http://165.22.222.169:8080/api/v1/resources/";

  bool parseAuctionsResponse(String jsonStr, List<AuctionItem> auctions) {
    Map<String, dynamic> myMap = json.decode(jsonStr);
    if (myMap["status"] < 200 || myMap["status"] >= 300) {
      print("Received error from backend: ${myMap['error_msg']}");
      return false;
    }
    final auctionsMap = myMap["data"]["auctions"];
    auctionsMap.forEach((auction) {
      auctions.add(
          AuctionItem(
              itemDescription: "Something",
              auctionID: auction["id"] as int,
              itemID: auction["item_id"] as String,
              itemName: auction["item_name"] as String,
              itemQty: auction["item_qty"] as int,
              auctionDurationSecs: Duration(seconds: auction["auction_duration_secs"]),
              auctionStartTime: DateTime.parse(auction["auction_start_time"]),
              minBid: double.parse(auction["min_bid"].toString()),
              maxBid: double.parse(auction["min_bid"].toString()),
              imageURL: auction["image_url"] as String));
    });
    return true;
  }

  Future<List<AuctionItem>> getAuctions(int startID, int numAuctions) async {
    final route = baseRoute + "auctions/fetch_all";
    try {
      final response = await http.post(route, body: json.encode({
        "start_id": startID, "num_auctions": numAuctions
      }));
      List<AuctionItem> auctions = [];
      if (parseAuctionsResponse(response.body, auctions)) {
        return auctions;
      } else {
        throw Future.error("Unable to parse response object");
      }
    } catch (error) {
      throw error;
    }
  }

}