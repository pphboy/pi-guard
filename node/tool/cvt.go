package tool

import (
	"go-node/models"
	"snproto"
)

func ConvertRpcInfoToNodeInfo(sp *snproto.NodeAppInfo) *models.NodeApp {
	at := models.AppType(sp.NodeAppType)

	ca := sp.CreatedAt.AsTime()
	da := sp.DeletedAt.AsTime()
	ua := sp.UpdatedAt.AsTime()
	return &models.NodeApp{
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
