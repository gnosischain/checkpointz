package specblock

import (
	"errors"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	dynssz "github.com/pk910/dynamic-ssz"
)

type SpecBlock struct {
	*spec.VersionedSignedBeaconBlock

	dynSsz func() (*dynssz.DynSsz, error)
}

func NewSpecBlock(block *spec.VersionedSignedBeaconBlock, dynssz func() (*dynssz.DynSsz, error)) *SpecBlock {
	return &SpecBlock{block, dynssz}
}

func (block *SpecBlock) Root() (phase0.Root, error) {
	ssz, err := block.dynSsz()
	if err != nil {
		return [32]byte{}, err
	}

	switch block.Version {
	case spec.DataVersionPhase0:
		return ssz.HashTreeRoot(block.Phase0.Message)
	case spec.DataVersionAltair:
		return ssz.HashTreeRoot(block.Altair.Message)
	case spec.DataVersionBellatrix:
		return ssz.HashTreeRoot(block.Bellatrix.Message)
	case spec.DataVersionCapella:
		return ssz.HashTreeRoot(block.Capella.Message)
	case spec.DataVersionDeneb:
		return ssz.HashTreeRoot(block.Deneb.Message)
	case spec.DataVersionElectra:
		return ssz.HashTreeRoot(block.Electra.Message)
	default:
		return [32]byte{}, errors.New("unknown version")
	}
}
