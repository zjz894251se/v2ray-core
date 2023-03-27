package custom

import (
	"context"
	"github.com/v2fly/v2ray-core/v5/common/session"
	"sync"
)

type SessionClientKV struct {
	sync.Mutex
	mps        map[session.ID]string
	onBindCall func(session.ID, string)
}

func (kv *SessionClientKV) Bind(sessionId session.ID, clientId string) {
	kv.Lock()
	defer kv.Unlock()
	kv.mps[sessionId] = clientId
}

func (kv *SessionClientKV) GetClientId(sessionId session.ID) string {
	kv.Lock()
	defer kv.Unlock()
	return kv.mps[sessionId]
}

func (kv *SessionClientKV) UnBind(sessionId session.ID) {
	kv.Lock()
	defer kv.Unlock()
	delete(kv.mps, sessionId)
}

func (kv *SessionClientKV) RegisterOnBindCall(call func(session.ID, string)) {
	kv.onBindCall = call
}

var _sessionClientKv = &SessionClientKV{mps: make(map[session.ID]string)}

func BindSessionClientId(ctx context.Context, clientId string) {
	_sessionClientKv.Bind(session.IDFromContext(ctx), clientId)
}

func ClientIdFromContext(ctx context.Context) string {
	return _sessionClientKv.GetClientId(session.IDFromContext(ctx))
}

func UnBindSessionClientId(ctx context.Context) {
	_sessionClientKv.UnBind(session.IDFromContext(ctx))
}

func RegisterSessionOnBindCall(call func(session.ID, string)) {
	_sessionClientKv.RegisterOnBindCall(call)
}
