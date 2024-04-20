package tool

import (
	nm "go-node/models"
	"snproto"
)

func ConvertRpcInfoToNodeInfo(sp *snproto.NodeAppInfo) *nm.NodeApp {
	at := nm.AppType(sp.NodeAppType)

	ca := sp.CreatedAt.AsTime()
	da := sp.DeletedAt.AsTime()
	ua := sp.UpdatedAt.AsTime()
	return &nm.NodeApp{
		NodeAppId:      sp.NodeAppId,
		NodeId:         sp.NodeId,
		NodeAppType:    &at,
		NodeAppName:    sp.NodeAppName,
		NodeAppIntro:   sp.NodeAppIntro,
		NodeAppPort:    sp.NodeAppPort,
		NodeAppDomain:  sp.NodeAppDomain,
		NodeAppStatus:  int(sp.NodeAppStatus),
		NodeAppVersion: int(sp.NodeAppVersion),
		CreatedAt:      &ca,
		DeletedAt:      &da,
		UpdatedAt:      &ua,
	}
}

func ConvertMsgToNodeInfo(s *snproto.NodeAppInfo) *nm.NodeApp {
	at := nm.AppType(s.NodeAppType)
	ca := s.CreatedAt.AsTime()
	da := s.DeletedAt.AsTime()
	ua := s.UpdatedAt.AsTime()
	n := &nm.NodeApp{
		NodeAppId:      s.NodeAppId,
		NodeId:         s.NodeId,
		NodeAppName:    s.NodeAppName,
		NodeAppType:    &at,
		NodeAppIntro:   s.NodeAppIntro,
		NodeAppPort:    s.NodeAppPort,
		NodeAppDomain:  s.NodeAppDomain,
		NodeAppStatus:  int(s.NodeAppStatus),
		NodeAppVersion: int(s.NodeAppVersion),
		CreatedAt:      &ca,
		DeletedAt:      &da,
		UpdatedAt:      &ua,
	}
	return n
}

func ConvertMsgToNodeSys(s *snproto.NodeSys) *nm.NodeSys {
	ca := s.CreatedAt.AsTime()
	da := s.DeletedAt.AsTime()
	ua := s.UpdatedAt.AsTime()
	return &nm.NodeSys{
		NodeId:     s.NodeId,
		NodeName:   s.NodeName,
		NodeStatus: int(s.NodeStatus),
		NodeDomain: s.NodeDomain,
		CreatedAt:  &ca,
		DeletedAt:  &da,
		UpdatedAt:  &ua,
	}
}
