package dashboard

import dao "git/inspursoft/board/src/common/dao/dashboard"

type NodeDs struct {
	nodeName   string
	recordTime string
	count      string
	timeUnit   string
	OutJson    []byte
}
type Model struct {
	NodeName     string `json:"node_name"`
	NodeTimeunit string `json:"node_timeunit"`
	NodeCount    string `json:"node_count"`
}

func (ds NodeDs) SetNode(nodeName string, recordTime string, count string, timeUnit string) {
	ds.nodeName = nodeName
	ds.recordTime = recordTime
	ds.count = count
	ds.timeUnit = timeUnit

}
func (ds NodeDs) genNodeList() []dao.Node {
	return dao.QueryNode(ds.nodeName, ds.recordTime, ds.count, ds.timeUnit)
}
func (ds NodeDs) preJson() {
	//for k, v := range ds.genNodeList() {

	//}
}
