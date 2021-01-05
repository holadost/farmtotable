import 'package:flutter/material.dart';

import '../widgets/side_drawer_widget.dart';

class HomeScreen extends StatelessWidget {
  static const routeName = "/";

  @override
  Widget build(BuildContext context) {
    final appBar = AppBar(
      title: Text('AlgoRhythm'),
    );
    final body = Container(
      child: Text('Welcome to AlgoRhythm. Still under construction!'),
    );
    return Scaffold(
      appBar: appBar,
      body: body,
      drawer: SideDrawerWidget(),
    );
  }
}
