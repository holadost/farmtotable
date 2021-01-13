import 'package:flutter/material.dart';

import '../screens/item_screen.dart';
import '../util/constants.dart';
import '../models/order.dart';

class OrderScreen extends StatelessWidget {
  static const routeName = '/order-screen';

  void _handlePayment() {}
  void _gotoItemScreen(BuildContext context, String itemID) {
    Navigator.of(context).pushNamed(
        ItemScreen.routeName,
        arguments: {'item_id': itemID, 'show_bid_button': false});
  }

  Widget _buildBody(BuildContext context, Order order) {
    final borderSide =
        BorderSide(color: Colors.grey, width: 2.0, style: BorderStyle.solid);
    var statusColor;
    if (order.status == OrderStatus.KOrderPaymentPending) {
      statusColor = Colors.red;
    } else if (order.status == OrderStatus.KOrderDeliveryPending) {
      statusColor = Colors.lightGreen;
    } else if (order.status == OrderStatus.KOrderComplete) {
      statusColor = Colors.green;
    } else {
      statusColor = Colors.grey;
    }

    final orderContents = Padding(
        padding: const EdgeInsets.all(20.0),
        child: Table(
          border: TableBorder(
              top: borderSide,
              bottom: borderSide,
              right: borderSide,
              left: borderSide),
          children: [
            TableRow(children: [
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "Order ID:",
                  textAlign: TextAlign.left,
                  style: TextStyle(fontSize: 16),
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "${order.orderID}",
                  textAlign: TextAlign.left,
                  style: TextStyle(fontSize: 16),
                ),
              ),
            ]),
            TableRow(children: [
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "Quantity:",
                  textAlign: TextAlign.left,
                  style: TextStyle(fontSize: 16),
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "${order.orderedQty}",
                  textAlign: TextAlign.left,
                  style: TextStyle(fontSize: 16),
                ),
              ),
            ]),
            TableRow(children: [
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "Item Price",
                  textAlign: TextAlign.left,
                  style: TextStyle(fontSize: 16),
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "$Rupee${order.itemPrice}",
                  textAlign: TextAlign.left,
                  style: TextStyle(fontSize: 16),
                ),
              ),
            ]),
            TableRow(children: [
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "Delivery Price:",
                  textAlign: TextAlign.left,
                  style: TextStyle(fontSize: 16),
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "$Rupee${order.deliveryPrice}",
                  textAlign: TextAlign.left,
                  style: TextStyle(fontSize: 16),
                ),
              ),
            ]),
            TableRow(children: [
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "Tax",
                  textAlign: TextAlign.left,
                  style: TextStyle(fontSize: 16),
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "$Rupee${order.taxPrice}",
                  textAlign: TextAlign.left,
                  style: TextStyle(fontSize: 16),
                ),
              ),
            ]),
            TableRow(children: [
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "Total",
                  textAlign: TextAlign.left,
                  style: TextStyle(fontSize: 16),
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "$Rupee${order.totalPrice}",
                  textAlign: TextAlign.left,
                  style: TextStyle(fontSize: 16),
                ),
              ),
            ]),
            TableRow(children: [
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "Status",
                  textAlign: TextAlign.left,
                  style: TextStyle(fontSize: 16),
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Text(
                  "${order.status.toString()}",
                  textAlign: TextAlign.left,
                  style: TextStyle(
                    fontSize: 16, color: statusColor
                  ),
                ),
              ),
            ]),
          ],
        ));
    final body = SingleChildScrollView(
      child: Column(
        children: [
          GestureDetector(
            onTap: () => _gotoItemScreen(context, order.itemID),
            child: Container(
                height: 200,
                width: 200,
                decoration: BoxDecoration(
                    shape: BoxShape.rectangle,
                    image: DecorationImage(
                        fit: BoxFit.fill, image: NetworkImage(order.imageURL)))),
          ),
          SizedBox(
            height: 20,
          ),
          orderContents,
          if (order.status == OrderStatus.KOrderPaymentPending)
            RaisedButton(
              onPressed: _handlePayment,
              child: const Text("Make payment"),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(16.0),
              ),
            )
        ],
      ),
    );
    return body;
  }

  @override
  Widget build(BuildContext context) {
    final order = ModalRoute.of(context).settings.arguments as Order;
    final appBar = AppBar(
      backgroundColor: PrimaryColor,
      title: Text(order.itemName),
    );
    return Scaffold(appBar: appBar, body: _buildBody(context, order));
  }
}
