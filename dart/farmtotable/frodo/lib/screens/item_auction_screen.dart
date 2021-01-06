import 'package:flutter/material.dart';

import '../models/auction_item.dart';

class ItemAuctionScreen extends StatelessWidget {
  static const routeName = '/item-auction-screen';

  @override
  Widget build(BuildContext context) {
    final auctionItem =
        ModalRoute.of(context).settings.arguments as AuctionItem;
    final themeData = Theme.of(context);
    final appBar = AppBar(
      title: Text(auctionItem.itemName),
    );
    final body = SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          SizedBox(height: 20,),
          Container(
              height: 250,
              width: double.infinity,
              decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  image: DecorationImage(
                      fit: BoxFit.fill,
                      image: NetworkImage(
                          "https://cdn.pixabay.com/photo/2018/07/11/21/51/toast-3532016_1280.jpg")))),
          ClipRRect(
            child: Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                border: Border.all(color: Colors.green),
                gradient: LinearGradient(
                  colors: [
                    themeData.primaryColor.withOpacity(0.7),
                    themeData.primaryColor],
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
                borderRadius: BorderRadius.circular(15),
              ),
              width: 300,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  SizedBox(height: 5,),
                  Text(
                    auctionItem.itemDescription,
                    style: TextStyle(
                        fontSize: 16,
                        fontFamily: 'Lato'),
                  ),
                ],
              ),
            ),
          ),
          SizedBox(height: 20,),
          ElevatedButton(
            child: Text("Bid now"),
            onPressed: () {},
          )
        ],
      ),
    );
    return Scaffold(
      appBar: appBar,
      body: body,
    );
  }
}
