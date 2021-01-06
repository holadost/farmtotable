import '../models/auction_item.dart';

var DUMMY_AUCTIONS = <AuctionItem>[
  AuctionItem(
      itemDescription:
          "Rice is the seed of the grass species Oryza glaberrima (African rice) or Oryza sativa (Asian rice). As a cereal grain, it is the most widely consumed staple food for a large part of the world's human population, especially in Asia and Africa. It is the agricultural commodity with the third-highest worldwide production (rice, 741.5 million tonnes in 2014), after sugarcane (1.9 billion tonnes) and maize (1.0 billion tonnes).",
      auctionID: 1,
      itemID: "Item1",
      itemName: "Rice grains",
      itemQty: 300,
      auctionDurationSecs: Duration(seconds: 3600),
      auctionStartTime: DateTime.now(),
      minBid: 10.00,
      maxBid: 10.00,
      imageURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/7/7b/White%2C_Brown%2C_Red_%26_Wild_rice.jpg/800px-White%2C_Brown%2C_Red_%26_Wild_rice.jpg"),
  AuctionItem(
      itemDescription:
          "Wheat is a grass widely cultivated for its seed, a cereal grain which is a worldwide staple food.[2][3][4] The many species of wheat together make up the genus Triticum; the most widely grown is common wheat (T. aestivum). The archaeological record suggests that wheat was first cultivated in the regions of the Fertile Crescent around 9600 BCE. Botanically, the wheat kernel is a type of fruit called a caryopsis.",
      auctionID: 2,
      itemID: "Item2",
      itemName: "Whole Wheat",
      itemQty: 200,
      auctionDurationSecs: Duration(seconds: 3600),
      auctionStartTime: DateTime.now(),
      minBid: 15.00,
      maxBid: 17.00,
      imageURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b4/Wheat_close-up.JPG/800px-Wheat_close-up.JPG"),
  AuctionItem(
      itemDescription:
          "The pea is most commonly the small spherical seed or the seed-pod of the pod fruit Pisum sativum. Each pod contains several peas, which can be green or yellow. Botanically, pea pods are fruit,[2] since they contain seeds and develop from the ovary of a (pea) flower. The name is also used to describe other edible seeds from the Fabaceae such as the pigeon pea (Cajanus cajan), the cowpea (Vigna unguiculata), and the seeds from several species of Lathyrus.",
      auctionID: 3,
      itemID: "Item3",
      itemName: "Peas and beans",
      itemQty: 300,
      auctionDurationSecs: Duration(seconds: 3600),
      auctionStartTime: DateTime.now(),
      minBid: 22.00,
      maxBid: 22.00,
      imageURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/1/11/Peas_in_pods_-_Studio.jpg/800px-Peas_in_pods_-_Studio.jpg"),
  AuctionItem(
      itemDescription:
          "The carrot (Daucus carota subsp. sativus) is a root vegetable, usually orange in color, though purple, black, red, white, and yellow cultivars exist.[2][3][4] They are a domesticated form of the wild carrot, Daucus carota, native to Europe and Southwestern Asia. The plant probably originated in Persia and was originally cultivated for its leaves and seeds. The most commonly eaten part of the plant is the taproot, although the stems and leaves are also eaten. The domestic carrot has been selectively bred for its greatly enlarged, more palatable, less woody-textured taproot.",
      auctionID: 4,
      itemID: "Item4",
      itemName: "Carrots",
      itemQty: 400,
      auctionDurationSecs: Duration(seconds: 3600),
      auctionStartTime: DateTime.now(),
      minBid: 14.00,
      maxBid: 15.25,
      imageURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/c/c3/Carrots_at_Ljubljana_Central_Market.JPG/1024px-Carrots_at_Ljubljana_Central_Market.JPG"),
  AuctionItem(
      itemDescription:
          "Vegetables are parts of plants that are consumed by humans or other animals as food. The original meaning is still commonly used and is applied to plants collectively to refer to all edible plant matter, including the flowers, fruits, stems, leaves, roots, and seeds. The alternate definition of the term is applied somewhat arbitrarily, often by culinary and cultural tradition. It may exclude foods derived from some plants that are fruits, flowers, nuts, and cereal grains, but include savoury fruits such as tomatoes and courgettes, flowers such as broccoli, and seeds such as pulses.",
      auctionID: 5,
      itemID: "Item5",
      itemName: "Vegetables",
      itemQty: 500,
      auctionDurationSecs: Duration(seconds: 3600),
      auctionStartTime: DateTime.now(),
      minBid: 20.00,
      maxBid: 32.33,
      imageURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/2/24/Marketvegetables.jpg/800px-Marketvegetables.jpg"),
];
