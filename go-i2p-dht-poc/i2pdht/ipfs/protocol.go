package ipfs

/*
import (
	"encoding/base32"
	"fmt"

	ma "github.com/multiformats/go-multiaddr"
)

var serviceIDEncoding = base32.StdEncoding.WithPadding(base32.NoPadding)
var garlicListenAddr ma.Multiaddr

const GARLIC_LISTEN_PROTO_CODE = 0x55

var garlicListenProto = ma.Protocol{
	"garlicListen", GARLIC_LISTEN_PROTO_CODE, ma.CodeToVarint(GARLIC_LISTEN_PROTO_CODE), 0, false, nil}

func init() {
	// Add the listen protocol
	if err := ma.AddProtocol(garlicListenProto); err != nil {
		panic(fmt.Errorf("Failed adding garlicListen protocol: %v", err))
	} else if garlicListenAddr, err = ma.NewMultiaddr("/garlicListen"); err != nil {
		panic(fmt.Errorf("Failed creating garlicListen addr: %v err))
	}
	// Change existing garlic protocol to support v3 and be more lenient when transcoding
	/*ma.TranscoderGarlic64 = ma.NewTranscoderFromFunctions(garlicStringToBytes, garlicBytesToString, nil)
	for _, p := range ma.Protocols {
		if p.Code == ma.P_GARLIC64 {
			p.Size = ma.LengthPrefixedVarSize
			p.Transcoder = ma.TranscoderGarlic64
			break
		}
	}
}
*/
/*
func garlicStringToBytes(str string) ([]byte, error) {
	// Just convert the whole thing for now
	return []byte(str), nil
}

func garlicBytesToString(byts []byte) (string, error) {
	// Just convert the whole thing for now
	return string(byts), nil
}
*/
