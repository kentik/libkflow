package flow

import (
	"fmt"

	capnp "zombiezen.com/go/capnproto2"

	"github.com/kentik/libkflow/chf"
)

// ToCapnProtoMessage converts a slice of Flow objects into a Cap'n Proto message.
func ToCapnProtoMessage(flows []Flow, segment *capnp.Segment) (*capnp.Message, error) {
	packedCHF, err := chf.NewRootPackedCHF(segment)
	if err != nil {
		return nil, fmt.Errorf("failed to create root packed CHF: %w", err)
	}

	chfList, err := packedCHF.NewMsgs(int32(len(flows)))
	if err != nil {
		return nil, fmt.Errorf("failed to create messages list: %w", err)
	}

	for i, f := range flows {
		var list chf.Custom_List
		if n := int32(len(f.Customs)); n > 0 {
			list, err = chf.NewCustom_List(segment, n)
			if err != nil {
				return nil, fmt.Errorf("failed to create custom list: %w", err)
			}
		}
		f.FillCHF(chfList.At(i), list)
	}

	err = packedCHF.SetMsgs(chfList)
	if err != nil {
		return nil, fmt.Errorf("failed to set messages: %w", err)
	}
	return segment.Message(), nil
}
