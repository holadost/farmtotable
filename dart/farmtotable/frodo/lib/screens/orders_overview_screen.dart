import 'package:flutter/material.dart';

import '../widgets/side_drawer_widget.dart';

class OrdersOverviewScreen extends StatelessWidget {
  static const routeName = '/orders-overview-screen';

  @override
  Widget build(BuildContext context) {
    final appBar = AppBar();
    final body = Container(
      child: Text("Your orders"),
    );
    return Scaffold(
      appBar: appBar,
      body: body,
      drawer: SideDrawerWidget(),
    );
  }
}
