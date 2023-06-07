package eventcheck

import (
	"github.com/skyhighblockchain/skyhigh/eventcheck/basiccheck"
	"github.com/skyhighblockchain/skyhigh/eventcheck/epochcheck"
	"github.com/skyhighblockchain/skyhigh/eventcheck/gaspowercheck"
	"github.com/skyhighblockchain/skyhigh/eventcheck/heavycheck"
	"github.com/skyhighblockchain/skyhigh/eventcheck/parentscheck"
	"github.com/skyhighblockchain/skyhigh/inter"
)

// Checkers is collection of all the checkers
type Checkers struct {
	Basiccheck    *basiccheck.Checker
	Epochcheck    *epochcheck.Checker
	Parentscheck  *parentscheck.Checker
	Gaspowercheck *gaspowercheck.Checker
	Heavycheck    *heavycheck.Checker
}

// Validate runs all the checks except Poset-related
func (v *Checkers) Validate(e inter.EventPayloadI, parents inter.EventIs) error {
	if err := v.Basiccheck.Validate(e); err != nil {
		return err
	}
	if err := v.Epochcheck.Validate(e); err != nil {
		return err
	}
	if err := v.Parentscheck.Validate(e, parents); err != nil {
		return err
	}
	var selfParent inter.EventI
	if e.SelfParent() != nil {
		selfParent = parents[0]
	}
	if err := v.Gaspowercheck.Validate(e, selfParent); err != nil {
		return err
	}
	if err := v.Heavycheck.Validate(e); err != nil {
		return err
	}
	return nil
}
