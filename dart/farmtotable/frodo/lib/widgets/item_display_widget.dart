import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../models/item.dart';
import '../util/constants.dart';

class ItemDisplayWidget extends StatelessWidget {
  final Function bidNow;
  final Item item;

  ItemDisplayWidget({this.bidNow, this.item});

  List<Widget> _buildBodyChildren() {
    double minWidth = 300;
    List<Widget> children = [
      SizedBox(
        height: 20,
      ),
      Container(
          height: 250,
          width: double.infinity,
          decoration: BoxDecoration(
              shape: BoxShape.circle,
              image: DecorationImage(
                  fit: BoxFit.fill, image: NetworkImage(item.imageURL)))),
      SizedBox(
        height: 20,
      ),
      Container(
        child: Text(
          "Min Price: $Rupee${item.minBidPrice.toStringAsPrecision(4)}",
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
      if (bidNow != null)
        Container(
          child: Text(
            "Deadline: ${DateFormat.yMMMMd().add_jm().format(item.auctionStartTime.add(item.auctionDurationSecs))}",
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
              item.itemDescription,
              style: TextStyle(fontSize: 16, fontFamily: 'Quicksand'),
            ),
          ],
        ),
      ),
      SizedBox(
        height: 10,
      ),
      if (bidNow != null)
        ElevatedButton(
          child: Text("Bid now"),
          onPressed: bidNow,
        ),
    ];
    return children;
  }

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: _buildBodyChildren(),
      ),
    );
  }
}
