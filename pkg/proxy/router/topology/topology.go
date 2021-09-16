// Copyright 2014 Wandoujia Inc. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package topology

import (
	"github.com/IceFireDB/kit/pkg/models"
	"github.com/IceFireDB/kit/pkg/models/client"
	"github.com/juju/errors"
	log "github.com/ngaut/logging"
)

type TopoUpdate interface {
	OnGroupChange(groupId int)
	OnSlotChange(slotId int)
	// todo if need on server changed
}

type Topology struct {
	ProductName string
	store       *models.Store
}

func (top *Topology) GetGroup(groupId int) (*models.ServerGroup, error) {
	return top.store.LoadGroup(groupId, true)
}

func (top *Topology) Exist(path string) (bool, error) {
	return top.store.Exists(path)
}

func (top *Topology) GetSlotByIndex(i int) (*models.Slot, *models.ServerGroup, error) {
	slot, err := top.store.GetSlot(i, true)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	log.Debugf("get slot %d : %+v", i, slot)
	if slot.State.Status != models.SLOT_STATUS_ONLINE && slot.State.Status != models.SLOT_STATUS_MIGRATE {
		log.Errorf("slot not online, %+v", slot)
	}

	groupServer, err := top.store.LoadGroup(slot.GroupId, true)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	return slot, groupServer, nil
}

func NewTopo(productName string, store *models.Store) *Topology {
	t := &Topology{ProductName: productName, store: store}
	return t
}

func (top *Topology) GetActionWithSeq(seq string) (*models.Action, error) {
	return top.store.GetActionWithSeq(seq)
}

func (top *Topology) GetActionWithSeqObject(seq string, act *models.Action) error {
	return top.store.GetActionObject(seq, act)
}

func (top *Topology) GetActionSeqList() ([]string, error) {
	return top.store.GetActionSeqList()
}

func (top *Topology) IsChildrenChangedEvent(e client.Event) bool {
	return e.Type == client.EventNodeChildrenChanged
}

func (top *Topology) CreateProxyInfo(pi *models.ProxyInfo) error {
	return top.store.UpdateProxy(pi)
}

func (top *Topology) GetProxyInfo(proxyName string) (*models.ProxyInfo, error) {
	return top.store.LoadProxy(proxyName)
}

func (top *Topology) GetActionResponsePath(seq string) string {
	return top.store.ActionPath(seq)
}

func (top *Topology) Close(proxyName string) {
	//zkhelper.DeleteRecursive(top.zkConn, path.Join(models.GetProxyPath(top.ProductName), proxyName), -1)
	//top.zkConn.Close()
}


func (top *Topology) doWatch(evtch <-chan client.Event, evtbus chan client.Event) {
	e := <-evtch
	log.Infof("topo event %+v", e)
	//if e.State == topo.StateExpired {
	//	log.Fatalf("session expired: %+v", e)
	//}

	switch e.Type {
	// case topo.EventNodeCreated:
	// case topo.EventNodeDataChanged:
	case client.EventNodeChildrenChanged: // only care children changed
		//todo:get changed node and decode event
	default:
		log.Warningf("%+v", e)
	}

	evtbus <- e
}

func (top *Topology) WatchAction(evtbus chan client.Event) ([]string, error) {
	evtch, content, err := top.store.WatchActions()
	if err != nil {
		return nil, errors.Trace(err)
	}

	go top.doWatch(evtch, evtbus)
	return content, nil
}
