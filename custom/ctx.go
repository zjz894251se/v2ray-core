package custom

import (
	"context"
	"github.com/v2fly/v2ray-core/v5/common/session"
	"sync"
	"time"
)

type clientInfo struct {
	ClientId     string    // id
	StartConnect time.Time // 开始链接时间
}

type SessionClientKV struct {
	sync.Mutex
	mps          map[session.ID]*clientInfo
	onBindCall   func(string)
	onUnbindCall func(string, time.Duration)
}

func (kv *SessionClientKV) Bind(sessionId session.ID, clientId string) {
	kv.Lock()
	defer kv.Unlock()
	kv.mps[sessionId] = &clientInfo{
		ClientId:     clientId,
		StartConnect: time.Now(),
	}
}

func (kv *SessionClientKV) GetClientId(sessionId session.ID) string {
	kv.Lock()
	defer kv.Unlock()
	info, ok := kv.mps[sessionId]
	if ok {
		return info.ClientId
	}
	return ""
}

func (kv *SessionClientKV) UnBind(sessionId session.ID) {
	kv.Lock()
	defer kv.Unlock()
	info, ok := kv.mps[sessionId]
	if ok {
		delete(kv.mps, sessionId)
		if kv.onUnbindCall != nil {
			kv.onUnbindCall(info.ClientId, time.Duration(time.Now().UnixNano()-info.StartConnect.UnixNano()))
		}
	}
}

func (kv *SessionClientKV) RegisterOnBindCall(call func(string)) {
	kv.onBindCall = call
}

func (kv *SessionClientKV) RegisterOnUnbindCall(call func(string, time.Duration)) {
	kv.onUnbindCall = call
}

var _sessionClientKv = &SessionClientKV{mps: make(map[session.ID]*clientInfo)}

func BindSessionClientId(ctx context.Context, clientId string) {
	_sessionClientKv.Bind(session.IDFromContext(ctx), clientId)
}

func ClientIdFromContext(ctx context.Context) string {
	return _sessionClientKv.GetClientId(session.IDFromContext(ctx))
}

func UnBindSessionClientId(ctx context.Context) {
	_sessionClientKv.UnBind(session.IDFromContext(ctx))
}

// RegisterSessionOnBindCall 客户链接回调
func RegisterSessionOnBindCall(call func(string)) {
	_sessionClientKv.RegisterOnBindCall(call)
}

// RegisterSessionOnUnbindCall 客户断开回调
func RegisterSessionOnUnbindCall(call func(string, time.Duration)) {
	_sessionClientKv.RegisterOnUnbindCall(call)
}
