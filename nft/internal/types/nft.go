package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/saisunkari19/modules/nft/common"
	"sort"
	"strings"

	"github.com/saisunkari19/modules/nft/exported"
)

var _ exported.NFT = (*BaseNFT)(nil)

type BaseNFT struct {
	ID           string               `json:"id,omitempty"`
	PrimaryNFTID string               `json:"primary_nftid,omitempty"`
	Type         string               `json:"type"`
	AssetID      string               `json:"asset_id"`
	FileHash     string               `json:"file_hash,omitempty"`
	Category     string               `json:"category"`
	Owner        sdk.AccAddress       `json:"owner"`
	Rights       common.RightsDetails `json:"rights"`
	TokenURI     string               `json:"token_uri"`
}

func NewBaseNFT(id, primaryNFTID, _type, assetID, fileHash, category string, owner sdk.AccAddress,
	rights common.RightsDetails, tokenURI string) BaseNFT {
	return BaseNFT{
		ID:           id,
		PrimaryNFTID: primaryNFTID,
		Type:         _type,
		AssetID:      assetID,
		FileHash:     fileHash,
		Category:     category,
		Owner:        owner,
		Rights:       rights,
		TokenURI:     strings.TrimSpace(tokenURI),
	}
}

func (bnft BaseNFT) GetID() string { return bnft.ID }

func (bnft BaseNFT) GetPrimaryNFTID() string { return bnft.PrimaryNFTID }

func (bnft BaseNFT) GetType() string { return bnft.Type }

func (bnft BaseNFT) GetAssetID() string { return bnft.AssetID }

func (bnft BaseNFT) GetFileHash() string { return bnft.FileHash }

func (bnft BaseNFT) GetCategory() string { return bnft.Category }

func (bnft BaseNFT) GetOwner() sdk.AccAddress { return bnft.Owner }

func (bnft BaseNFT) GetRights() common.RightsDetails { return bnft.Rights }

func (bnft *BaseNFT) SetOwner(address sdk.AccAddress) {
	bnft.Owner = address
}

func (bnft *BaseNFT) SetRights(rights common.RightsDetails) {
	bnft.Rights = rights
}

func (bnft *BaseNFT) SetAssetID(id string) {
	bnft.AssetID = id
}

// GetTokenURI returns the path to optional extra properties
func (bnft BaseNFT) GetTokenURI() string { return bnft.TokenURI }

// EditMetadata edits metadata of an nft
func (bnft *BaseNFT) EditMetadata(tokenURI string) {
	bnft.TokenURI = tokenURI
}

func (bnft BaseNFT) String() string {
	return fmt.Sprintf(`ID:				%s
Owner:			%s
TokenURI:		%s`,
		bnft.ID,
		bnft.Owner,
		bnft.TokenURI,
	)
}

// ----------------------------------------------------------------------------
// NFT

// NFTs define a list of NFT
type NFTs []exported.NFT

// NewNFTs creates a new set of NFTs
func NewNFTs(nfts ...exported.NFT) NFTs {
	if len(nfts) == 0 {
		return NFTs{}
	}
	return NFTs(nfts).Sort()
}

// Append appends two sets of NFTs
func (nfts NFTs) Append(nftsB ...exported.NFT) NFTs {
	return append(nfts, nftsB...).Sort()
}

// Find returns the searched collection from the set
func (nfts NFTs) Find(id string) (nft exported.NFT, found bool) {
	index := nfts.find(id)
	if index == -1 {
		return nft, false
	}
	return nfts[index], true
}

// Update removes and replaces an NFT from the set
func (nfts NFTs) Update(id string, nft exported.NFT) (NFTs, bool) {
	index := nfts.find(id)
	if index == -1 {
		return nfts, false
	}

	return append(append(nfts[:index], nft), nfts[index+1:]...), true
}

// Remove removes an NFT from the set of NFTs
func (nfts NFTs) Remove(id string) (NFTs, bool) {
	index := nfts.find(id)
	if index == -1 {
		return nfts, false
	}

	return append(nfts[:index], nfts[index+1:]...), true
}

// String follows stringer interface
func (nfts NFTs) String() string {
	if len(nfts) == 0 {
		return ""
	}

	out := ""
	for _, nft := range nfts {
		out += fmt.Sprintf("%v\n", nft.String())
	}
	return out[:len(out)-1]
}

// Empty returns true if there are no NFTs and false otherwise.
func (nfts NFTs) Empty() bool {
	return len(nfts) == 0
}

func (nfts NFTs) find(id string) int {
	return FindUtil(nfts, id)
}

// ----------------------------------------------------------------------------
// Encoding

// NFTJSON is the exported NFT format for clients
type NFTJSON map[string]BaseNFT

// MarshalJSON for NFTs
func (nfts NFTs) MarshalJSON() ([]byte, error) {
	nftJSON := make(NFTJSON)
	for _, nft := range nfts {
		id := nft.GetID()
		bnft := NewBaseNFT(id, nft.GetPrimaryNFTID(), nft.GetType(), nft.GetAssetID(), nft.GetFileHash(),
			nft.GetCategory(), nft.GetOwner(), nft.GetRights(), nft.GetTokenURI())
		nftJSON[id] = bnft
	}
	return json.Marshal(nftJSON)
}

// UnmarshalJSON for NFTs
func (nfts *NFTs) UnmarshalJSON(b []byte) error {
	nftJSON := make(NFTJSON)
	if err := json.Unmarshal(b, &nftJSON); err != nil {
		return err
	}

	for id, nft := range nftJSON {
		bnft := NewBaseNFT(id, nft.GetPrimaryNFTID(), nft.GetType(), nft.GetAssetID(), nft.GetFileHash(),
			nft.GetCategory(), nft.GetOwner(), nft.GetRights(), nft.GetTokenURI())
		*nfts = append(*nfts, &bnft)
	}
	return nil
}

// Findable and Sort interfaces
func (nfts NFTs) ElAtIndex(index int) string { return nfts[index].GetID() }
func (nfts NFTs) Len() int                   { return len(nfts) }
func (nfts NFTs) Less(i, j int) bool         { return strings.Compare(nfts[i].GetID(), nfts[j].GetID()) == -1 }
func (nfts NFTs) Swap(i, j int)              { nfts[i], nfts[j] = nfts[j], nfts[i] }

var _ sort.Interface = NFTs{}

// Sort is a helper function to sort the set of coins in place
func (nfts NFTs) Sort() NFTs {
	sort.Sort(nfts)
	return nfts
}
