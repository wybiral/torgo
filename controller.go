package torgo

import (
	"net/textproto"
	"fmt"
	"io/ioutil"
	"strings"
)

type Controller struct {
	conn *textproto.Conn
}

func NewController(addr string) (*Controller, error) {
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Controller{conn: conn}, nil
}

func (c *Controller) makeRequest(request string) (int, string, error) {
	id, err := c.conn.Cmd(request)
	if err != nil {
		return 0, "", err
	}
	c.conn.StartResponse(id)
	defer c.conn.EndResponse(id)
	return c.conn.ReadResponse(-1)
}

func (c *Controller) AuthenticateCookie() error {
	cookie, err := ioutil.ReadFile("/var/run/tor/control.authcookie")
	if err != nil {
		return err
	}
	code, msg, err := c.makeRequest("AUTHENTICATE \"" + string(cookie) + `"`)
	if err != nil {
		return err
	}
	if code != 250 {
		return fmt.Errorf("%d %s", code, msg)
	}
	return nil
}

func (c *Controller) Add(onion *Onion) error {
	req := "ADD_ONION "
	if len(onion.PrivateKey) == 0 {
		onion.PrivateKeyType = "NEW"
		onion.PrivateKey = "BEST"
	}
	req += fmt.Sprintf("%s:%s ", onion.PrivateKeyType, onion.PrivateKey)
	for remotePort, localPort := range onion.Ports {
		req += fmt.Sprintf("Port=%d,%d ", remotePort, localPort)
	}
	code, msg, err := c.makeRequest(req)
	if err != nil {
		return err
	}
	if code != 250 {
		return fmt.Errorf("%d %s", code, msg)
	}
	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if parts[0] == "ServiceID" {
			onion.ServiceID = parts[1]
		} else if parts[0] == "PrivateKey" {
			key := strings.SplitN(parts[1], ":", 2)
			onion.PrivateKeyType = key[0]
			onion.PrivateKey = key[1]
		}
	}
	return nil
}

func (c *Controller) Remove(onion *Onion) error {
	code, msg, err := c.makeRequest("DEL_ONION " + onion.ServiceID)
	if err != nil {
		return err
	}
	if code != 250 {
		return fmt.Errorf("%d %s", code, msg)
	}
	return nil
}
