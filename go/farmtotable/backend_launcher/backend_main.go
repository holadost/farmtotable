package main

import (
	"farmtotable/aragorn"
	"farmtotable/gandalf"
	"farmtotable/legolas"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"strconv"
	"time"
)

/* Backend backend_launcher. This can be used for dev/test purposes when we want both
Aragorn and Legolas running as the same service allowing us to use Sqlite for
Gandalf's backend. THIS MUST NOT BE USED IN PRODUCTION. */
func main() {
	flag.Parse()
	g := gandalf.NewSqliteGandalf()
	prepareDB(g)
	go aragorn.NewAragornWithGandalf(g).Run()
	go legolas.NewLegolasWithGandalf(g).Run()
	// Block forever
	select {}
}

func prepareDB(gnd *gandalf.Gandalf) {
	glog.Info("Preparing backend")
	addDummyUsers(gnd, "user", 5)
	addDummySuppliers(gnd, "supplier", 3)
	addDummyItems(gnd)
	addDummyOrders(gnd)
}

func addDummyUsers(gnd *gandalf.Gandalf, userIDPrefix string, numUsers uint) {
	glog.Info("Adding dummy users")
	// Add users.
	for ii := 0; ii < int(numUsers); ii++ {
		userID := userIDPrefix + strconv.Itoa(ii)
		name := "Nikhil_" + strconv.Itoa(ii)
		emailID := fmt.Sprintf("kjahd@lkaj_%s.com", strconv.Itoa(ii))
		err := gnd.RegisterUser(userID, name, emailID,
			"9873981799", "khadkjhadkha")
		if err != nil {
			glog.Fatalf("Unable to register user")
		}
	}
}

func addDummySuppliers(gnd *gandalf.Gandalf, supplierPrefix string, numSuppliers uint) {
	glog.Info("Adding dummy suppliers")
	for ii := 0; ii < int(numSuppliers); ii++ {
		name := supplierPrefix + strconv.Itoa(ii)
		emailID := fmt.Sprintf("nikhil.%s@gmail.com", strconv.Itoa(ii))
		err := gnd.RegisterSupplier(name, emailID,
			"9198029973", "Mera Naam Joker",
			"This supplier is a god amongst humans",
			"tag1,tag2,tag3")
		if err != nil {
			glog.Fatalf("Unable to register supplier")
		}
	}

}

