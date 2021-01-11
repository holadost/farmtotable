import 'package:flutter/material.dart';

import '../screens/order_screen.dart';
import '../models/order.dart';
import '../util/constants.dart';
import '../util/custom_icons.dart';

class OrdersListWidget extends StatelessWidget {
  final List<Order> _orders;
  OrdersListWidget(this._orders);

  @override
  Widget build(BuildContext context) {
    final themeData = Theme.of(context);
    return ListView.builder(
      itemBuilder: (ctx, ii) {
        var statusColor;
        if (_orders[ii].status.toInt() ==
            OrderStatus.KOrderPaymentPending.toInt()) {
          statusColor = Colors.red;
        } else {
          statusColor = Colors.green;
        }
        return Container(
          height: 100,
          child: ListTile(
              onTap: () {
                Navigator.of(ctx).pushNamed(OrderScreen.routeName,
                    arguments: _orders[ii]);
              },
              leading: Container(
                  height: 250,
                  width: double.infinity,
                  decoration: BoxDecoration(
                      image: DecorationImage(
                          fit: BoxFit.fill,
                          image: NetworkImage(_orders[ii].imageURL)))),
              title: Padding(
                padding:
                    const EdgeInsets.symmetric(horizontal: 2.0, vertical: 10.0),
                child: Text(
                  _orders[ii].itemName,
                  style: themeData.textTheme.headline6,
                  textAlign: TextAlign.left,
                ),
              ),
              subtitle: Text(
                '${_orders[ii].status}',
                style: TextStyle(
                  fontSize: 12,
                  color: statusColor,
                ),
                textAlign: TextAlign.left,
              ),
              trailing: IconButton(
                icon: Icon(CustomIcons.chevron_right),
                onPressed: () {},
              )),
        );
      },
      itemCount: _orders.length,
    );
  }
}
