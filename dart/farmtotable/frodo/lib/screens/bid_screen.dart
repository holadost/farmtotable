import 'package:flutter/material.dart';

import '../models/item.dart';
import '../widgets/register_bid_widget.dart';

class BidScreen extends StatelessWidget {
  static const String routeName = "/bid-screen";

  @override
  Widget build(BuildContext context) {
    final item = ModalRoute.of(context).settings.arguments as Item;
    return Scaffold(
      appBar: AppBar(
        title: Text(item.itemName),
      ),
      body: RegisterBidWidget(item)
    );
  }
}
