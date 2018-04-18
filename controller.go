package torgo

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/textproto"
	"strconv"
	"strings"
)

type Controller struct {
	conn        *textproto.Conn
	AuthMethods []string
	CookieFile  string
}

func NewController(addr string) (*Controller, error) {
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	c := &Controller{conn: conn}
	err = c.getProtocolInfo()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Controller) makeRequest(request string) (int, string, error) {
	id, err := c.conn.Cmd(request)
	if err != nil {
		return 0, "", err
	}
	c.conn.StartResponse(id)
	defer c.conn.EndResponse(id)
	return c.conn.ReadResponse(250)
}

func (c *Controller) getProtocolInfo() error {
	_, msg, err := c.makeRequest("PROTOCOLINFO 1")
	if err != nil {
		return err
	}
	lines := strings.Split(msg, "\n")
	authPrefix := "AUTH METHODS="
	cookiePrefix := "COOKIEFILE="
	for _, line := range lines {
		// Check for AUTH METHODS line
		if strings.HasPrefix(line, authPrefix) {
			// This line should be in a format like:
			/// AUTH METHODS=method1,method2,methodN COOKIEFILE=cookiefilepath
			line = line[len(authPrefix):]
			parts := strings.SplitN(line, " ", 2)
			c.AuthMethods = strings.Split(parts[0], ",")
			// Check gor COOKIEFILE key/value
			if len(parts) == 2 && strings.HasPrefix(parts[1], cookiePrefix) {
				raw := parts[1][len(cookiePrefix):]
				c.CookieFile, err = strconv.Unquote(raw)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c *Controller) GetVersion() (string, error) {
	_, msg, err := c.makeRequest("GETINFO version")
	if err != nil {
		return "", err
	}
	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if parts[0] == "version" {
			return parts[1], nil
		}
	}
	return "", fmt.Errorf("version not found")
}

func (c *Controller) AuthenticateNone() error {
	_, _, err := c.makeRequest("AUTHENTICATE")
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) AuthenticatePassword(password string) error {
	quoted := strconv.Quote(password)
	_, _, err := c.makeRequest("AUTHENTICATE " + quoted)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) AuthenticateCookie() error {
	rawCookie, err := ioutil.ReadFile(c.CookieFile)
	if err != nil {
		return err
	}
	cookie := hex.EncodeToString(rawCookie)
	_, _, err = c.makeRequest("AUTHENTICATE " + cookie)
	if err != nil {
		return err
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
	for remotePort, localAddr := range onion.Ports {
		req += fmt.Sprintf("Port=%d,%s ", remotePort, localAddr)
	}
	_, msg, err := c.makeRequest(req)
	if err != nil {
		return err
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
	_, _, err := c.makeRequest("DEL_ONION " + onion.ServiceID)
	if err != nil {
		return err
	}
	return nil
}
