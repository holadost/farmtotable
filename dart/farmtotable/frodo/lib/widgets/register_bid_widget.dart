import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/item.dart';
import '../providers/aragorn_client_provider.dart';
import '../util/constants.dart';

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
  bool _bidButtonDisabled = false;
  double _totalAmount = 0;
  AragornClientProvider apiClient;

  @override
  void didChangeDependencies() {
    apiClient = Provider.of<AragornClientProvider>(context, listen: false);
    super.didChangeDependencies();
  }

  void _updateTotal() {
    int qty;
    try {
      qty = int.parse(_qtyController.text);
    } catch (error) {
      setState(() {
        _totalAmount = 0;
      });
      return;
    }

    double amount;
    try {
      amount = double.parse(_amountController.text);
    } catch (error) {
      setState(() {
        _totalAmount = 0;
      });
      return;
    }
    if ((qty <= 0) || (amount <= 0)) {
      return;
    }
    setState(() {
      _totalAmount = qty * amount;
    });
  }

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
      _showAlert("Invalid quantity", "Quantity must be a whole number");
      return;
    }
    const alertTitle = "Invalid bid quantity";
    final alertContent = "The bid quantity must be: "
        "\n  1. >= ${widget.item.minBidQty}${widget.item.itemUnit} "
        "\n  2. <= ${widget.item.maxBidQty}${widget.item.itemUnit} "
        "\n  3. A multiple of ${widget.item.minBidQty}${widget.item.itemUnit}";
    if (qty < widget.item.minBidQty ||
        qty > widget.item.maxBidQty ||
        qty > widget.item.itemQty ||
        qty % widget.item.minBidQty != 0) {
      _showAlert(alertTitle, alertContent);
      return;
    }

    if (qty < widget.item.minBidQty ||
        qty > widget.item.maxBidQty ||
        qty > widget.item.itemQty ||
        qty % widget.item.minBidQty != 0) {
      _showAlert(alertTitle, alertContent);
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

    try {
      // Set state so that user can see a loading spinner while we process
      // the request.
      setState(() {
        _isBeingSubmitted = true;
      });
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
      _bidButtonDisabled = true;
    } catch (error) {
      _showAlert(
          "Bid failed",
          "The bid was invalid. Ensure that the current bid is greater than "
              "your previous bid");
      return;
    } finally {
      // Set state so that we can disable the loading spinner and user can
      // see the response.
      setState(() {
        _isBeingSubmitted = false;
      });
    }
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
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      Container(
                          height: 250,
                          width: double.infinity,
                          decoration: BoxDecoration(
                              shape: BoxShape.rectangle,
                              image: DecorationImage(
                                  fit: BoxFit.fill,
                                  image: NetworkImage(widget.item.imageURL)))),
                      SizedBox(
                        height: 20,
                      ),
                      TextField(
                        decoration: InputDecoration(
                          labelText: "Quantity(${widget.item.itemUnit})",
                          labelStyle: TextStyle(fontSize: 20),
                        ),
                        controller: _qtyController,
                        keyboardType: TextInputType.number,
                        onSubmitted: (_) => _submitData(),
                        onChanged: (x) {
                          _updateTotal();
                        },
                      ),
                      TextField(
                        decoration: InputDecoration(
                            labelText:
                                "Amount(per ${widget.item.minBidQty}${widget.item.itemUnit})",
                            labelStyle: TextStyle(fontSize: 20)),
                        controller: _amountController,
                        keyboardType: TextInputType.number,
                        onSubmitted: (_) => _submitData(),
                        onChanged: (x) {
                          _updateTotal();
                        },
                      ),
                      SizedBox(
                        height: 20,
                      ),
                      if (_totalAmount > 0)
                        Text(
                          "Total: $Rupee$_totalAmount",
                          style: TextStyle(
                            color: Colors.green,
                            fontSize: 20,
                            fontFamily: 'Anton'
                          ),
                        ),
                      SizedBox(
                        height: 20,
                      ),
                      RaisedButton(
                        onPressed: _bidButtonDisabled ? null : _submitData,
                        child: const Text("Bid now!"),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(16.0),
                        ),
                      )
                    ],
                  ),
                )),
          );
  }
}
