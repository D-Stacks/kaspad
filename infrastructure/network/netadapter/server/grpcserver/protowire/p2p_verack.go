package protowire

import (
	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/pkg/errors"
)

func (x *KaspadP2PMessage_Verack) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "KaspadMessage_Verack is nil")
	}
	return &appmessage.MsgVerAck{}, nil
}

func (x *KaspadP2PMessage_Verack) fromAppMessage(_ *appmessage.MsgVerAck) error {
	return nil
}
