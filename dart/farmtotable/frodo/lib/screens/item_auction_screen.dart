import 'package:flutter/material.dart';
import 'package:frodo/util/constants.dart';
import 'package:intl/intl.dart';

import '../models/auction_item.dart';
import '../widgets/item_display_widget.dart';

class ItemAuctionScreen extends StatelessWidget {
  static const routeName = '/item-auction-screen';

  void _bidNow() {

  }

  @override
  Widget build(BuildContext context) {
    final auctionItem =
        ModalRoute.of(context).settings.arguments as AuctionItem;
    final appBar = AppBar(
      backgroundColor: PrimaryColor,
      title: Text(auctionItem.itemName),
      actions: [
        IconButton(onPressed: _bidNow, icon: Icon(Icons.shopping_cart),)
      ],
    );
    return Scaffold(
      appBar: appBar,
      body: ItemDisplayWidget(
        auctionItem: auctionItem, bidNow: _bidNow,
      ),
    );
  }
}
