// Package layertree provides the Chrome DevTools Protocol
// commands, types, and events for the LayerTree domain.
//
// Generated by the cdproto-gen command.
package layertree

// Code generated by cdproto-gen. DO NOT EDIT.

import (
	"context"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/mailru/easyjson"
)

// CompositingReasonsParams provides the reasons why the given layer was
// composited.
type CompositingReasonsParams struct {
	LayerID LayerID `json:"layerId"` // The id of the layer for which we want to get the reasons it was composited.
}

// CompositingReasons provides the reasons why the given layer was
// composited.
//
// parameters:
//   layerID - The id of the layer for which we want to get the reasons it was composited.
func CompositingReasons(layerID LayerID) *CompositingReasonsParams {
	return &CompositingReasonsParams{
		LayerID: layerID,
	}
}

// CompositingReasonsReturns return values.
type CompositingReasonsReturns struct {
	CompositingReasons []string `json:"compositingReasons,omitempty"` // A list of strings specifying reasons for the given layer to become composited.
}

// Do executes LayerTree.compositingReasons against the provided context.
//
// returns:
//   compositingReasons - A list of strings specifying reasons for the given layer to become composited.
func (p *CompositingReasonsParams) Do(ctxt context.Context) (compositingReasons []string, err error) {
	// execute
	var res CompositingReasonsReturns
	err = cdp.Execute(ctxt, CommandCompositingReasons, p, &res)
	if err != nil {
		return nil, err
	}

	return res.CompositingReasons, nil
}

// DisableParams disables compositing tree inspection.
type DisableParams struct{}

// Disable disables compositing tree inspection.
func Disable() *DisableParams {
	return &DisableParams{}
}

// Do executes LayerTree.disable against the provided context.
func (p *DisableParams) Do(ctxt context.Context) (err error) {
	return cdp.Execute(ctxt, CommandDisable, nil, nil)
}

// EnableParams enables compositing tree inspection.
type EnableParams struct{}

// Enable enables compositing tree inspection.
func Enable() *EnableParams {
	return &EnableParams{}
}

// Do executes LayerTree.enable against the provided context.
func (p *EnableParams) Do(ctxt context.Context) (err error) {
	return cdp.Execute(ctxt, CommandEnable, nil, nil)
}

// LoadSnapshotParams returns the snapshot identifier.
type LoadSnapshotParams struct {
	Tiles []*PictureTile `json:"tiles"` // An array of tiles composing the snapshot.
}

// LoadSnapshot returns the snapshot identifier.
//
// parameters:
//   tiles - An array of tiles composing the snapshot.
func LoadSnapshot(tiles []*PictureTile) *LoadSnapshotParams {
	return &LoadSnapshotParams{
		Tiles: tiles,
	}
}

// LoadSnapshotReturns return values.
type LoadSnapshotReturns struct {
	SnapshotID SnapshotID `json:"snapshotId,omitempty"` // The id of the snapshot.
}

// Do executes LayerTree.loadSnapshot against the provided context.
//
// returns:
//   snapshotID - The id of the snapshot.
func (p *LoadSnapshotParams) Do(ctxt context.Context) (snapshotID SnapshotID, err error) {
	// execute
	var res LoadSnapshotReturns
	err = cdp.Execute(ctxt, CommandLoadSnapshot, p, &res)
	if err != nil {
		return "", err
	}

	return res.SnapshotID, nil
}

// MakeSnapshotParams returns the layer snapshot identifier.
type MakeSnapshotParams struct {
	LayerID LayerID `json:"layerId"` // The id of the layer.
}

// MakeSnapshot returns the layer snapshot identifier.
//
// parameters:
//   layerID - The id of the layer.
func MakeSnapshot(layerID LayerID) *MakeSnapshotParams {
	return &MakeSnapshotParams{
		LayerID: layerID,
	}
}

