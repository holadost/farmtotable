import 'package:flutter/material.dart';

TextStyle getAppBarTextStyle() {
  return TextStyle(
      fontFamily: 'Lato',
      fontWeight: FontWeight.bold,
      color: Colors.white
  );
}

ThemeData getAppTheme() {
  return ThemeData(
    brightness: Brightness.dark,
    primarySwatch: Colors.deepOrange,
    primaryColor: Colors.deepOrange);
}