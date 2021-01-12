import 'package:flutter/material.dart';

import '../util/logging.dart';
import '../models/auction_item.dart';
import '../net/aragorn_rest_client.dart';
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
  final apiClient = AragornRestClient();
  bool _isLoading = false;
  int _lastID = -1;
  int _numItemsPerPage = 8;
  final List<AuctionItem> _auctions = [];

  @override
  void didChangeDependencies() {
    _cleanLoad();
    super.didChangeDependencies();
  }

  void _cleanLoad() {
    _fetchData(true);
  }

  void _loadMore() {
    _fetchData(false);
  }

  void _fetchLatestBids() {
    // TODO: Implement this.
  }

  void _fetchData(bool clean) async {
    // Loads all the required auctions.
    List<AuctionItem> auctions = [];
    try {
      setState(() {
        if (clean) {
          _lastID = -1;
          // Avoid reloading screen if we are not doing a
          // clean fetch. This ensures that during loadMore,
          // we do not go back to the top of the screen.
          _isLoading = true;
        }
      });
      info("Fetching auction starting from ID: ${_lastID + 1}");
      auctions = await apiClient.getAuctions(_lastID + 1, _numItemsPerPage);
      if ((auctions != null) && (auctions.length > 0)) {
        _lastID = _lastID + auctions.length;
      }
      print("Last ID fetched: $_lastID");
    } catch (error) {
      print("Failed to fetch data due to error: $error");
    } finally {
      setState(() {
        if (clean) {
          _auctions.clear();
        }
        _isLoading = false;
        if (auctions != null) {
          auctions.forEach((element) {
            _auctions.add(element);
          });
        }

      });
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
              _cleanLoad();
            }),
      ],
    );
    return appBar;
  }

  Widget _buildBody() {
    final body = _isLoading
        ? Center(child: CircularProgressIndicator())
        : AuctionsListWidget(_auctions, _loadMore);
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
