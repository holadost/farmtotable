import 'package:flutter/material.dart';

class RegisterBidWidget extends StatefulWidget {

  @override
  _RegisterBidWidgetState createState() => _RegisterBidWidgetState();
}

class _RegisterBidWidgetState extends State<RegisterBidWidget> {
  final _qtyController = TextEditingController();
  final _amountController = TextEditingController();

  void _submitData() {
    final qty = _qtyController.text;
    final amount = double.parse(_amountController.text);
    print("Quantity: $qty, Price: $amount");
    Navigator.of(context).pop();
  }

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: Card(
          elevation: 5,
          child: Container(
            padding: EdgeInsets.only(
                top: 10,
                left: 10,
                right: 10,
                bottom: MediaQuery.of(context).viewInsets.bottom + 10),
            child:Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                TextField(
                  decoration: InputDecoration(
                      labelText: "Quantity"
                  ),
                  controller: _qtyController,
                  keyboardType: TextInputType.number,
                  onSubmitted: (_) => _submitData(),
                ),
                TextField(
                  decoration: InputDecoration(
                      labelText: "Amount"
                  ),
                  controller: _amountController,
                  keyboardType: TextInputType.number,
                  onSubmitted: (_) => _submitData(),
                ),
                ElevatedButton(
                  onPressed: _submitData,
                  child: const Text("Bid now!"),
                )
              ],
            ),
          )
      ),
    );
  }
}