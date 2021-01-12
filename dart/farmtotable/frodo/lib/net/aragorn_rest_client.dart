import 'dart:convert';
import 'package:http/http.dart' as http;

import '../models/auction_item.dart';
import '../models/order.dart';
import '../models/item.dart';

class AragornRestClient {
  static const baseRoute = "http://165.22.222.169:8080/api/v1/resources/";

  bool _parseAuctionsResponse(String jsonStr, List<AuctionItem> auctions) {
    Map<String, dynamic> myMap = json.decode(jsonStr);
    if (myMap["status"] < 200 || myMap["status"] >= 300) {
      print("Received error from backend: ${myMap['error_msg']}");
      return false;
    }
    final auctionsMap = myMap["data"]["auctions"];
    if (auctionsMap != null) {
      auctionsMap.forEach((auction) {
        auctions.add(AuctionItem(
            itemDescription: "Something",
            auctionID: auction["id"] as int,
            itemID: auction["item_id"] as String,
            itemName: auction["item_name"] as String,
            itemQty: auction["item_qty"] as int,
            auctionDurationSecs:
            Duration(seconds: auction["auction_duration_secs"]),
            auctionStartTime: DateTime.parse(auction["auction_start_time"]),
            minBid: double.parse(auction["min_bid"].toString()),
            maxBid: double.parse(auction["max_bid"].toString()),
            imageURL: auction["image_url"] as String));
      });
    }
    return true;
  }

  Future<List<AuctionItem>> getAuctions(int startID, int numAuctions) async {
    final route = baseRoute + "auctions/fetch_all";
    try {
      final response = await http.post(route,
          body:
              json.encode({"start_id": startID, "num_auctions": numAuctions}));
      List<AuctionItem> auctions = [];
      if (_parseAuctionsResponse(response.body, auctions)) {
        return auctions;
      } else {
        throw Future.error("Unable to parse response object");
      }
    } catch (error) {
      throw error;
    }
  }

  Item _parseItemResponse(String jsonStr) {
    Map<String, dynamic> myMap = json.decode(jsonStr);
    if (myMap["status"] < 200 || myMap["status"] >= 300) {
      print("Received error from backend: ${myMap['error_msg']}");
      throw Future.error("Failure while fetching item data from backend");
    }
    final itemDeets = myMap["data"];
    return Item(
      itemID: itemDeets['item_id'],
      itemName: itemDeets['item_name'],
      itemDescription: itemDeets['item_description'],
      itemQty: itemDeets['item_qty'],
      itemUnit: itemDeets['item_unit'],
      minBidPrice: double.parse(itemDeets['min_price'].toString()),
      minBidQty: itemDeets['min_bid_qty'],
      maxBidQty: itemDeets['max_bid_qty'],
      imageURL: itemDeets['image_url'],
      auctionDurationSecs:
          Duration(seconds: itemDeets['auction_duration_secs']),
      auctionStartTime: DateTime.parse(itemDeets['auction_start_time']),
    );
  }

  Future<Item> getItem(String itemID) async {
    final route = baseRoute + "items/get";
    try {
      final response = await http.post(route,
          body: json.encode({
            "item_id": itemID,
          }));
      final item = _parseItemResponse(response.body);
      return item;
    } catch (error) {
      throw error;
    }
  }

  void _parseRegisterBidResponse(String jsonStr) {
    Map<String, dynamic> myMap = json.decode(jsonStr);
    if (myMap["status"] < 200 || myMap["status"] >= 300) {
      print("Received error from backend: ${myMap['error_msg']}");
      throw Future.error("${myMap['error_msg']}");
    }
  }

  Future<void> registerBid(String itemID, double bidAmt, int bidQty) async {
    final route = baseRoute + "auctions/register_bid";
    try {
      final response = await http.post(route,
          body: json.encode({
            "item_id": itemID,
            "user_id": "user1",
            "bid_amount": bidAmt,
            "bid_qty": bidQty,
          }));
      _parseRegisterBidResponse(response.body);
    } catch (error) {
      throw error;
    }
  }

  bool _parseOrdersResponse(String jsonStr, List<Order> orders) {
    print(jsonStr);
    Map<String, dynamic> myMap = json.decode(jsonStr);
    if (myMap["status"] < 200 || myMap["status"] >= 300) {
      print("Received error from backend: ${myMap['error_msg']}");
      return false;
    }
    final ordersMap = myMap["data"]["orders"];
    ordersMap.forEach((order) {
      double ip = order["item_price"] as num;
      double dp = double.parse(order["delivery_price"].toString());
      double txp = double.parse(order["tax_price"].toString());
      double ttp = double.parse(order["total_price"].toString());
      orders.add(Order(
          orderID: order["order_id"] as String,
          itemID: order["item_id"] as String,
          itemName: order["item_name"] as String,
          orderedQty: (order["item_qty"] as num) as int,
          itemPrice: ip,
          deliveryPrice: dp,
          taxPrice: txp,
          totalPrice: ttp,
          status: OrderStatus((order["status"] as num) as int),
          imageURL: order["image_url"] as String));
    });
    return true;
  }

  Future<List<Order>> getUserOrders(String userID) async {
    final route = baseRoute + "orders/get_user_orders";
    try {
      final response = await http.post(route,
          body: json.encode({
            "user_id": "user1",
          }));
      List<Order> orders = [];
      _parseOrdersResponse(response.body, orders);
      return orders;
    } finally {
    }
  }

  Order _parseOrderInfo(String jsonStr) {
    print(jsonStr);
    Map<String, dynamic> myMap = json.decode(jsonStr);
    if (myMap["status"] < 200 || myMap["status"] >= 300) {
      print("Received error from backend: ${myMap['error_msg']}");
      throw Future.error("Did not get expected response");
    }
    final orderMap = myMap["data"];
    double ip = orderMap["item_price"] as num;
    double dp = double.parse(orderMap["delivery_price"].toString());
    double txp = double.parse(orderMap["tax_price"].toString());
    double ttp = double.parse(orderMap["total_price"].toString());
    final order = Order(
        orderID: orderMap["order_id"] as String,
        itemID: orderMap["item_id"] as String,
        itemName: orderMap["item_name"] as String,
        orderedQty: orderMap["item_qty"] as num,
        itemPrice: ip,
        deliveryPrice: dp,
        taxPrice: txp,
        totalPrice: ttp,
        status: OrderStatus(orderMap["status"] as int),
        imageURL: orderMap["image_url"] as String);
    return order;
  }

  Future<Order> getUserOrder(String userID, String itemID) async {
    final route = baseRoute + "orders/get_order";
    try {
      final response = await http.post(route,
          body: json.encode({
            "user_id": "user1",
          }));
      return _parseOrderInfo(response.body);
    } catch (error) {
      throw error;
    }
  }

  double _parseGetUserBidInfo(String jsonStr) {
    print(jsonStr);
    Map<String, dynamic> myMap = json.decode(jsonStr);
    if (myMap["status"] < 200 || myMap["status"] >= 300) {
      print("Received error from backend: ${myMap['error_msg']}");
      throw Future.error("Did not get expected response");
    }
    final userBid = myMap["data"] as double;
    return userBid;
  }

  Future<double> getUserBid(String userID, String itemID) async {
    final route = baseRoute + "auctions/get_user_bid";
    try {
      final response = await http.post(route,
          body: json.encode({
            "user_id": "user1",
          }));
      return _parseGetUserBidInfo(response.body);
    } catch (error) {
      throw error;
    }
  }
}
