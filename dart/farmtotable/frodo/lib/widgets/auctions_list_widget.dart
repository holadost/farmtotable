import 'package:flutter/material.dart';
import 'package:frodo/util/logging.dart';
import 'package:intl/intl.dart';

import '../screens/item_screen.dart';
import '../models/auction_item.dart';
import '../util/constants.dart';

class AuctionsListWidget extends StatelessWidget {
  final List<AuctionItem> _auctions;
  final Function loadMore;
  AuctionsListWidget(this._auctions, this.loadMore);

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      itemBuilder: (ctx, ii) {
        return (ii == _auctions.length) ? Container(
          height: 50,
          color: PrimaryColor,
          child: FlatButton(
            child: Text(
              "Load More",
              style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold),),
            onPressed: () {
              loadMore();
            },
          ),
        ) :
        Container(
          height: 100,
          child: ListTile(
            onTap: () {
              Navigator.of(ctx).pushNamed(ItemScreen.routeName, arguments: {
                "item_id": _auctions[ii].itemID,
                "show_bid_button": true
              });
            },
            leading: CircleAvatar(
                backgroundColor: PrimaryColor,
                radius: 30,
                child: Container(
                    height: 250,
                    width: double.infinity,
                    decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        image: DecorationImage(
                            fit: BoxFit.fill,
                            image: NetworkImage(_auctions[ii].imageURL))))),
            title: Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 2.0, vertical: 10.0),
              child: Text(
                _auctions[ii].itemName,
                style: Theme.of(context).textTheme.headline6,
                textAlign: TextAlign.left,
              ),
            ),
            subtitle: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Deadline: ${DateFormat.yMMMMd().add_jm().format(_auctions[ii].auctionStartTime.add(_auctions[ii].auctionDurationSecs))}',
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
                  'Min price: $Rupee${_auctions[ii].minBid.toStringAsFixed(2)}',
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
                SizedBox(
                  height: 20,
                ),
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
                      "$Rupee${_auctions[ii].maxBid.toStringAsFixed(2)}",
                      style: TextStyle(fontSize: 18),
                    ),
                  ),
                ),
              ],
            ),
          ),
        );
      },
      itemCount: _auctions.length+1,
    );
  }
}
