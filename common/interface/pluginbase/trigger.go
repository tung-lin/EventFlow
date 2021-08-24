package pluginbase

import (
	"EventFlow/common/tool/pipelinetool"
)

//ITriggerFactory interfacr for factory
type ITriggerFactory interface {
	GetIdentifyName() string
	CreateTrigger(config interface{}) ITriggerPlugin
}

//ITriggerPlugin interface for trigger plugin
type ITriggerPlugin interface {
	IActionHandler
	Start()
	Stop()
}

//IActionHandler interface for action handler
type IActionHandler interface {
	ActionHandleFunc(triggerPlugin *ITriggerPlugin, pipelineID string, f func(pipelineID string, triggerPlugin *ITriggerPlugin, messageFromTrigger *string))
	FireAction(triggerPlugin *ITriggerPlugin, messageFromTrigger *string)
	CloseChannel()
}

const (
	pipelineChannelBufferSize = 300
)

//ActionHandler struct for action handler
type ActionHandler struct {
	pipelineID       string
	actionHandleFunc func(pipelineID string, triggerPlugin *ITriggerPlugin, messageFromTrigger *string)
	channel          chan *string
}

//ActionHandleFunc add handler function
func (h *ActionHandler) ActionHandleFunc(triggerPlugin *ITriggerPlugin, pipelineID string, f func(pipelineID string, triggerPlugin *ITriggerPlugin, messageFromTrigger *string)) {
	h.pipelineID = pipelineID
	h.actionHandleFunc = f
	h.channel = make(chan *string, pipelineChannelBufferSize)

	go func() {
		for {
			message, moreData := <-h.channel

			if !moreData {
				break
			}

			pipelinetool.GlobalPipelineWaitGroup.Add(1)

			//execute filters and actions
			go func() {

				defer func() {
					pipelinetool.GlobalPipelineWaitGroup.Done()
				}()

				if h.actionHandleFunc != nil {
					h.actionHandleFunc(pipelineID, triggerPlugin, message)
				}
			}()
		}
	}()
}

//FireAction fire action
func (h *ActionHandler) FireAction(triggerPlugin *ITriggerPlugin, messageFromTrigger *string) {
	//if h.actionHandleFunc != nil {
	//	h.actionHandleFunc(h.pipelineID, triggerPlugin, messageFromTrigger)
	//}
	h.channel <- messageFromTrigger
}

//CloseChannel gracefully close channel
func (h *ActionHandler) CloseChannel() {
	close(h.channel)
}
