import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../widgets/side_drawer_widget.dart';
import '../data/dummy_auctions.dart';
import '../screens/item_auction_screen.dart';

class AuctionsOverviewScreen extends StatelessWidget {
  static const routeName = '/auctions-overview-screen';

  @override
  Widget build(BuildContext context) {
    var auctions = [...DUMMY_AUCTIONS];
    final themeData = Theme.of(context);
    final appBar = AppBar(
      title: Text('Auctions', style: themeData.textTheme.headline6),
    );
    final body = Container(
        child: ListView.builder(
      itemBuilder: (ctx, ii) {
        return ListTile(
          onTap: () {
            Navigator.of(ctx).pushNamed(
                ItemAuctionScreen.routeName,
                arguments: auctions[ii]);
          },
          leading: CircleAvatar(
            radius: 30,
            child: Padding(
                padding: const EdgeInsets.all(6.0),
                child: FittedBox(
                  child: Text(
                    "Rs " +
                        auctions[ii]
                            .minBid
                            .toStringAsPrecision(4),
                    style: TextStyle(fontSize: 20),
                  ),
                )),
          ),
          title: Text(
            auctions[ii].itemName,
            style: Theme.of(context).textTheme.title,
            textAlign: TextAlign.left,
          ),
          subtitle: Text(
            DateFormat.yMMMMEEEEd()
                .add_jm()
                .format(auctions[ii].auctionStartTime),
            style: TextStyle(
              fontSize: 12,
              color: Colors.grey,
            ),
            textAlign: TextAlign.left,
          ),
        );
      },
      itemCount: auctions.length,
    ));
    return Scaffold(
      appBar: appBar,
      body: body,
      drawer: SideDrawerWidget(),
    );
  }
}
