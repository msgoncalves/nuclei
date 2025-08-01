package rdp

import (
	"context"
	"fmt"
	"time"

	"github.com/praetorian-inc/fingerprintx/pkg/plugins"
	"github.com/praetorian-inc/fingerprintx/pkg/plugins/services/rdp"
	"github.com/projectdiscovery/nuclei/v3/pkg/protocols/common/protocolstate"
)

type (
	// IsRDPResponse is the response from the IsRDP function.
	// this is returned by IsRDP function.
	// @example
	// ```javascript
	// const rdp = require('nuclei/rdp');
	// const isRDP = rdp.IsRDP('acme.com', 3389);
	// log(toJSON(isRDP));
	// ```
	IsRDPResponse struct {
		IsRDP bool
		OS    string
	}
)

// IsRDP checks if the given host and port are running rdp server.
// If connection is successful, it returns true.
// If connection is unsuccessful, it returns false and error.
// The Name of the OS is also returned if the connection is successful.
// @example
// ```javascript
// const rdp = require('nuclei/rdp');
// const isRDP = rdp.IsRDP('acme.com', 3389);
// log(toJSON(isRDP));
// ```
func IsRDP(ctx context.Context, host string, port int) (IsRDPResponse, error) {
	executionId := ctx.Value("executionId").(string)
	return memoizedisRDP(executionId, host, port)
}

// @memo
func isRDP(executionId string, host string, port int) (IsRDPResponse, error) {
	resp := IsRDPResponse{}

	dialer := protocolstate.GetDialersWithId(executionId)
	if dialer == nil {
		return IsRDPResponse{}, fmt.Errorf("dialers not initialized for %s", executionId)
	}

	timeout := 5 * time.Second
	conn, err := dialer.Fastdialer.Dial(context.TODO(), "tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return resp, err
	}
	defer func() {
		_ = conn.Close()
	}()

	server, isRDP, err := rdp.DetectRDP(conn, timeout)
	if err != nil {
		return resp, err
	}
	if !isRDP {
		return resp, nil
	}
	resp.IsRDP = true
	resp.OS = server
	return resp, nil
}

type (
	// CheckRDPAuthResponse is the response from the CheckRDPAuth function.
	// this is returned by CheckRDPAuth function.
	// @example
	// ```javascript
	// const rdp = require('nuclei/rdp');
	// const checkRDPAuth = rdp.CheckRDPAuth('acme.com', 3389);
	// log(toJSON(checkRDPAuth));
	// ```
	CheckRDPAuthResponse struct {
		PluginInfo *plugins.ServiceRDP
		Auth       bool
	}
)

// CheckRDPAuth checks if the given host and port are running rdp server
// with authentication and returns their metadata.
// If connection is successful, it returns true.
// @example
// ```javascript
// const rdp = require('nuclei/rdp');
// const checkRDPAuth = rdp.CheckRDPAuth('acme.com', 3389);
// log(toJSON(checkRDPAuth));
// ```
func CheckRDPAuth(ctx context.Context, host string, port int) (CheckRDPAuthResponse, error) {
	executionId := ctx.Value("executionId").(string)
	return memoizedcheckRDPAuth(executionId, host, port)
}

// @memo
func checkRDPAuth(executionId string, host string, port int) (CheckRDPAuthResponse, error) {
	resp := CheckRDPAuthResponse{}

	dialer := protocolstate.GetDialersWithId(executionId)
	if dialer == nil {
		return CheckRDPAuthResponse{}, fmt.Errorf("dialers not initialized for %s", executionId)
	}
	timeout := 5 * time.Second
	conn, err := dialer.Fastdialer.Dial(context.TODO(), "tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return resp, err
	}
	defer func() {
		_ = conn.Close()
	}()

	pluginInfo, auth, err := rdp.DetectRDPAuth(conn, timeout)
	if err != nil {
		return resp, err
	}
	if !auth {
		return resp, nil
	}
	resp.Auth = true
	resp.PluginInfo = pluginInfo
	return resp, nil
}
