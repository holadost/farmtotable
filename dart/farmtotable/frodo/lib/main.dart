import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import './screens/home_screen.dart';
import './screens/auctions_overview_screen.dart';
import './screens/orders_overview_screen.dart';


void main() => runApp(MyApp());

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'AlgoRhythm',
      theme: ThemeData(
        primarySwatch: Colors.purple,
        accentColor: Colors.deepOrange,
        fontFamily: 'Lato',
      ),
      home: HomeScreen(),
      routes: {
        OrdersOverviewScreen.routeName: (ctx) => OrdersOverviewScreen(),
        AuctionsOverviewScreen.routeName: (ctx) => AuctionsOverviewScreen(),
      },
    );
  }
}
