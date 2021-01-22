import 'package:flutter/foundation.dart';
import 'package:provider/provider.dart';

import '../util/logging.dart';

class AuthProvider with ChangeNotifier {
  String _idToken = "";
  String _userID = "";
  String _userName = "";
  String _userEmailAddress = "";

  String get token {
    return _idToken;
  }

  Future<bool> login(String userEmail, String password) async {
    info("Successfully logged in");
    notifyListeners();
    return true;
  }

  Future<bool> signup(String userEmail, String password) async {
    info("Successfully signed up");
    notifyListeners();
    return true;
  }
}
