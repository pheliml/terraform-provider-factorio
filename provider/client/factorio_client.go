package client

import (
	"encoding/json"
	"fmt"
)

type FactorioClient struct {
	conn *RCON
}

func NewFactorioClient(rconHost string, rconPassword string) (*FactorioClient, error) {
	r, err := Dial(rconHost)
	if err != nil {
		return nil, err
	}
	err = r.Authenticate(rconPassword)
	if err != nil {
		return nil, err
	}
	var c FactorioClient
	c.conn = r
	err = c.DoHandShake()
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (client *FactorioClient) DoHandShake() error {
	var result string
	// Ignore the first error.
	// Execute the ping twice in order to skip past the warning about
	// how Lua console commands will disable achievements
	client.doCall(&result, "ping")
	err := client.doCall(&result, "ping")
	if err != nil {
		return err
	}
	if result != "pong" {
		return fmt.Errorf("expected \"pong\" from handshake but got \"%s\"", result)
	}
	return nil
}

func (client *FactorioClient) Read(resourceType string, query interface{}, resultOut interface{}) error {
	return client.doCall(resultOut, "read", resourceType, query)
}

func (client *FactorioClient) Create(resourceType string, createConfig interface{}, resultOut interface{}) error {
	return client.doCall(resultOut, "create", resourceType, createConfig)
}

// Perhaps Update should just return success / failure?
func (client *FactorioClient) Update(resourceType string, resourceID string, updateOpts interface{}, resultOut interface{}) error {
	return client.doCall(resultOut, "update", resourceType, resourceID, updateOpts)
}

func (client *FactorioClient) Delete(resourceType string, resourceID string) error {
	result := struct {
		ResourceExists bool `json:"resource_exists"`
	}{
		ResourceExists: true,
	}
	err := client.doCall(&result, "delete", resourceType, resourceID)
	if err != nil {
		return err
	}
	if result.ResourceExists {
		return fmt.Errorf("resource still exists")
	}
	return nil
}

type RPCRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RPCResponse struct {
	Result *json.RawMessage `json:"result"`
	Error  *RPCError        `json:"error"`
}

func (client *FactorioClient) doCall(result interface{}, method string, params ...interface{}) error {
	if params == nil {
		params = []interface{}{}
	}
	req := RPCRequest{
		Method: method,
		Params: params,
	}
	requestBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	// Use single quotes around request_bytes
	// to avoid conflict with json double quotes
	// TODO: Escape single quotes in request_bytes
	command := fmt.Sprintf("/silent-command rcon.print(remote.call('terraform-crud-api', 'call', '%s'))", requestBytes)
	executeResponse, err := client.conn.Execute(command)
	if err != nil {
		return err
	}
	var response RPCResponse
	err = json.Unmarshal([]byte(executeResponse), &response)
	if err != nil {
		return fmt.Errorf("unmarshalling \"%v\": %v", executeResponse, err)
	}
	if response.Error != nil {
		return fmt.Errorf(
			"error from api, code: %d, message: \"%s\", details: %+v",
			response.Error.Code,
			response.Error.Message,
			response.Error.Data)
	}
	// Lua nil does not get serialized to null, so a missing Result
	// is interpreted as a null Result
	if response.Result == nil {
		null := json.RawMessage("null")
		response.Result = &null
	}
	return json.Unmarshal(*response.Result, result)
}
