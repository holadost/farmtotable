import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../providers/auth_provider.dart';
import '../screens/auth_screen.dart';
import '../screens/welcome_screen.dart';

class HomeScreen extends StatelessWidget {
  static const routeName = "/";

  @override
  Widget build(BuildContext context) {
    final authProv = Provider.of<AuthProvider>(context, listen: true);
    if (authProv.isAuthorized()) {
      return WelcomeScreen();
    }
    return AuthScreen();
  }
}
