# EventFlow ![Go Report Card](https://goreportcard.com/badge/github.com/tung-lin/EventFlow)

## Dependencies

- Go 1.9.6

## Introduction

事件通報服務使用模組化開發方式，分為觸發(trigger)、資料處理(filter)、行為(action)三個執行階段，並使用yaml定義執行的工作。

- **trigger:** 事件觸發來源，使用udp、http等方式接收事件觸發。

- **filter:** 事件資料處理，包含json格式解析、觸發頻率限制等。

- **action:** 事件通報方式，包含email通報、line通報等。

## Configuration

- **config/config.yaml:** 服務設定檔。

- **config/pipeline/:** 工作pipeline設定檔，可同時有多個檔案。





