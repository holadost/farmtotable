import 'package:flutter/material.dart';

import '../widgets/side_drawer_widget.dart';
import '../util/styles.dart';
import '../util/constants.dart';
import '../data/dummy_orders.dart';
import '../screens/order_screen.dart';

class OrdersOverviewScreen extends StatelessWidget {
  static const routeName = '/orders-overview-screen';

  @override
  Widget build(BuildContext context) {
    var orders = [...DUMMY_ORDERS];
    final appBar = AppBar(
      backgroundColor: PrimaryColor,
      title: Text(
        'Orders',
        style: getAppBarTextStyle(),
      ),
    );
    final body = ListView.builder(
      itemBuilder: (ctx, ii) {
        return Container(
          height: 100,
          child: ListTile(
            onTap: () {
              Navigator.of(ctx).pushNamed(OrderScreen.routeName,
                  arguments: orders[ii]);
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
                            image: NetworkImage(orders[ii].itemImageURL))))),
            title: Padding(
              padding:
              const EdgeInsets.symmetric(horizontal: 2.0, vertical: 10.0),
              child: Text(
                orders[ii].itemName,
                style: Theme.of(context).textTheme.headline6,
                textAlign: TextAlign.left,
              ),
            ),
            subtitle: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Quantity: ',
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey,
                  ),
                  textAlign: TextAlign.left,
                ),
                SizedBox(
                  height: 3,
                ),
                Text(
                  'Total Price: $Rupee${orders[ii].price.toStringAsFixed(2)}',
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey,
                  ),
                  textAlign: TextAlign.left,
                )
              ],
            ),
          ),
        );
      },
      itemCount: orders.length,
    );
    return Scaffold(
      appBar: appBar,
      body: body,
      drawer: SideDrawerWidget(),
    );
  }
}
