import 'package:flutter/foundation.dart';
import 'package:provider/provider.dart';

class RestClientProvider with ChangeNotifier {
  final _userToken;
  RestClientProvider(this._userToken);
}