func addDummyItems(gnd *gandalf.Gandalf) {
	glog.Info("Adding dummy items")
	// This function adds 6 products (2 times each(A and B)) such that
	// 12 products can go on auction everyday starting from now.
	// The auction duration is set to 10 days for all products.
	for ii := 0; ii < 10; ii++ {
		startDate := time.Now().Add(time.Duration(86400 * ii))
		err := gnd.RegisterItem(
			"supplier1",
			"Rice-A "+strconv.Itoa(ii),
			"Rice is the seed of the grass species Oryza glaberrima (African rice) or Oryza sativa (Asian rice). As a cereal grain, it is the most widely consumed staple food for a large part of the world's human population, especially in Asia and Africa. It is the agricultural commodity with the third-highest worldwide production (rice, 741.5 million tonnes in 2014), after sugarcane (1.9 billion tonnes) and maize (1.0 billion tonnes).",
			100,
			startDate,
			15.00,
			10*86400,
			"https://upload.wikimedia.org/wikipedia/commons/thumb/7/7b/White%2C_Brown%2C_Red_%26_Wild_rice.jpg/800px-White%2C_Brown%2C_Red_%26_Wild_rice.jpg",
			1,
			50,
			"kg")
		if err != nil {
			glog.Fatalf("Unable to register item 1")
		}

		err = gnd.RegisterItem(
			"supplier1",
			"Wheat A-"+strconv.Itoa(ii),
			"Rice is the seed of the grass species Oryza glaberrima (African rice) or Oryza sativa (Asian rice). As a cereal grain, it is the most widely consumed staple food for a large part of the world's human population, especially in Asia and Africa. It is the agricultural commodity with the third-highest worldwide production (rice, 741.5 million tonnes in 2014), after sugarcane (1.9 billion tonnes) and maize (1.0 billion tonnes).",
			300,
			startDate,
			20.00,
			10*86400,
			"https://upload.wikimedia.org/wikipedia/commons/thumb/b/b4/Wheat_close-up.JPG/800px-Wheat_close-up.JPG",
			1,
			50,
			"kg")
		if err != nil {
			glog.Fatalf("Unable to register item 1")
		}

		err = gnd.RegisterItem(
			"supplier2",
			"Peas A-"+strconv.Itoa(ii),
			"The pea is most commonly the small spherical seed or the seed-pod of the pod fruit Pisum sativum. Each pod contains several peas, which can be green or yellow. Botanically, pea pods are fruit,[2] since they contain seeds and develop from the ovary of a (pea) flower. The name is also used to describe other edible seeds from the Fabaceae such as the pigeon pea (Cajanus cajan), the cowpea (Vigna unguiculata), and the seeds from several species of Lathyrus.",
			100,
			startDate,
			22.50,
			10*86400,
			"https://upload.wikimedia.org/wikipedia/commons/thumb/1/11/Peas_in_pods_-_Studio.jpg/800px-Peas_in_pods_-_Studio.jpg",
			1,
			10,
			"kg")
		if err != nil {
			glog.Fatalf("Unable to register item 1")
		}

		err = gnd.RegisterItem(
			"supplier2",
			"Carrots A-"+strconv.Itoa(ii),
			"The carrot (Daucus carota subsp. sativus) is a root vegetable, usually orange in color, though purple, black, red, white, and yellow cultivars exist.[2][3][4] They are a domesticated form of the wild carrot, Daucus carota, native to Europe and Southwestern Asia. The plant probably originated in Persia and was originally cultivated for its leaves and seeds. The most commonly eaten part of the plant is the taproot, although the stems and leaves are also eaten. The domestic carrot has been selectively bred for its greatly enlarged, more palatable, less woody-textured taproot.",
			75,
			startDate,
			15.50,
			10*86400,
			"https://upload.wikimedia.org/wikipedia/commons/thumb/c/c3/Carrots_at_Ljubljana_Central_Market.JPG/1024px-Carrots_at_Ljubljana_Central_Market.JPG",
			1,
			20,
			"kg")
		if err != nil {
			glog.Fatalf("Unable to register item 1")
		}

		err = gnd.RegisterItem(
			"supplier3",
			"Quinoa A-"+strconv.Itoa(ii),
			"Quinoa (Chenopodium quinoa; /ˈkiːnwɑː/ or /kɪˈnoʊ.ə/, from Quechua kinwa or kinuwa)[2] is a flowering plant in the amaranth family. It is a herbaceous annual plant grown as a crop primarily for its edible seeds; the seeds are rich in protein, dietary fiber, B vitamins, and dietary minerals in amounts greater than in many grains.[3] Quinoa is not a grass, but rather a pseudocereal botanically related to spinach and amaranth (Amaranthus spp.), and originated in the Andean region of northwestern South America.[4] It was first used to feed livestock 5.2–7.0 thousand years ago, and for human consumption 3–4 thousand years ago in the Lake Titicaca basin of Peru and Bolivia.[5]\n\nToday, almost all production in the Andean region is done by small farms and associations. Its cultivation has spread to more than 70 countries, including Kenya, India, the United States, and several European countries.[6] As a result of increased popularity and consumption in North America, Europe, and Australasia, quinoa crop prices tripled between 2006 and 2013.[7][8]",
			75,
			startDate,
			30.50,
			10*86400,
			"https://upload.wikimedia.org/wikipedia/commons/thumb/4/43/Red_quinoa.png/1024px-Red_quinoa.png",
			1,
			20,
			"kg")
		if err != nil {
			glog.Fatalf("Unable to register item 1")
		}

		err = gnd.RegisterItem(
			"supplier3",
			"Kale A-"+strconv.Itoa(ii),
			"Kale (/keɪl/), or leaf cabbage, belongs to a group of cabbage (Brassica oleracea) cultivars grown for their edible leaves, although some are used as ornamentals. Kale plants have green or purple leaves, and the central leaves do not form a head (as with headed cabbage). Kales are considered to be closer to wild cabbage than most of the many domesticated forms of Brassica oleracea.[1]\n\nKale originated in the eastern Mediterranean and Asia Minor, where it was cultivated for food beginning by 2000 BCE at the latest.[3] Curly-leaved varieties of cabbage already existed along with flat-leaved varieties in Greece in the 4th century BC. These forms, which were referred to by the Romans as Sabellian kale, are considered to be the ancestors of modern kales.\n\nThe earliest record of cabbages in western Europe is of hard-heading cabbage in the 13th century.[3] Records in 14th-century England distinguish between hard-heading cabbage and loose-leaf kale.[3]",
			75,
			startDate,
			12.25,
			10*86400,
			"https://upload.wikimedia.org/wikipedia/commons/thumb/4/4b/Kale-Bundle.jpg/1024px-Kale-Bundle.jpg",
			1,
			20,
			"kg")
		if err != nil {
			glog.Fatalf("Unable to register item 1")
		}

		err = gnd.RegisterItem(
			"supplier1",
			"Rice B-"+strconv.Itoa(ii),
			"Rice is the seed of the grass species Oryza glaberrima (African rice) or Oryza sativa (Asian rice). As a cereal grain, it is the most widely consumed staple food for a large part of the world's human population, especially in Asia and Africa. It is the agricultural commodity with the third-highest worldwide production (rice, 741.5 million tonnes in 2014), after sugarcane (1.9 billion tonnes) and maize (1.0 billion tonnes).",
			100,
			startDate,
			16.00,
			10*86400,
			"https://upload.wikimedia.org/wikipedia/commons/thumb/7/7b/White%2C_Brown%2C_Red_%26_Wild_rice.jpg/800px-White%2C_Brown%2C_Red_%26_Wild_rice.jpg",
			1,
			50,
			"kg")
		if err != nil {
			glog.Fatalf("Unable to register item 1")
		}

		err = gnd.RegisterItem(
			"supplier1",
			"Wheat B-"+strconv.Itoa(ii),
			"Rice is the seed of the grass species Oryza glaberrima (African rice) or Oryza sativa (Asian rice). As a cereal grain, it is the most widely consumed staple food for a large part of the world's human population, especially in Asia and Africa. It is the agricultural commodity with the third-highest worldwide production (rice, 741.5 million tonnes in 2014), after sugarcane (1.9 billion tonnes) and maize (1.0 billion tonnes).",
			300,
			startDate,
			22.00,
			10*86400,
			"https://upload.wikimedia.org/wikipedia/commons/thumb/b/b4/Wheat_close-up.JPG/800px-Wheat_close-up.JPG",
			1,
			50,
			"kg")
		if err != nil {
			glog.Fatalf("Unable to register item 1")
		}

		err = gnd.RegisterItem(
			"supplier2",
			"Peas B-"+strconv.Itoa(ii),
			"The pea is most commonly the small spherical seed or the seed-pod of the pod fruit Pisum sativum. Each pod contains several peas, which can be green or yellow. Botanically, pea pods are fruit,[2] since they contain seeds and develop from the ovary of a (pea) flower. The name is also used to describe other edible seeds from the Fabaceae such as the pigeon pea (Cajanus cajan), the cowpea (Vigna unguiculata), and the seeds from several species of Lathyrus.",
			100,
			startDate,
			25.50,
			10*86400,
			"https://upload.wikimedia.org/wikipedia/commons/thumb/1/11/Peas_in_pods_-_Studio.jpg/800px-Peas_in_pods_-_Studio.jpg",
			1,
			10,
			"kg")
		if err != nil {
			glog.Fatalf("Unable to register item 1")
		}

		err = gnd.RegisterItem(
			"supplier2",
			"Carrots B-"+strconv.Itoa(ii),
			"The carrot (Daucus carota subsp. sativus) is a root vegetable, usually orange in color, though purple, black, red, white, and yellow cultivars exist.[2][3][4] They are a domesticated form of the wild carrot, Daucus carota, native to Europe and Southwestern Asia. The plant probably originated in Persia and was originally cultivated for its leaves and seeds. The most commonly eaten part of the plant is the taproot, although the stems and leaves are also eaten. The domestic carrot has been selectively bred for its greatly enlarged, more palatable, less woody-textured taproot.",
			75,
			startDate,
			12.50,
			10*86400,
			"https://upload.wikimedia.org/wikipedia/commons/thumb/c/c3/Carrots_at_Ljubljana_Central_Market.JPG/1024px-Carrots_at_Ljubljana_Central_Market.JPG",
			1,
			20,
			"kg")
		if err != nil {
			glog.Fatalf("Unable to register item 1")
		}

		err = gnd.RegisterItem(
			"supplier3",
			"Quinoa B-"+strconv.Itoa(ii),
			"Quinoa (Chenopodium quinoa; /ˈkiːnwɑː/ or /kɪˈnoʊ.ə/, from Quechua kinwa or kinuwa)[2] is a flowering plant in the amaranth family. It is a herbaceous annual plant grown as a crop primarily for its edible seeds; the seeds are rich in protein, dietary fiber, B vitamins, and dietary minerals in amounts greater than in many grains.[3] Quinoa is not a grass, but rather a pseudocereal botanically related to spinach and amaranth (Amaranthus spp.), and originated in the Andean region of northwestern South America.[4] It was first used to feed livestock 5.2–7.0 thousand years ago, and for human consumption 3–4 thousand years ago in the Lake Titicaca basin of Peru and Bolivia.[5]\n\nToday, almost all production in the Andean region is done by small farms and associations. Its cultivation has spread to more than 70 countries, including Kenya, India, the United States, and several European countries.[6] As a result of increased popularity and consumption in North America, Europe, and Australasia, quinoa crop prices tripled between 2006 and 2013.[7][8]",
			75,
			startDate,
			28.50,
			10*86400,
			"https://upload.wikimedia.org/wikipedia/commons/thumb/4/43/Red_quinoa.png/1024px-Red_quinoa.png",
			1,
			20,
			"kg")
		if err != nil {
			glog.Fatalf("Unable to register item 1")
		}

		err = gnd.RegisterItem(
			"supplier3",
			"Kale B-"+strconv.Itoa(ii),
			"Kale (/keɪl/), or leaf cabbage, belongs to a group of cabbage (Brassica oleracea) cultivars grown for their edible leaves, although some are used as ornamentals. Kale plants have green or purple leaves, and the central leaves do not form a head (as with headed cabbage). Kales are considered to be closer to wild cabbage than most of the many domesticated forms of Brassica oleracea.[1]\n\nKale originated in the eastern Mediterranean and Asia Minor, where it was cultivated for food beginning by 2000 BCE at the latest.[3] Curly-leaved varieties of cabbage already existed along with flat-leaved varieties in Greece in the 4th century BC. These forms, which were referred to by the Romans as Sabellian kale, are considered to be the ancestors of modern kales.\n\nThe earliest record of cabbages in western Europe is of hard-heading cabbage in the 13th century.[3] Records in 14th-century England distinguish between hard-heading cabbage and loose-leaf kale.[3]",
			75,
			startDate,
			14.25,
			10*86400,
			"https://upload.wikimedia.org/wikipedia/commons/thumb/4/4b/Kale-Bundle.jpg/1024px-Kale-Bundle.jpg",
			1,
			20,
			"kg")
		if err != nil {
			glog.Fatalf("Unable to register item 1")
		}
	}
}

func addDummyOrders(gnd *gandalf.Gandalf) {
	glog.Info("Adding dummy orders")
}
