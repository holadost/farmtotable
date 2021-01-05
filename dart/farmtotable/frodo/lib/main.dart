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
        brightness: Brightness.light,
        primarySwatch: Colors.purple,
        fontFamily: 'Lato',
      ),
      darkTheme: ThemeData(
        brightness: Brightness.dark,
        primarySwatch: Colors.orange,
        fontFamily: 'Lato',
      ),
      themeMode: ThemeMode.dark,
      home: HomeScreen(),
      routes: {
        OrdersOverviewScreen.routeName: (ctx) => OrdersOverviewScreen(),
        AuctionsOverviewScreen.routeName: (ctx) => AuctionsOverviewScreen(),
      },
    );
  }
}
