import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../util/constants.dart';
import '../models/order.dart';

class OrderScreen extends StatelessWidget {
  static const routeName = '/order-screen';

  void _handlePayment() {

  }

  @override
  Widget build(BuildContext context) {
    final auctionItem = ModalRoute.of(context).settings.arguments as Order;
    final appBar = AppBar(
      backgroundColor: PrimaryColor,
      title: Text(auctionItem.itemName),
      actions: [
        IconButton(onPressed: _handlePayment, icon: Icon(Icons.shopping_cart),)
      ],
    );
    return Scaffold(
      appBar: appBar,
      body: Container(),
    );
  }
}
