/*
 * Copyright (c) 2018 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package offline

import (
	"github.com/ortuman/jackal/log"
	"github.com/ortuman/jackal/module/xep0030"
	"github.com/ortuman/jackal/storage"
	"github.com/ortuman/jackal/stream"
	"github.com/ortuman/jackal/xml"
)

const offlineNamespace = "msgoffline"

// Config represents Offline Storage module configuration.
type Config struct {
	QueueSize int `yaml:"queue_size"`
}

// Offline represents an offline server stream module.
type Offline struct {
	cfg     *Config
	stm     stream.C2S
	actorCh chan func()
}

// New returns an offline server stream module.
func New(config *Config, stm stream.C2S) *Offline {
	r := &Offline{
		cfg:     config,
		stm:     stm,
		actorCh: make(chan func(), 32),
	}
	go r.actorLoop(stm.Context().Done())
	return r
}

// RegisterDisco registers disco entity features/items
// associated to offline module.
func (o *Offline) RegisterDisco(discoInfo *xep0030.DiscoInfo) {
	discoInfo.Entity(o.stm.Domain(), "").AddFeature(offlineNamespace)
}

// ArchiveMessage archives a new offline messages into the storage.
func (o *Offline) ArchiveMessage(message *xml.Message) {
	o.actorCh <- func() {
		o.archiveMessage(message)
	}
}

// DeliverOfflineMessages delivers every archived offline messages to the peer
// deleting them from storage.
func (o *Offline) DeliverOfflineMessages() {
	o.actorCh <- func() {
		o.deliverOfflineMessages()
	}
}

func (o *Offline) actorLoop(doneCh <-chan struct{}) {
	for {
		select {
		case f := <-o.actorCh:
			f()
		case <-doneCh:
			return
		}
	}
}

func (o *Offline) archiveMessage(message *xml.Message) {
	toJid := message.ToJID()
	queueSize, err := storage.Instance().CountOfflineMessages(toJid.Node())
	if err != nil {
		log.Error(err)
		return
	}
	if queueSize >= o.cfg.QueueSize {
		response := xml.NewElementFromElement(message)
		response.SetFrom(toJid.String())
		response.SetTo(o.stm.JID().String())
		o.stm.SendElement(response.ServiceUnavailableError())
		return
	}
	delayed := xml.NewElementFromElement(message)
	delayed.Delay(o.stm.Domain(), "Offline Storage")
	if _, _, err := storage.Instance().WriteMsgToDB(delayed.Text(), toJid.Node(), message.Text(), 0, 0); err != nil {
		log.Errorf("%v", err)
		return
	}
	log.Infof("archived offline message... id: %s", message.ID())
}

func (o *Offline) deliverOfflineMessages() {
	messages, err := storage.Instance().FetchOfflineMessages(o.stm.Username())
	if err != nil {
		log.Error(err)
		return
	}
	if len(messages) == 0 {
		return
	}
	log.Infof("delivering offline messages... count: %d", len(messages))

	for _, m := range messages {
		o.stm.SendElement(m)
	}
	if err := storage.Instance().DeleteOfflineMessages(o.stm.Username()); err != nil {
		log.Error(err)
	}
}
