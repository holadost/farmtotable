import 'package:flutter/material.dart';

import '../models/item.dart';
import '../util/constants.dart';
import '../widgets/item_display_widget.dart';

class ItemScreen extends StatelessWidget {
  static const routeName = '/item-auction-screen';
  bool showBiddingButton;
  Item item;

  void _bidNow() {
    print("Bidding now");
  }

  @override
  Widget build(BuildContext context) {
    final args =
        ModalRoute.of(context).settings.arguments as Map<String, dynamic>;
    showBiddingButton = args['show_bid_button'];
    item = args['item'];
    final appBar = AppBar(
      backgroundColor: PrimaryColor,
      title: Text(item.itemName),
      actions: [
        if (showBiddingButton)
          IconButton(
            onPressed: _bidNow,
            icon: Icon(Icons.shopping_cart),
          )
      ],
    );
    Function bidNow;
    if (showBiddingButton) {
      bidNow = _bidNow;
    }
    return Scaffold(
      appBar: appBar,
      body: ItemDisplayWidget(
        item: item,
        bidNow: bidNow,
      ),
    );
  }
}