// MakeSnapshotReturns return values.
type MakeSnapshotReturns struct {
	SnapshotID SnapshotID `json:"snapshotId,omitempty"` // The id of the layer snapshot.
}

// Do executes LayerTree.makeSnapshot against the provided context.
//
// returns:
//   snapshotID - The id of the layer snapshot.
func (p *MakeSnapshotParams) Do(ctxt context.Context) (snapshotID SnapshotID, err error) {
	// execute
	var res MakeSnapshotReturns
	err = cdp.Execute(ctxt, CommandMakeSnapshot, p, &res)
	if err != nil {
		return "", err
	}

	return res.SnapshotID, nil
}

// ProfileSnapshotParams [no description].
type ProfileSnapshotParams struct {
	SnapshotID     SnapshotID `json:"snapshotId"`               // The id of the layer snapshot.
	MinRepeatCount int64      `json:"minRepeatCount,omitempty"` // The maximum number of times to replay the snapshot (1, if not specified).
	MinDuration    float64    `json:"minDuration,omitempty"`    // The minimum duration (in seconds) to replay the snapshot.
	ClipRect       *dom.Rect  `json:"clipRect,omitempty"`       // The clip rectangle to apply when replaying the snapshot.
}

// ProfileSnapshot [no description].
//
// parameters:
//   snapshotID - The id of the layer snapshot.
func ProfileSnapshot(snapshotID SnapshotID) *ProfileSnapshotParams {
	return &ProfileSnapshotParams{
		SnapshotID: snapshotID,
	}
}

// WithMinRepeatCount the maximum number of times to replay the snapshot (1,
// if not specified).
func (p ProfileSnapshotParams) WithMinRepeatCount(minRepeatCount int64) *ProfileSnapshotParams {
	p.MinRepeatCount = minRepeatCount
	return &p
}

// WithMinDuration the minimum duration (in seconds) to replay the snapshot.
func (p ProfileSnapshotParams) WithMinDuration(minDuration float64) *ProfileSnapshotParams {
	p.MinDuration = minDuration
	return &p
}

// WithClipRect the clip rectangle to apply when replaying the snapshot.
func (p ProfileSnapshotParams) WithClipRect(clipRect *dom.Rect) *ProfileSnapshotParams {
	p.ClipRect = clipRect
	return &p
}

// ProfileSnapshotReturns return values.
type ProfileSnapshotReturns struct {
	Timings []PaintProfile `json:"timings,omitempty"` // The array of paint profiles, one per run.
}

// Do executes LayerTree.profileSnapshot against the provided context.
//
// returns:
//   timings - The array of paint profiles, one per run.
func (p *ProfileSnapshotParams) Do(ctxt context.Context) (timings []PaintProfile, err error) {
	// execute
	var res ProfileSnapshotReturns
	err = cdp.Execute(ctxt, CommandProfileSnapshot, p, &res)
	if err != nil {
		return nil, err
	}

	return res.Timings, nil
}

// ReleaseSnapshotParams releases layer snapshot captured by the back-end.
type ReleaseSnapshotParams struct {
	SnapshotID SnapshotID `json:"snapshotId"` // The id of the layer snapshot.
}

// ReleaseSnapshot releases layer snapshot captured by the back-end.
//
// parameters:
//   snapshotID - The id of the layer snapshot.
func ReleaseSnapshot(snapshotID SnapshotID) *ReleaseSnapshotParams {
	return &ReleaseSnapshotParams{
		SnapshotID: snapshotID,
	}
}

// Do executes LayerTree.releaseSnapshot against the provided context.
func (p *ReleaseSnapshotParams) Do(ctxt context.Context) (err error) {
	return cdp.Execute(ctxt, CommandReleaseSnapshot, p, nil)
}

