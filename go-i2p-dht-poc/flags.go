package main

import (
	"flag"
)

var (
	/* INERT SO FAR
	      	saveFile = flag.Bool("t", false,
	   		"Use saved key files and persist tunnels(If false, tunnel will not persist after program is stopped.")
	   	encryptLeaseSet = flag.Bool("l", true,
	   		"Use an encrypted leaseset(true or false)")
	   	encryptKeyFiles = flag.String("cr", "",
	   		"Encrypt/decrypt the key files with a passfile")
	*/
	debug = flag.Bool("d", false,
		"enable debug mode")
)

var (
	participatingPeerCount = flag.Int("p", 2,
		"Set base number of participating peers for providers (0 to 10)")
)

var (
	/* INERT FOR NOW
	   	inAllowZeroHop = flag.Bool("zi", false,
	   		"Allow zero-hop, non-anonymous tunnels in(true or false)")
	   	outAllowZeroHop = flag.Bool("zo", false,
	   		"Allow zero-hop, non-anonymous tunnels out(true or false)")
	   	useCompression = flag.Bool("z", false,
	   		"Uze gzip(true or false)")
	       samHost = flag.String("sh", "127.0.0.1",
	   		"SAM host")
	   	samPort = flag.String("sp", "7656",
	   		"SAM port")
	*/
	dataID = flag.String("n", "i2pdht",
		"Tunnel name, this must be unique but can be anything.")

/* INERT FOR NOW
accessListType = flag.String("a", "none",
	"Type of access list to use, can be \"whitelist\" \"blacklist\" or \"none\".")
inLength = flag.Int("il", 3,
	"Set inbound tunnel length(0 to 7)")
outLength = flag.Int("ol", 3,
	"Set outbound tunnel length(0 to 7)")
inQuantity = flag.Int("iq", 6,
	"Set inbound tunnel quantity(0 to 15)")
outQuantity = flag.Int("oq", 6,
	"Set outbound tunnel quantity(0 to 15)")
inVariance = flag.Int("iv", 0,
	"Set inbound tunnel length variance(-7 to 7)")
outVariance = flag.Int("ov", 0,
	"Set outbound tunnel length variance(-7 to 7)")
inBackupQuantity = flag.Int("ib", 2,
	"Set inbound tunnel backup quantity(0 to 5)")
outBackupQuantity = flag.Int("ob", 2,
	"Set outbound tunnel backup quantity(0 to 5)")
leaseSetKey = flag.String("k", "none",
	"key for encrypted leaseset")
leaseSetPrivateKey = flag.String("pk", "none",
	"private key for encrypted leaseset")
leaseSetPrivateSigningKey = flag.String("psk", "none",
	"private signing key for encrypted leaseset")
*/
)

//var (
/*
	webAdmin = flag.Bool("w", true,
		"Start web administration interface")
	sigType = flag.String("st", "",
		"Signature type")
	webPort = flag.String("wp", "7957",
		"Web port")
	webUser = flag.String("webuser", "samcatd",
		"Web interface username")
	webPass = flag.String("webpass", "",
		"Web interface password")
	webCSS = flag.String("css", "css/styles.css",
		"custom CSS for web interface")
	webJS = flag.String("js", "js/scripts.js",
		"custom JS for web interface")
	targetDir = flag.String("d", "",
		"Directory to save tunnel configuration file in.")
	targetDest = flag.String("de", "",
		"Destination to connect client's to by default.")
	iniFile = flag.String("f", "none",
		"Use an ini file for configuration(config file options override passed arguments for now.)")
	targetDestination = flag.String("i", "none",
        "Destination for queries. Invalid for providers.")
	targetHost = flag.String("h", "127.0.0.1",
		"Target host(Host of service to forward to i2p)")
	targetPort = flag.String("p", "57890",
		"Target port(Port of service to forward to i2p)")
*/
//)
