import 'package:flutter/material.dart';

import '../widgets/side_drawer_widget.dart';

class ItemAuctionScreen extends StatelessWidget {
  static const routeName = '/item-auction-screen';

  @override
  Widget build(BuildContext context) {
    final appBar = AppBar();
    final body = Container(
      child: Text("Item Auction"),
    );
    return Scaffold(
      appBar: appBar,
      body: body,
    );
  }
}
