package protowire

import (
	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/pkg/errors"
)

func (x *KaspadMessage_GetBlockTemplatesRequest) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "KaspadMessage_GetBlockTemplatesRequest is nil")
	}
	return x.GetBlockTemplatesRequest.toAppMessage()
}

func (x *KaspadMessage_GetBlockTemplatesRequest) fromAppMessage(message *appmessage.GetBlockTemplatesRequestMessage) error {
	templateRequests := make([]*TemplateRequest, len(message.TemplatesRequest))

	for i, templateRequest := range message.TemplatesRequest {
		templateRequests[i].fromAppMessage(templateRequest)
	}

	x.GetBlockTemplatesRequest = &GetBlockTemplatesRequestMessage{
		TemplateRequests: templateRequests,
	}
	return nil
}

func (x *GetBlockTemplatesRequestMessage) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "GetBlockTemplatesRequestMessage is nil")
	}
	templateRequests := make([]*appmessage.TemplatesRequest, len(x.TemplateRequests))

	for i, templateRequest := range x.TemplateRequests {
		templateRequests[i] = templateRequest.toAppMessage()
	}

	return &appmessage.GetBlockTemplatesRequestMessage{
		TemplatesRequest: templateRequests,
	}, nil
}

func (x *KaspadMessage_GetBlockTemplatesResponse) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "KaspadMessage_GetBlockTemplateResponse is nil")
	}
	return x.GetBlockTemplatesResponse.toAppMessage()
}

func (x *KaspadMessage_GetBlockTemplatesResponse) fromAppMessage(message *appmessage.GetBlockTemplatesResponseMessage) error {
	var err *RPCError
	if message.Error != nil {
		err = &RPCError{Message: message.Error.Message}
	}

	templateResponses := make([]*TemplateReponse, len(message.TemplatesResponse))

	for i, templateResponse := range message.TemplatesResponse {
		templateResponses[i].fromAppMessage(templateResponse)
	}

	x.GetBlockTemplatesResponse = &GetBlockTemplatesResponseMessage{
		TemplateResponses: templateResponses,
		IsSynced:          message.IsSynced,
		Error:             err,
	}
	return nil
}

func (x *GetBlockTemplatesResponseMessage) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "GetBlockTemplatesResponseMessage is nil")
	}
	var rpcError *appmessage.RPCError
	if x.Error != nil {
		var err error
		rpcError, err = x.Error.toAppMessage()
		if err != nil {
			return nil, err
		}
	}

	templatesResponses := make([]*appmessage.TemplatesResponse, len(x.TemplateResponses))
	var err error
	for i, templateResponse := range x.TemplateResponses {
		templatesResponses[i], err = templateResponse.toAppMessage()
		if err != nil {
			return nil, err
		}
	}

	return &appmessage.GetBlockTemplatesResponseMessage{
		TemplatesResponse: templatesResponses,
		IsSynced:          x.IsSynced,
		Error:             rpcError,
	}, nil
}

func (x *TemplateReponse) toAppMessage() (*appmessage.TemplatesResponse, error) {
	blockTemplate, err := x.BlockTemplate.toAppMessage()
	if err != nil {
		return nil, err
	}

	return &appmessage.TemplatesResponse{
		Address: x.Address,
		Block:   blockTemplate,
	}, nil
}

func (x *TemplateRequest) toAppMessage() *appmessage.TemplatesRequest {
	return &appmessage.TemplatesRequest{
		PayAddress: x.PayAddress,
		ExtraData:  x.ExtraData,
	}

}

func (x *TemplateReponse) fromAppMessage(message *appmessage.TemplatesResponse) error {
	var block *RpcBlock
	if message.Block != nil {
		protoBlock := &RpcBlock{}
		err := protoBlock.fromAppMessage(message.Block)
		if err != nil {
			return err
		}
		block = protoBlock
	}

	x.BlockTemplate = block
	x.Address = message.Address
	return nil
}

func (x *TemplateRequest) fromAppMessage(message *appmessage.TemplatesRequest) {
	x.PayAddress = message.PayAddress
	x.ExtraData = message.ExtraData
}