// ReplaySnapshotParams replays the layer snapshot and returns the resulting
// bitmap.
type ReplaySnapshotParams struct {
	SnapshotID SnapshotID `json:"snapshotId"`         // The id of the layer snapshot.
	FromStep   int64      `json:"fromStep,omitempty"` // The first step to replay from (replay from the very start if not specified).
	ToStep     int64      `json:"toStep,omitempty"`   // The last step to replay to (replay till the end if not specified).
	Scale      float64    `json:"scale,omitempty"`    // The scale to apply while replaying (defaults to 1).
}

// ReplaySnapshot replays the layer snapshot and returns the resulting
// bitmap.
//
// parameters:
//   snapshotID - The id of the layer snapshot.
func ReplaySnapshot(snapshotID SnapshotID) *ReplaySnapshotParams {
	return &ReplaySnapshotParams{
		SnapshotID: snapshotID,
	}
}

// WithFromStep the first step to replay from (replay from the very start if
// not specified).
func (p ReplaySnapshotParams) WithFromStep(fromStep int64) *ReplaySnapshotParams {
	p.FromStep = fromStep
	return &p
}

// WithToStep the last step to replay to (replay till the end if not
// specified).
func (p ReplaySnapshotParams) WithToStep(toStep int64) *ReplaySnapshotParams {
	p.ToStep = toStep
	return &p
}

// WithScale the scale to apply while replaying (defaults to 1).
func (p ReplaySnapshotParams) WithScale(scale float64) *ReplaySnapshotParams {
	p.Scale = scale
	return &p
}

// ReplaySnapshotReturns return values.
type ReplaySnapshotReturns struct {
	DataURL string `json:"dataURL,omitempty"` // A data: URL for resulting image.
}

// Do executes LayerTree.replaySnapshot against the provided context.
//
// returns:
//   dataURL - A data: URL for resulting image.
func (p *ReplaySnapshotParams) Do(ctxt context.Context) (dataURL string, err error) {
	// execute
	var res ReplaySnapshotReturns
	err = cdp.Execute(ctxt, CommandReplaySnapshot, p, &res)
	if err != nil {
		return "", err
	}

	return res.DataURL, nil
}

// SnapshotCommandLogParams replays the layer snapshot and returns canvas
// log.
type SnapshotCommandLogParams struct {
	SnapshotID SnapshotID `json:"snapshotId"` // The id of the layer snapshot.
}

// SnapshotCommandLog replays the layer snapshot and returns canvas log.
//
// parameters:
//   snapshotID - The id of the layer snapshot.
func SnapshotCommandLog(snapshotID SnapshotID) *SnapshotCommandLogParams {
	return &SnapshotCommandLogParams{
		SnapshotID: snapshotID,
	}
}

// SnapshotCommandLogReturns return values.
type SnapshotCommandLogReturns struct {
	CommandLog []easyjson.RawMessage `json:"commandLog,omitempty"` // The array of canvas function calls.
}

// Do executes LayerTree.snapshotCommandLog against the provided context.
//
// returns:
//   commandLog - The array of canvas function calls.
func (p *SnapshotCommandLogParams) Do(ctxt context.Context) (commandLog []easyjson.RawMessage, err error) {
	// execute
	var res SnapshotCommandLogReturns
	err = cdp.Execute(ctxt, CommandSnapshotCommandLog, p, &res)
	if err != nil {
		return nil, err
	}

	return res.CommandLog, nil
}

// Command names.
const (
	CommandCompositingReasons = "LayerTree.compositingReasons"
	CommandDisable            = "LayerTree.disable"
	CommandEnable             = "LayerTree.enable"
	CommandLoadSnapshot       = "LayerTree.loadSnapshot"
	CommandMakeSnapshot       = "LayerTree.makeSnapshot"
	CommandProfileSnapshot    = "LayerTree.profileSnapshot"
	CommandReleaseSnapshot    = "LayerTree.releaseSnapshot"
	CommandReplaySnapshot     = "LayerTree.replaySnapshot"
	CommandSnapshotCommandLog = "LayerTree.snapshotCommandLog"
)
