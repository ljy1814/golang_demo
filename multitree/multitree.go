package multitree

type MTreeNode struct {
	ID  int64
	PID int64
	//	Value interface{}
	Value string
}

type MTree struct {
	Children []*MTreeNode
	Parent   *MTreeNode
}

func MergeTree(mns []MTreeNode) map[int64]map[int64]*MTreeNode {
	if len(mns) <= 0 {
		return nil
	}
	// PID : ID -> NODE
	var mtp = make(map[int64]map[int64]*MTreeNode)

	for i, v := range mns {
		if len(mtp[v.PID]) <= 0 {
			mtp[v.PID] = make(map[int64]*MTreeNode)
		}
		if v.PID == 0 {
			//mtp[v.PID][v.ID] = &mns[i]
			//mtp[v.PID][v.ID] = append(mtp[v.PID][v.ID], &mns[i])
			mtp[v.PID][v.ID] = &mns[i]
		} else {
			mtp[v.PID][v.ID] = &mns[i]
		}
	}
	return mtp
}
