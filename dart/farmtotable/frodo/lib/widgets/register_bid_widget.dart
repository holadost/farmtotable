import 'package:flutter/material.dart';
import 'package:frodo/net/rest_api_client.dart';

import '../models/item.dart';

class RegisterBidWidget extends StatefulWidget {
  final Item item;

  RegisterBidWidget(this.item);

  @override
  _RegisterBidWidgetState createState() => _RegisterBidWidgetState();
}

class _RegisterBidWidgetState extends State<RegisterBidWidget> {
  final _qtyController = TextEditingController();
  final _amountController = TextEditingController();
  bool _isBeingSubmitted = false;
  var apiClient = RestApiClient();

  void _submitData() async {
    if (_qtyController.text == "" || _amountController.text == "") {
      print("error. Did not submit data");
      return;
    }
    final qty = int.parse(_qtyController.text);
    final amount = double.parse(_amountController.text);
    setState(() {
      _isBeingSubmitted = true;
    });
    try {
      await apiClient.registerBid(widget.item.itemID, amount, qty);
      Scaffold.of(context).showSnackBar(SnackBar(
        content: Text("Registered bid successfully"),
        duration: Duration(seconds: 2),
        backgroundColor: Colors.green,
      ));
    } catch (error) {
      showDialog(
          context: context,
          builder: (ctx) => AlertDialog(
                title: Text(
                  "Bid failed",
                  textAlign: TextAlign.center,
                  style: TextStyle(color: Colors.red),
                ),
                content: Text(
                  "The bid was invalid. Ensure that the current bid is greater than your previous bid",
                ),
                actions: [
                  FlatButton(
                      onPressed: () {
                        Navigator.of(context).pop();
                      },
                      child: const Text("OK"))
                ],
              ));
    }
    setState(() {
      _isBeingSubmitted = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    return _isBeingSubmitted
        ? Container(
            height: 300, child: Center(child: CircularProgressIndicator()))
        : SingleChildScrollView(
            child: Card(
                elevation: 5,
                child: Container(
                  padding: EdgeInsets.only(
                      top: 10,
                      left: 10,
                      right: 10,
                      bottom: MediaQuery.of(context).viewInsets.bottom + 10),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      TextField(
                        decoration: InputDecoration(labelText: "Quantity"),
                        controller: _qtyController,
                        keyboardType: TextInputType.number,
                        onSubmitted: (_) => _submitData(),
                      ),
                      TextField(
                        decoration: InputDecoration(labelText: "Amount"),
                        controller: _amountController,
                        keyboardType: TextInputType.number,
                        onSubmitted: (_) => _submitData(),
                      ),
                      SizedBox(
                        height: 20,
                      ),
                      ElevatedButton(
                        onPressed: () => _submitData(),
                        child: const Text("Bid now!"),
                      )
                    ],
                  ),
                )),
          );
  }
}
