import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import './screens/home_screen.dart';
import './screens/auctions_overview_screen.dart';
import './screens/item_screen.dart';
import './screens/orders_overview_screen.dart';
import './screens/order_screen.dart';
import './screens/bid_screen.dart';
import './util/constants.dart';
import './providers/auth_provider.dart';
import './providers/rest_client_provider.dart';


void main() => runApp(MyApp());

class MyApp extends StatelessWidget {

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider.value(value: AuthProvider()),
        // ignore: missing_required_param
        ChangeNotifierProxyProvider<AuthProvider, RestClientProvider>(
            update: (ctx, auth, prevClient){
              return RestClientProvider(auth.token);
            }),
      ],
      child: MaterialApp(
        title: AppName,
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
          ItemScreen.routeName: (ctx) => ItemScreen(),
          OrderScreen.routeName: (ctx) => OrderScreen(),
          BidScreen.routeName: (ctx) => BidScreen(),
        },
      ),
    );
  }
}
