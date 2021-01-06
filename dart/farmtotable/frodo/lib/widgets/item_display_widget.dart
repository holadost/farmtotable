import 'package:flutter/material.dart';
import 'package:frodo/util/constants.dart';
import 'package:intl/intl.dart';

import '../models/auction_item.dart';

class ItemDisplayWidget extends StatelessWidget {
  final Function bidNow;
  final AuctionItem auctionItem;
  ItemDisplayWidget({this.bidNow, this.auctionItem});

  @override
  Widget build(BuildContext context) {
    double minWidth = 300;
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          SizedBox(
            height: 20,
          ),
          Container(
              height: 250,
              width: double.infinity,
              decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  image: DecorationImage(
                      fit: BoxFit.fill,
                      image: NetworkImage(auctionItem.imageURL)))),
          SizedBox(
            height: 20,
          ),
          Container(
            child: Text(
              "Min Price: $Rupee${auctionItem.minBid.toStringAsPrecision(4)}",
              style: TextStyle(
                  fontFamily: "Quicksand",
                  fontWeight: FontWeight.bold,
                  color: Colors.green,
                  fontSize: 20),
            ),
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(15),
            ),
            width: minWidth,
          ),
          Container(
            child: Text(
              "Deadline: ${DateFormat.yMMMMd().add_jm().format(auctionItem.auctionStartTime)}",
              style: TextStyle(
                fontFamily: "Quicksand",
                fontWeight: FontWeight.bold,
                color: Colors.red,
                fontSize: 18,
              ),
            ),
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(15),
            ),
            width: minWidth,
          ),
          Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(15),
            ),
            width: minWidth,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                SizedBox(
                  height: 5,
                ),
                Text(
                  auctionItem.itemDescription,
                  style: TextStyle(fontSize: 16, fontFamily: 'Quicksand'),
                ),
              ],
            ),
          ),
          SizedBox(
            height: 10,
          ),
          ElevatedButton(
            child: Text("Bid now"),
            onPressed: bidNow,
          )
        ],
      ),
    );
  }
}
