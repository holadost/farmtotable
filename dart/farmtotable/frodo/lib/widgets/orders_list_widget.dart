import 'package:flutter/material.dart';

import '../screens/item_screen.dart';
import '../models/order.dart';
import '../util/constants.dart';

class OrdersListWidget extends StatelessWidget {
  final List<Order> _orders;
  OrdersListWidget(this._orders);

  @override
  Widget build(BuildContext context) {
    final themeData = Theme.of(context);
    return ListView.builder(
      itemBuilder: (ctx, ii) {
        return Container(
          height: 100,
          child: ListTile(
            onTap: () {
              Navigator.of(ctx).pushNamed(ItemScreen.routeName, arguments: {
                "item_id": _orders[ii].itemID,
                "show_bid_button": true
              });
            },
            leading: CircleAvatar(
                backgroundColor: PrimaryColor,
                radius: 30,
                child: Container(
                    height: 250,
                    width: double.infinity,
                    decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        image: DecorationImage(
                            fit: BoxFit.fill,
                            image: NetworkImage(_orders[ii].imageURL))))),
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
              'Min price: $Rupee${_orders[ii].status}',
              style: TextStyle(
                fontSize: 12,
                color: Colors.grey,
              ),
              textAlign: TextAlign.left,
            ),
            trailing: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                SizedBox(
                  height: 20,
                ),
                Container(
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(10),
                    color: Colors.green,
                  ),
                  padding: const EdgeInsets.all(3.0),
                  height: 30,
                  width: 80,
                  child: FittedBox(
                    fit: BoxFit.contain,
                    child: Text(
                      "$Rupee${_auctions[ii].maxBid.toStringAsFixed(2)}",
                      style: TextStyle(fontSize: 18),
                    ),
                  ),
                ),
              ],
            ),
          ),
        );
      },
      itemCount: _auctions.length,
    );
  }
}
