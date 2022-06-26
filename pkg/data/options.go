package data

type Options struct {
	Handshaker string  `short:"H" long:"handshaker" description:"Handshaker service address to use for search/await handshaking" default:"gamu.group:20220"`
	Host       string  `short:"h" long:"host" description:"Directly hosting on an address"`
	Join       string  `short:"j" long:"join" description:"Directly join an address"`
	Search     string  `short:"s" long:"search" description:"Search for a given user using external handshaking"`
	Await      bool    `short:"a" long:"await" description:"Await for a player search"`
	Map        string  `short:"m" long:"map" description:"Map to start the game on" default:"001"`
	Name       string  `short:"n" long:"name" description:"Name to user in multiplayer"`
	Speed      float64 `short:"S" long:"speed" description:"Game speed multiplier" default:"1.0"`
}
