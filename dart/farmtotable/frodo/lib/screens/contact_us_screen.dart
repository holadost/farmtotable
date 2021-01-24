import 'package:flutter/material.dart';

import '../widgets/side_drawer_widget.dart';
import '../util/constants.dart';
import '../util/styles.dart';

class ContactUsScreen extends StatelessWidget {
  static const routeName = "/contact-us-screen";

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
      child: Text('Contact Us'),
    );
    return Scaffold(
      appBar: appBar,
      body: body,
      drawer: SideDrawerWidget(),
    );
  }
}
