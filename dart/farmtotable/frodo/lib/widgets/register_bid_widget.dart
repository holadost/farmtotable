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

  void _showAlert(String title, String content) {
    showDialog(
        context: context,
        builder: (ctx) => AlertDialog(
              title: Text(
                title,
                textAlign: TextAlign.center,
                style: TextStyle(color: Colors.red),
              ),
              content: Text(content),
              actions: [
                FlatButton(
                    onPressed: () {
                      Navigator.of(context).pop();
                    },
                    child: const Text("OK"))
              ],
            ));
  }

  void _submitData() async {
    if (_qtyController.text == "" || _amountController.text == "") {
      _showAlert("Invalid inputs", "Enter an amount and quantity");
      return;
    }

    // Check quantity.
    int qty;
    try {
      qty = int.parse(_qtyController.text);
    } catch (error) {
      _showAlert("Invalid quanity", "Quantity must be a whole number");
      return;
    }
    if (qty < widget.item.minBidQty ||
        qty > widget.item.maxBidQty ||
        qty > widget.item.itemQty ||
        qty % widget.item.minBidQty != 0) {
      _showAlert(
          "Invalid bid quantity",
          "The bid quantity must be >= than "
              "${widget.item.minBidQty}${widget.item.itemUnit} "
              "and must be lesser than ${widget.item.maxBidQty}${widget.item.itemUnit}"
              "and must be a multiple of ${widget.item.minBidQty}${widget.item.itemUnit}");
      return;
    }

    if (qty < widget.item.minBidQty ||
        qty > widget.item.maxBidQty ||
        qty > widget.item.itemQty ||
        qty % widget.item.minBidQty != 0) {
      _showAlert(
          "Invalid bid quantity",
          "The bid quantity must be >= than "
              "${widget.item.minBidQty}${widget.item.itemUnit} "
              "and must be lesser than ${widget.item.maxBidQty}${widget.item.itemUnit}"
              "and must be a multiple of ${widget.item.minBidQty}${widget.item.itemUnit}");
      return;
    }

    // Check amount.
    double amount;
    try {
      amount = double.parse(_amountController.text);
    } catch (error) {
      _showAlert("Invalid amount", "Amount must be a non-zero number");
      return;
    }
    if ((qty <= 0) || (amount <= 0)) {
      _showAlert("Invalid amount/quantity", "Quantity and Amount must be > 0");
      return;
    }

    // Set state so that user can see a loading spinner while we process
    // the request.
    setState(() {
      _isBeingSubmitted = true;
    });

    try {
      await apiClient.registerBid(widget.item.itemID, amount, qty);
      Scaffold.of(context).showSnackBar(SnackBar(
        content: Text(
          "Registered bid successfully",
          textAlign: TextAlign.center,
          style: TextStyle(color: Colors.black, fontSize: 16),
        ),
        duration: Duration(seconds: 2),
        backgroundColor: Colors.green,
      ));
    } catch (error) {
      _showAlert(
          "Bid failed",
          "The bid was invalid. Ensure that the current bid is greater than "
              "your previous bid");
    }
    // Set state so that we can disable the loading spinner and user can
    // see the response.
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
                      Container(
                          height: 250,
                          width: double.infinity,
                          decoration: BoxDecoration(
                              shape: BoxShape.rectangle,
                              image: DecorationImage(
                                  fit: BoxFit.fill, image: NetworkImage(widget.item.imageURL)))),
                      SizedBox(
                        height: 20,
                      ),
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
