import 'package:flutter/material.dart';

import '../widgets/side_drawer_widget.dart';

class AuctionsOverviewScreen extends StatelessWidget {
  static const routeName = '/auctions-overview-screen';

  @override
  Widget build(BuildContext context) {
    final appBar = AppBar();
    final body = Container(
      child: Text("Your auctions"),
    );
    return Scaffold(
      appBar: appBar,
      body: body,
      drawer: SideDrawerWidget(),
    );
  }
}
