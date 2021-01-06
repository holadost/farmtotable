import 'package:flutter/material.dart';
import 'package:frodo/util/constants.dart';

import './screens/home_screen.dart';
import './screens/auctions_overview_screen.dart';
import './screens/orders_overview_screen.dart';
import './screens/item_auction_screen.dart';


void main() => runApp(MyApp());

class MyApp extends StatelessWidget {

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'AlgoRhythm',
      theme: ThemeData(
        brightness: Brightness.dark,
        primaryColor: PrimaryColor,
        primarySwatch: Colors.deepPurple,
        accentColor: AccentColor,
        fontFamily: 'Lato',
      ),
      home: HomeScreen(),
      routes: {
        OrdersOverviewScreen.routeName: (ctx) => OrdersOverviewScreen(),
        AuctionsOverviewScreen.routeName: (ctx) => AuctionsOverviewScreen(),
        ItemAuctionScreen.routeName: (ctx) => ItemAuctionScreen(),
      },
    );
  }
}
