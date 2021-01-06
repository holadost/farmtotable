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
                child: Container(
                    height: 250,
                    width: double.infinity,
                    decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        image: DecorationImage(
                            fit: BoxFit.fill,
                            image: NetworkImage(auctions[ii].imageURL))))),
            title: Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 2.0, vertical: 10.0),
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
                  'Deadline: ${DateFormat.yMMMMd().add_jm().format(auctions[ii].auctionStartTime)}',
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
                  'Min price: $Rupee${auctions[ii].minBid.toStringAsPrecision(4)}',
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.grey,
                  ),
                  textAlign: TextAlign.left,
                )
              ],
            ),
            trailing: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                SizedBox(height: 20,),
                Container(
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(10),
                    color: Colors.green,
                  ),
                  padding: const EdgeInsets.all(3.0),
                  height: 30,
                  width: 80,
                  child: FittedBox(
                    fit: BoxFit.contain,
                    child: Text(
                      "$Rupee${auctions[ii].maxBid.toStringAsFixed(2)}",
                      style: TextStyle(fontSize: 18),
                    ),
                  ),
                ),
              ],
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
