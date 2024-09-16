package list

import (
	"sync/atomic"
	"unsafe"
)

type NodeList struct {
	firstNode unsafe.Pointer
}

type Node struct {
	Data any
	Next unsafe.Pointer
}

func NewNodeList() NodeList {
	return NodeList{}
}

func (l *NodeList) AddRecord(item any) {
	newNode := Node{Data: item, Next: l.firstNode}
	atomic.StorePointer(&l.firstNode, unsafe.Pointer(&newNode))
}

func (l *NodeList) GetRecords(count int) (result []any) {
	nodeLink := l.firstNode
	for nodeLink != nil && count > 0 {
		val := (*Node)(nodeLink)
		result = append(result, val.Data)
		nodeLink = val.Next
		count--
	}

	return
}
