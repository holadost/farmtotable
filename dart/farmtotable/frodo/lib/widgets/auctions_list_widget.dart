import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../screens/item_screen.dart';
import '../models/auction_item.dart';
import '../util/constants.dart';

class AuctionsListWidget extends StatelessWidget {

  final List<AuctionItem> _auctions;
  AuctionsListWidget(this._auctions);

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      itemBuilder: (ctx, ii) {
        return Container(
          height: 100,
          child: ListTile(
            onTap: () {
              Navigator.of(ctx).pushNamed(ItemScreen.routeName,
                  arguments: _auctions[ii]);
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
                  'Deadline: ${DateFormat.yMMMMd().add_jm().format(_auctions[ii].auctionStartTime)}',
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
      itemCount: _auctions.length,
    );
  }
}
