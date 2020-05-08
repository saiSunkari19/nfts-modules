package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/saisunkari19/modules/nft/common"
)

// NFT non fungible token interface
type NFT interface {
	GetID() string
	GetPrimaryNFTID() string
	GetType() string
	GetAssetID() string
	GetFileHash() string
	GetCategory() string
	GetRights() common.RightsDetails
	GetOwner() sdk.AccAddress
	SetOwner(address sdk.AccAddress)
	SetRights(rights common.RightsDetails)
	SetAssetID(id string)
	GetTokenURI() string
	EditMetadata(tokenURI string)
	String() string
}
