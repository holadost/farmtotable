import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/order.dart';
import '../providers/aragorn_client_provider.dart';
import '../widgets/orders_list_widget.dart';
import '../widgets/side_drawer_widget.dart';
import '../util/styles.dart';
import '../util/constants.dart';

class OrdersOverviewScreen extends StatefulWidget {
  static const routeName = '/orders-overview-screen';

  @override
  _OrdersOverviewScreenState createState() => _OrdersOverviewScreenState();
}

class _OrdersOverviewScreenState extends State<OrdersOverviewScreen> {
  AragornClientProvider apiClient;
  List<Order> _ordersList = [];
  bool _isLoading = false;

  void _loadData() async {
    // Loads all the required auctions.
    List<Order> orders = [];
    try {
      setState(() {
        _isLoading = true;
      });
      orders = await apiClient.getUserOrders("user1");
    } finally {
      setState(() {
        _ordersList = [...orders];
        _isLoading = false;
      });
    }
  }

  @override
  void didChangeDependencies() {
    apiClient = Provider.of<AragornClientProvider>(context, listen: false);
    _loadData();
    super.didChangeDependencies();
  }

  @override
  Widget build(BuildContext context) {
    final appBar = AppBar(
      backgroundColor: PrimaryColor,
      title: Text(
        'Orders',
        style: getAppBarTextStyle(),
      ),
    );
    final body = _isLoading
        ? Center(child: CircularProgressIndicator())
        : Padding(
            padding: const EdgeInsets.all(8.0),
            child: OrdersListWidget(_ordersList),
          );
    return Scaffold(
      appBar: appBar,
      body: body,
      drawer: SideDrawerWidget(),
    );
  }
}
