import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_auth/firebase_auth.dart';

import '../util/logging.dart';

class AuthProvider with ChangeNotifier {
  final FirebaseAuth _auth = FirebaseAuth.instance;
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
    try {
      final result = await _auth.createUserWithEmailAndPassword(
          email: userEmail, password: password);
      info("Signed up!");
      return true;
    } catch (error) {
      info("Caught error: ${error.toString()}");
     return false;
    }
  }
}
