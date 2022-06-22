package data

type Options struct {
	Handshaker string `short:"H" long:"handshaker" description:"Handshaker service address to use for search/await handshaking"`
	Host       string `short:"h" long:"host" description:"Directly hosting on an address"`
	Join       string `short:"j" long:"join" description:"Directly join an address"`
	Search     string `short:"s" long:"search" description:"Search for a given user using external handshaking"`
	Await      string `short:"a" long:"await" description:"Await for a player search using the given user in external handshaking"`
	Map        string `short:"m" long:"map" description:"Map to start the game on" default:"001"`
}
