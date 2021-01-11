import 'package:flutter/foundation.dart';

class OrderStatus {
  final int _value;

  const OrderStatus(this._value);

  int toInt() {
    return _value;
  }
  @override
  String toString() {
    if (_value == 0) {
      return "Payment Pending";
    } else if (_value == 1) {
      return "Pending delivery";
    } else if (_value == 2) {
      return "Delivery complete";
    } else if (_value == 3) {
      return "Cancelled";
    } else {
      throw new FormatException("Invalid order status");
    }
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
  final double price;
  final String imageURL;
  final OrderStatus status;

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
