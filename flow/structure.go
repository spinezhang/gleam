package flow

import (
	"io"
	"sync"
	"time"

	"github.com/chrislusf/gleam/instruction"
	"github.com/chrislusf/gleam/script"
	"github.com/chrislusf/gleam/util"
)

type NetworkType int
type ModeIO int

const (
	OneShardToOneShard NetworkType = iota
	OneShardToAllShard
	AllShardToOneShard
	OneShardToEveryNShard
	LinkedNShardToOneShard
	MergeTwoShardToOneShard
)

const (
	ModeInMemory ModeIO = iota
	ModeOnDisk
)

type DasetsetMetadata struct {
	TotalSize  int64
	OnDisk     ModeIO
	Datacenter string
	Rack       string
}

type StepMetadata struct {
	IsRestartable bool
}

type FlowContext struct {
	PrevScriptType string
	PrevScriptPart string
	Scripts        map[string]func() script.Script
	Steps          []*Step
	Datasets       []*Dataset
	HashCode       uint32
}

type Dataset struct {
	FlowContext     *FlowContext
	Id              int
	Shards          []*DatasetShard
	Step            *Step
	ReadingSteps    []*Step
	IsPartitionedBy []int
	IsLocalSorted   []instruction.OrderBy
	Meta            *DasetsetMetadata
	RunLocked
}

type DatasetShard struct {
	Id            int
	Dataset       *Dataset
	ReadingTasks  []*Task
	IncomingChan  *util.Piper
	OutgoingChans []*util.Piper
	Counter       int64
	ReadyTime     time.Time
	CloseTime     time.Time
}

type Step struct {
	Id             int
	FlowContext    *FlowContext
	InputDatasets  []*Dataset
	OutputDataset  *Dataset
	Function       func([]io.Reader, []io.Writer, *instruction.Stats) error
	Instruction    instruction.Instruction
	Tasks          []*Task
	Name           string
	NetworkType    NetworkType
	IsOnDriverSide bool
	IsPipe         bool
	Script         script.Script
	Command        *script.Command // used in Pipe()
	Meta           *StepMetadata
	Params         map[string]interface{}
	RunLocked
}

type Task struct {
	Id           int
	Step         *Step
	InputShards  []*DatasetShard
	InputChans   []*util.Piper // task specific input chans. InputShard may have multiple reading tasks
	OutputShards []*DatasetShard
	Stats        *instruction.Stats
}

type RunLocked struct {
	sync.Mutex
	StartTime time.Time
}
