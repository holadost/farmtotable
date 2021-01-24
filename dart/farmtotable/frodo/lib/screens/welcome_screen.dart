import 'package:flutter/material.dart';

import '../widgets/side_drawer_widget.dart';
import '../util/constants.dart';
import '../util/styles.dart';

class WelcomeScreen extends StatelessWidget {
  static const routeName = "/welcome-screen";

  @override
  Widget build(BuildContext context) {
    final appBar = AppBar(
      backgroundColor: PrimaryColor,
      title: Text(
        AppName,
        style: getAppBarTextStyle(),
      ),
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
