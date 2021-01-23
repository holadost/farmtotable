import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:intl/intl.dart';
import 'package:firebase_auth/firebase_auth.dart';

import '../util/logging.dart';

class AuthProvider with ChangeNotifier {
  final FirebaseAuth _auth = FirebaseAuth.instance;
  String _idToken = "";
  DateTime _expTime;

  String get token {
    return _idToken;
  }

  bool isAuthorized() {
    if (_idToken == "" || _idToken == null) {
      return false;
    }
    if (_expTime == null || _expTime.isBefore(DateTime.now())) {
      return false;
    }
    return true;
  }

  Future<bool> login(String userEmail, String password) async {
    try {
      info("Logging in user: $userEmail");
      final result = await _auth.signInWithEmailAndPassword(
          email: userEmail, password: password);
      final user = _auth.currentUser;
      final tokenRes = await user.getIdTokenResult(true);
      _idToken = tokenRes.token;
      _expTime = tokenRes.expirationTime;
      notifyListeners();
      info("Successfully logged in. User ID: ${user.uid}."
          "Expiry Time: ${DateFormat.yMMMMd().add_jm().format(_expTime)}");
      return true;
    } catch (error) {
      info("Caught error while logging in: ${error.toString()}");
      return false;
    }
  }

  Future<bool> signup(String userEmail, String password) async {
    try {
      info("Signing up user: $userEmail");
      await _auth.createUserWithEmailAndPassword(
          email: userEmail, password: password);
      info("Successfully signed up new user! UID: "
           "${_auth.currentUser.uid}");
      final res = await login(userEmail, password);
      if (!res) {
        return false;
      }
      return true;
    } catch (error) {
      info("Caught error while signing up: ${error.toString()}");
     return false;
    }
  }
}
