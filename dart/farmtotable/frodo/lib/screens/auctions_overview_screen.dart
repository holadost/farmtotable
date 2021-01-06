import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../widgets/side_drawer_widget.dart';
import '../data/dummy_auctions.dart';
import '../screens/item_auction_screen.dart';
import '../util/styles.dart';
import '../util/constants.dart';

class AuctionsOverviewScreen extends StatelessWidget {
  static const routeName = '/auctions-overview-screen';

  @override
  Widget build(BuildContext context) {
    var auctions = [...DUMMY_AUCTIONS];
    final appBar = AppBar(
      backgroundColor: PrimaryColor,
      title: Text(
        'Auctions',
        style: getAppBarTextStyle(),
      ),
    );
    final body = ListView.builder(
      itemBuilder: (ctx, ii) {
        return Container(
          height: 100,
          child: ListTile(
            onTap: () {
              Navigator.of(ctx).pushNamed(ItemAuctionScreen.routeName,
                  arguments: auctions[ii]);
            },
            leading: CircleAvatar(
              backgroundColor: Colors.green,
              radius: 30,
              child: Padding(
                  padding: const EdgeInsets.all(6.0),
                  child: FittedBox(
                    child: Text(
                      "Rs " + auctions[ii].minBid.toStringAsPrecision(4),
                      style: TextStyle(
                          fontSize: 20, color: Colors.black, fontFamily: 'Lato'),
                    ),
                  )),
            ),
            title: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 2.0, vertical: 10.0),
              child: Text(
                auctions[ii].itemName,
                style: Theme.of(context).textTheme.headline6,
                textAlign: TextAlign.left,
              ),
            ),
            subtitle: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Deadline: ${DateFormat.yMMMMd()
                      .add_jm()
                      .format(auctions[ii].auctionStartTime)}',
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey,
                  ),
                  textAlign: TextAlign.left,
                ),
                SizedBox(
                  height: 3,
                ),
                Text(
                  'Min price: Rs ${auctions[ii].minBid.toStringAsPrecision(4)}',
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey,
                  ),
                  textAlign: TextAlign.left,
                )
              ],
            ),
            trailing: Container(
              width: 80.0,
              padding: const EdgeInsets.all(8.0),
              child: RaisedButton(
                shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(10.0)
                ),
                onPressed: () {},
                child: Text(
                  "Bid", style: TextStyle(fontSize: 16),),
              ),
            ),
          ),
        );
      },
      itemCount: auctions.length,
    );
    return Scaffold(
      appBar: appBar,
      body: body,
      drawer: SideDrawerWidget(),
    );
  }
}
