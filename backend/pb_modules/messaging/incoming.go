package messaging

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"strings"

	"github.com/labstack/echo/v5"
)

type incomingMessage struct {
	From string
	Body string
}

// TwiMLResponse is the XML envelope Twilio expects back.
type TwiMLResponse struct {
	XMLName xml.Name `xml:"Response"`
	Message string   `xml:"Message"`
}

func parseIncoming(c echo.Context) (incomingMessage, error) {
	ct := c.Request().Header.Get("Content-Type")

	if strings.Contains(ct, "application/x-www-form-urlencoded") {
		// Twilio format
		return incomingMessage{
			From: c.FormValue("From"),
			Body: strings.TrimSpace(c.FormValue("Body")),
		}, nil
	}

	// JSON format (simulator / curl)
	raw, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return incomingMessage{}, err
	}
	var msg incomingMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return incomingMessage{}, err
	}
	msg.Body = strings.TrimSpace(msg.Body)
	return msg, nil
}
