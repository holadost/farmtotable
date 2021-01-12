import 'package:flutter/foundation.dart';

class OrderStatus {
  final int value;

  const OrderStatus(this.value);

  int toInt() {
    return value;
  }

  @override
  String toString() {
    if (value == 0) {
      return "Payment Pending";
    } else if (value == 1) {
      return "Pending delivery";
    } else if (value == 2) {
      return "Delivery complete";
    } else if (value == 3) {
      return "Cancelled";
    } else {
      throw new FormatException("Invalid order status");
    }
  }

  @override
  bool operator ==(other) {
    if (value == other.value) {
      return true;
    }
    return false;
  }

  static const OrderStatus KOrderPaymentPending = OrderStatus(0);
  static const OrderStatus KOrderDeliveryPending = OrderStatus(1);
  static const OrderStatus KOrderComplete = OrderStatus(2);
  static const OrderStatus KOrderCancelled = OrderStatus(3);
}

class Order {
  final String orderID;
  final String itemID;
  final String itemName;
  final int orderedQty;
  final double itemPrice;
  final String imageURL;
  final double deliveryPrice;
  final double taxPrice;
  final double totalPrice;
  final OrderStatus status;

  Order({
    @required this.orderID,
    @required this.itemID,
    @required this.itemName,
    @required this.orderedQty,
    @required this.itemPrice,
    @required this.imageURL,
    @required this.status,
    @required this.deliveryPrice,
    @required this.taxPrice,
    @required this.totalPrice,
  });
}
