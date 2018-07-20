/*
 * Copyright (c) 2018 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package xml_test

import (
	"testing"

	"github.com/ortuman/jackal/xml"
	"github.com/ortuman/jackal/xml/jid"
	"github.com/stretchr/testify/require"
)

func TestMessageBuild(t *testing.T) {
	j, _ := jid.New("ortuman", "example.org", "balcony", false)

	elem := xml.NewElementName("iq")
	_, err := xml.NewMessageFromElement(elem, j, j) // wrong name...
	require.NotNil(t, err)

	elem.SetName("message")
	elem.SetType("invalid")
	_, err = xml.NewMessageFromElement(elem, j, j) // invalid type...
	require.NotNil(t, err)

	// valid message...
	elem.SetType(xml.ChatType)
	elem.AppendElement(xml.NewElementName("body"))
	message, err := xml.NewMessageFromElement(elem, j, j)
	require.Nil(t, err)
	require.NotNil(t, message)
	require.True(t, message.IsMessageWithBody())

	msg2 := xml.NewMessageType("an-id123", xml.GroupChatType)
	require.Equal(t, "an-id123", msg2.ID())
	require.Equal(t, xml.GroupChatType, msg2.Type())
}

func TestMessageType(t *testing.T) {
	message, _ := xml.NewMessageFromElement(xml.NewElementName("message"), &jid.JID{}, &jid.JID{})
	require.True(t, message.IsNormal())

	message.SetType(xml.NormalType)
	require.True(t, message.IsNormal())

	message.SetType(xml.HeadlineType)
	require.True(t, message.IsHeadline())

	message.SetType(xml.ChatType)
	require.True(t, message.IsChat())

	message.SetType(xml.GroupChatType)
	require.True(t, message.IsGroupChat())
}

func TestMessageJID(t *testing.T) {
	from, _ := jid.New("ortuman", "test.org", "balcony", false)
	to, _ := jid.New("ortuman", "example.org", "garden", false)
	message, _ := xml.NewMessageFromElement(xml.NewElementName("message"), &jid.JID{}, &jid.JID{})
	message.SetFromJID(from)
	require.Equal(t, message.FromJID().String(), message.From())
	message.SetToJID(to)
	require.Equal(t, message.ToJID().String(), message.To())
}
