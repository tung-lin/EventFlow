package udp

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/logtool"
	"fmt"
	"net"
)

type UDPPlugin struct {
	currentConnection *net.UDPConn
	pluginbase.ActionHandler
	Setting SettingConfig
}

func (trigger *UDPPlugin) Start() {

	listenAddress := fmt.Sprintf(":%d", trigger.Setting.Port)
	udpAddress, err := net.ResolveUDPAddr("udp", listenAddress)

	if err != nil {
		logtool.Error("trigger", "udp", fmt.Sprintf("resolve udp endpoint failed : %v", err))
		return
	}

	connection, err := net.ListenUDP("udp", udpAddress)

	if err != nil {
		logtool.Error("trigger", "udp", fmt.Sprintf("listen udp endpoint(%s) failed : %v", udpAddress.String(), err))
		return
	}

	trigger.currentConnection = connection

	go readFromUDP(trigger)
}

func readFromUDP(trigger *UDPPlugin) {
	receiveBuffer := make([]byte, 1024)

	for {
		logtool.Debug("trigger", "udp", fmt.Sprintf("udp(%s) ready to receive message", trigger.currentConnection.LocalAddr().String()))

		received, remoteAddress, err := trigger.currentConnection.ReadFromUDP(receiveBuffer)

		if err != nil {
			logtool.Error("trigger", "udp", fmt.Sprintf("read message from udp(remote address: %s) failed : %v", remoteAddress.String(), err))
			break
		}

		message := string(receiveBuffer[:received])

		logtool.Debug("trigger", "udp", fmt.Sprintf("received message from %s, message: %s", trigger.currentConnection.LocalAddr().String(), message))

		var triggerPlugin pluginbase.ITriggerPlugin = trigger
		trigger.FireAction(&triggerPlugin, &message)
	}
}

func (trigger *UDPPlugin) Stop() {
	if trigger.currentConnection != nil {
		trigger.currentConnection.Close()
	}
}
