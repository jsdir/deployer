package names

import (
	"math/rand"
)

func NewRandomName(joinString string) string {
	word1 := adjectives[rand.Intn(len(adjectives))]
	word2 := nouns[rand.Intn(len(nouns))]
	return word1 + joinString + word2
}

var adjectives = []string{
	"adventurous",
	"adorable",
	"aggressive",
	"agile",
	"agitated",
	"animated",
	"beautiful",
	"bubbly",
	"carefree",
	"charming",
	"cheery",
	"clever",
	"clueless",
	"clumsy",
	"courageous",
	"crazy",
	"cuddly",
	"demanding",
	"determined",
	"dizzy",
	"dopey",
	"elaborate",
	"elegant",
	"energetic",
	"excitable",
	"fabulous",
	"fancy",
	"fearless",
	"feisty",
	"friendly",
	"funky",
	"gorgeous",
	"graceful",
	"gruesome",
	"hopeful",
	"irritating",
	"kindhearted",
	"klutzy",
	"knowledgeable",
	"lazy",
	"likable",
	"likely",
	"lively",
	"loud",
	"lovable",
	"lucky",
	"majestic",
	"mysterious",
	"noisy",
	"optimistic",
	"orderly",
	"outgoing",
	"playful",
	"puzzled",
	"sarcastic",
	"scary",
	"scholarly",
	"slimy",
	"stunning",
	"talkative",
	"unlucky",
	"weird",
	"wicked",
	"wild",
}

var nouns = []string{
	"aardvark",
	"albatross",
	"alligator",
	"alpaca",
	"ant",
	"anteater",
	"antelope",
	"ape",
	"armadillo",
	"baboon",
	"badger",
	"baracida",
	"bat",
	"bear",
	"beaver",
	"bee",
	"bison",
	"boar",
	"buffalo",
	"bunny",
	"butterfly",
	"calf",
	"camel",
	"capybara",
	"caribou",
	"cat",
	"caterpilla",
	"cheetah",
	"chick",
	"chicken",
	"chinchilla",
	"clam",
	"cockroach",
	"cormorant",
	"cow",
	"coyote",
	"crab",
	"crane",
	"crocodile",
	"crow",
	"cub",
	"cygnet",
	"deer",
	"dinosaur",
	"dogfish",
	"dolphin",
	"donkey",
	"dove",
	"dragon",
	"dragonfly",
	"duck",
	"duckling",
	"eel",
	"elephant",
	"elk",
	"emu",
	"falcon",
	"ferret",
	"finch",
	"fish",
	"flamingo",
	"fledgling",
	"fly",
	"foal",
	"fox",
	"frog",
	"gazelle",
	"gerbil",
	"gibbon",
	"giraffe",
	"gnu",
	"goat",
	"goldfinch",
	"goldfish",
	"goose",
	"gorilla",
	"gorilla",
	"gosling",
	"grasshopper",
	"grouse",
	"gull",
	"hamster",
	"hatchling",
	"hawk",
	"hedgehog",
	"heron",
	"herring",
	"hippo",
	"hornet",
	"horse",
	"hummingbird",
	"hyena",
	"jackal",
	"jaguar",
	"jay",
	"jellyfish",
	"kangaroo",
	"kingfisher",
	"kitten",
	"koala",
	"kookabura",
	"laark",
	"lamb",
	"lapwing",
	"larva",
	"lemur",
	"leopard",
	"lion",
	"llama",
	"lobster",
	"locust",
	"louse",
	"maggot",
	"magpie",
	"mallard",
	"manatee",
	"mantis",
	"marten",
	"meerkat",
	"mink",
	"mole",
	"mongoose",
	"monkey",
	"moose",
	"mosquito",
	"mouse",
	"mule",
	"narwhal",
	"newt",
	"nightingale",
	"ocelot",
	"octopus",
	"oryx",
	"ostrich",
	"otter",
	"oyster",
	"panda",
	"panther",
	"parrot",
	"partridge",
	"peafowl",
	"pelican",
	"penguin",
	"pheasant",
	"pig",
	"pigeon",
	"piglet",
	"pony",
	"porcupine",
	"porpoise",
	"puppy",
	"quail",
	"quetzal",
	"rabbit",
	"raccoon",
	"ram",
	"rat",
	"raven",
	"reindeer",
	"rhino",
	"rook",
	"salamander",
	"salmon",
	"sandpiper",
	"sardine",
	"scorpion",
	"seahorse",
	"seal",
	"shark",
	"sheep",
	"shrew",
	"skunk",
	"snail",
	"snake",
	"sparrow",
	"spider",
	"spoonbill",
	"squid",
	"squirrel",
	"starling",
	"stingray",
	"stork",
	"swallow",
	"swan",
	"tadpole",
	"tapier",
	"termite",
	"tiger",
	"toad",
	"tortoise",
	"trout",
	"turkey",
	"turtle",
	"unicorn",
	"viper",
	"vulture",
	"wallaby",
	"walrus",
	"wasp",
	"weasel",
	"whale",
	"wildcat",
	"wolf",
	"wolverine",
	"wombat",
	"woodcock",
	"woodpecker",
	"worm",
	"wren",
	"yak",
	"zebra",
	"cyclops",
	"gremlin",
	"sphynx",
	"centaur",
	"dog",
	"bull",
	"bee",
	"sloth",
	"axolotl",
	"tarsier",
}
