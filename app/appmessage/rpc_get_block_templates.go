package appmessage

// GetBlockTemplatesRequestMessage is an appmessage corresponding to
// its respective RPC message
type GetBlockTemplatesRequestMessage struct {
	baseMessage
	TemplatesRequest []*TemplatesRequest
}

// TemplatesRequest is an appmessage corresponding to
// its respective RPC message
type TemplatesRequest struct {
	PayAddress string
	ExtraData  string
}

// GetBlockTemplatesResponseMessage is an appmessage corresponding to
// its respective RPC message
type GetBlockTemplatesResponseMessage struct {
	baseMessage
	TemplatesResponse []*TemplatesResponse
	IsSynced          bool
	Error             *RPCError
}

// TemplatesResponse is an appmessage corresponding to
// its respective RPC message
type TemplatesResponse struct {
	Address string
	Block   *RPCBlock
}

// Command returns the protocol command string for the message
func (msg *GetBlockTemplatesRequestMessage) Command() MessageCommand {
	return CmdGetBlockTemplatesRequestMessage
}

// NewGetBlockTemplatesRequestMessage returns a instance of the message
func NewGetBlockTemplatesRequestMessage(templatesRequest []*TemplatesRequest) *GetBlockTemplatesRequestMessage {
	return &GetBlockTemplatesRequestMessage{
		TemplatesRequest: templatesRequest,
	}
}

// Command returns the protocol command string for the message
func (msg *GetBlockTemplatesResponseMessage) Command() MessageCommand {
	return CmdGetBlockTemplatesResponseMessage
}

// NewGetBlockTemplatesResponseMessage returns a instance of the message
func NewGetBlockTemplatesResponseMessage(templatesResponses []*TemplatesResponse, isSynced bool) *GetBlockTemplatesResponseMessage {
	return &GetBlockTemplatesResponseMessage{
		TemplatesResponse: templatesResponses,
		IsSynced:          isSynced,
	}
}
