import 'package:flutter/material.dart';

import '../models/item.dart';
import '../net/aragorn_rest_client.dart';
import '../screens/bid_screen.dart';
import '../util/constants.dart';
import '../widgets/item_display_widget.dart';

class ItemScreen extends StatefulWidget {
  static const routeName = '/item-auction-screen';

  @override
  _ItemScreenState createState() => _ItemScreenState();
}

class _ItemScreenState extends State<ItemScreen> {
  bool _showBiddingButton = false;
  String _itemID;
  bool _isLoading = false;
  var _apiClient = AragornRestClient();
  Item _item;
  bool _gatheredArgs = false;

  @override
  void didChangeDependencies() {
    if (!_gatheredArgs) {
      final args = ModalRoute.of(context).settings.arguments as Map<String, dynamic>;
      _showBiddingButton = args['show_bid_button'];
      _itemID = args['item_id'];
      _gatheredArgs = true;
    }
    print("Show Bidding: $_showBiddingButton, Item ID: $_itemID");
    _loadData();
    super.didChangeDependencies();
  }

  void _gotoBidScreen() {
    Navigator.of(context).pushNamed(BidScreen.routeName, arguments: _item);
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
      _item = item;
      print("Successfully fetched item from backend");
    } catch (error) {
      print("Unable to load data due to error: $error");
    } finally {
      setState(() {
        _isLoading = false;
        print("Finished loading");
      });
    }
  }

  AppBar _buildAppBar(BuildContext context) {
    final appBar = AppBar(
      backgroundColor: PrimaryColor,
      title: _isLoading ? const Text("") : Text(_item.itemName),
      actions: [
        if (_showBiddingButton)
          IconButton(
            onPressed: _gotoBidScreen,
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
      bidNow = _gotoBidScreen;
    }
    return Scaffold(
      appBar: _buildAppBar(context),
      body: _isLoading
          ? Center(
              child: CircularProgressIndicator(),
            )
          : ItemDisplayWidget(
              item: _item,
              bidNow: bidNow,
            ),
    );
  }
}
