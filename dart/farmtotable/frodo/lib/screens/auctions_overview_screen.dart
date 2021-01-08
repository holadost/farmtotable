import 'package:flutter/material.dart';

import '../models/auction_item.dart';
import '../net/rest_api_client.dart';
import '../util/styles.dart';
import '../util/constants.dart';
import '../widgets/side_drawer_widget.dart';
import '../widgets/auctions_list_widget.dart';

class AuctionsOverviewScreen extends StatefulWidget {
  static const routeName = '/auctions-overview-screen';

  @override
  _AuctionsOverviewScreenState createState() => _AuctionsOverviewScreenState();
}

class _AuctionsOverviewScreenState extends State<AuctionsOverviewScreen> {
  final apiClient = RestApiClient();
  bool _isLoading = false;
  int _lastID = -1;
  int _numItemsPerPage = 8;
  List<AuctionItem> _auctions = [];

  @override
  void didChangeDependencies() {
    _loadData();
    super.didChangeDependencies();
  }

  void _loadData() async {
    // Loads all the required auctions.
    print("Fetching data from backend");
    try {
      setState(() {
        print("Currently loading");
        _isLoading = true;
      });

      final auctions =
          await apiClient.getAuctions(_lastID + 1, _numItemsPerPage);
      setState(() {
        _auctions = [...auctions];
        print("Finished loading");
        _isLoading = false;
      });
      print("Successfully fetched data from backend");
    } catch (error) {

    }
  }

  AppBar _buildAppBar() {
    final appBar = AppBar(
      backgroundColor: PrimaryColor,
      title: Text(
        'Auctions',
        style: getAppBarTextStyle(),
      ),
      actions: [
        IconButton(
            icon: Icon(Icons.refresh),
            onPressed: () {
              // Refresh page.
              _loadData();
            }),
      ],
    );
    return appBar;
  }


  Widget _buildBody() {
    final body = _isLoading
        ? Center(child: CircularProgressIndicator())
        : AuctionsListWidget(_auctions);
    return body;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildAppBar(),
      body: _buildBody(),
      drawer: SideDrawerWidget(),
    );
  }
}
