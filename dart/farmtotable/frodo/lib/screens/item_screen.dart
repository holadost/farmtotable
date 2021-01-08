import 'package:flutter/material.dart';
import 'package:frodo/net/rest_api_client.dart';

import '../models/item.dart';
import '../util/constants.dart';
import '../widgets/item_display_widget.dart';

class ItemScreen extends StatefulWidget {
  static const routeName = '/item-auction-screen';

  @override
  _ItemScreenState createState() => _ItemScreenState();
}

class _ItemScreenState extends State<ItemScreen> {
  bool _showBiddingButton;
  String _itemID;
  bool _isLoading;
  var _apiClient = RestApiClient();
  Item _item;

  @override
  void didChangeDependencies() {
    final args = ModalRoute.of(context).settings.arguments as Map<String, dynamic>;
    _showBiddingButton = args['show_bid_button'];
    _itemID = args['item_id'];
    _loadData();
    super.didChangeDependencies();
  }

  void _bidNow() {
    print("Bidding now");
  }

  void _loadData() async {
    // Loads all the required auctions.
    print("Fetching item from backend");
    try {
      setState(() {
        print("Currently loading");
        _isLoading = true;
      });

      final item = await _apiClient.getItem(_itemID);
      setState(() {
        _item = item;
        _isLoading = false;
        print("Finished loading");
      });
      print("Successfully fetched item from backend");
    } catch (error) {
      print("Unable to load data due to error: $error");
    }
  }

  AppBar _buildAppBar() {
    final appBar = AppBar(
      backgroundColor: PrimaryColor,
      title: Text(_item.itemName),
      actions: [
        if (_showBiddingButton)
          IconButton(
            onPressed: _bidNow,
            icon: Icon(Icons.shopping_cart),
          )
      ],
    );
    return appBar;
  }

  @override
  Widget build(BuildContext context) {
    Function bidNow;
    if (_showBiddingButton) {
      bidNow = _bidNow;
    }
    return Scaffold(
      appBar: _buildAppBar(),
      body: ItemDisplayWidget(
        item: _item,
        bidNow: bidNow,
      ),
    );
  }
}
