import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../providers/auth_provider.dart';
import '../screens/home_screen.dart';
import '../screens/contact_us_screen.dart';
import '../screens/orders_overview_screen.dart';
import '../screens/auctions_overview_screen.dart';
import '../util/constants.dart';

class SideDrawerWidget extends StatelessWidget {

  @override
  Widget build(BuildContext context) {
    var auth = Provider.of<AuthProvider>(context);
    return Drawer(
      child: Column(
        children: [
          AppBar(
            backgroundColor: PrimaryColor,
            title: Text(AppName),
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
          Divider(),
          ListTile(
            leading: Icon(Icons.people), title: Text('Contact Us'),
            onTap: () {
              // Go back to home screen.
              Navigator.of(context).pushReplacementNamed(
                  ContactUsScreen.routeName);
            },
          ),
          Divider(),
          ListTile(
            leading: Icon(Icons.logout), title: Text('Log out'),
            onTap: () {
              auth.signout();
              // Go back to home screen.
              Navigator.of(context).popUntil((route) => route.isFirst);
              Navigator.of(context).pushReplacementNamed(HomeScreen.routeName);
            },
          ),
        ],
      ),
    );
  }
}
