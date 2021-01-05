import 'package:flutter/material.dart';
import '../screens/home_screen.dart';
import '../screens/orders_overview_screen.dart';
import '../screens/auctions_overview_screen.dart';

class SideDrawerWidget extends StatelessWidget {

  @override
  Widget build(BuildContext context) {
    return Drawer(
      child: Column(
        children: [
          AppBar(
            title: Text('Hello Friend'),
            automaticallyImplyLeading: false,
          ),
          Divider(),
          ListTile(
            leading: Icon(Icons.home), title: Text('Home'),
            onTap: () {
              Navigator.of(context).pushReplacementNamed(
                  HomeScreen.routeName);
            },
          ),
          Divider(),
          ListTile(
            leading: Icon(Icons.shop), title: Text('Auctions'),
            onTap: () {
              Navigator.of(context).pushReplacementNamed(
                  AuctionsOverviewScreen.routeName);
            },
          ),
          Divider(),
          ListTile(
            leading: Icon(Icons.payment), title: Text('Orders'),
            onTap: () {
              Navigator.of(context).pushReplacementNamed(
                  OrdersOverviewScreen.routeName);
            },
          ),
        ],
      ),
    );
  }
}
