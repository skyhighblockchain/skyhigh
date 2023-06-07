package eventmodule

import (
	"github.com/skyhighblockchain/skyhigh/gossip/blockproc"
	"github.com/skyhighblockchain/skyhigh/inter"
)

type ValidatorEventsModule struct{}

func New() *ValidatorEventsModule {
	return &ValidatorEventsModule{}
}

func (m *ValidatorEventsModule) Start(bs blockproc.BlockState, es blockproc.EpochState) blockproc.ConfirmedEventsProcessor {
	return &ValidatorEventsProcessor{
		es:                     es,
		bs:                     bs,
		validatorHighestEvents: make(inter.EventIs, es.Validators.Len()),
	}
}

type ValidatorEventsProcessor struct {
	es                     blockproc.EpochState
	bs                     blockproc.BlockState
	validatorHighestEvents inter.EventIs
}

func (p *ValidatorEventsProcessor) ProcessConfirmedEvent(e inter.EventI) {
	creatorIdx := p.es.Validators.GetIdx(e.Creator())
	prev := p.validatorHighestEvents[creatorIdx]
	if prev == nil || e.Seq() > prev.Seq() {
		p.validatorHighestEvents[creatorIdx] = e
	}
	p.bs.EpochGas += e.GasPowerUsed()
}

func (p *ValidatorEventsProcessor) Finalize(block blockproc.BlockCtx, _ bool) blockproc.BlockState {
	for _, v := range p.bs.EpochCheaters {
		creatorIdx := p.es.Validators.GetIdx(v)
		p.validatorHighestEvents[creatorIdx] = nil
	}
	for creatorIdx, e := range p.validatorHighestEvents {
		if e == nil {
			continue
		}
		info := p.bs.ValidatorStates[creatorIdx]
		if block.Idx <= info.LastBlock+p.es.Rules.Economy.BlockMissedSlack {
			prevOnlineTime := info.LastOnlineTime
			if p.es.Rules.Upgrades.Berlin {
				prevOnlineTime = inter.MaxTimestamp(info.LastOnlineTime, p.es.EpochStart)
			}
			if e.MedianTime() > prevOnlineTime {
				info.Uptime += e.MedianTime() - prevOnlineTime
			}
		}
		info.LastGasPowerLeft = e.GasPowerLeft()
		info.LastOnlineTime = e.MedianTime()
		info.LastBlock = block.Idx
		info.LastEvent = e.ID()
		p.bs.ValidatorStates[creatorIdx] = info
	}
	return p.bs
}
